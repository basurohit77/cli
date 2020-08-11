package v7action

import (
	"code.cloudfoundry.org/cli/actor/actionerror"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccerror"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3"
	"code.cloudfoundry.org/cli/api/cloudcontroller/ccv3/constant"
	"code.cloudfoundry.org/cli/resources"
	"code.cloudfoundry.org/cli/types"
	"code.cloudfoundry.org/cli/util/railway"
)

type ServiceInstanceUpdateManagedParams struct {
	ServicePlanName types.OptionalString
	Tags            types.OptionalStringSlice
	Parameters      types.OptionalObject
}

func (actor Actor) GetServiceInstanceByNameAndSpace(serviceInstanceName string, spaceGUID string) (resources.ServiceInstance, Warnings, error) {
	serviceInstance, _, warnings, err := actor.CloudControllerClient.GetServiceInstanceByNameAndSpace(serviceInstanceName, spaceGUID)
	switch e := err.(type) {
	case ccerror.ServiceInstanceNotFoundError:
		return serviceInstance, Warnings(warnings), actionerror.ServiceInstanceNotFoundError{Name: e.Name}
	default:
		return serviceInstance, Warnings(warnings), err
	}
}

func (actor Actor) UnshareServiceInstanceByServiceInstanceAndSpace(serviceInstanceGUID string, sharedToSpaceGUID string) (Warnings, error) {
	warnings, err := actor.CloudControllerClient.DeleteServiceInstanceRelationshipsSharedSpace(serviceInstanceGUID, sharedToSpaceGUID)
	return Warnings(warnings), err
}

func (actor Actor) CreateUserProvidedServiceInstance(serviceInstance resources.ServiceInstance) (Warnings, error) {
	serviceInstance.Type = resources.UserProvidedServiceInstance
	_, warnings, err := actor.CloudControllerClient.CreateServiceInstance(serviceInstance)
	return Warnings(warnings), err
}

func (actor Actor) UpdateUserProvidedServiceInstance(serviceInstanceName, spaceGUID string, serviceInstanceUpdates resources.ServiceInstance) (Warnings, error) {
	var original resources.ServiceInstance

	warnings, err := railway.Sequentially(
		func() (warnings ccv3.Warnings, err error) {
			original, _, warnings, err = actor.CloudControllerClient.GetServiceInstanceByNameAndSpace(serviceInstanceName, spaceGUID)
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			err = assertServiceInstanceType(resources.UserProvidedServiceInstance, original)
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			_, warnings, err = actor.CloudControllerClient.UpdateServiceInstance(original.GUID, serviceInstanceUpdates)
			return
		},
	)

	return Warnings(warnings), err
}

func (actor Actor) UpdateManagedServiceInstance(serviceInstanceName, spaceGUID string, serviceInstanceUpdates ServiceInstanceUpdateManagedParams) (Warnings, error) {
	var (
		original             resources.ServiceInstance
		jobURL               ccv3.JobURL
		allWarnings          Warnings
		includedResources    ccv3.IncludedResources
		serviceInstanceQuery []ccv3.Query
	)

	updates := resources.ServiceInstance{
		Tags:       serviceInstanceUpdates.Tags,
		Parameters: serviceInstanceUpdates.Parameters,
	}

	if serviceInstanceUpdates.ServicePlanName.IsSet {
		serviceInstanceQuery = []ccv3.Query{
			{
				Key:    ccv3.FieldsServicePlanServiceOffering,
				Values: []string{"name"},
			},
			{
				Key:    ccv3.FieldsServicePlanServiceOfferingServiceBroker,
				Values: []string{"name"},
			},
		}
	}
	warnings, err := handleServiceInstanceErrors(railway.Sequentially(
		func() (warnings ccv3.Warnings, err error) {
			original, includedResources, warnings, err = actor.CloudControllerClient.GetServiceInstanceByNameAndSpace(serviceInstanceName, spaceGUID, serviceInstanceQuery...)
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			err = assertServiceInstanceType(resources.ManagedServiceInstance, original)
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			if !serviceInstanceUpdates.ServicePlanName.IsSet {
				return
			}
			_, actorWarnings, err := actor.GetServicePlanByNameOfferingAndBroker(
				serviceInstanceUpdates.ServicePlanName.Value,
				includedResources.ServiceOfferings[0].Name,
				includedResources.ServiceBrokers[0].Name,
			)
			if err != nil {
				return ccv3.Warnings(actorWarnings), err
			}
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			jobURL, warnings, err = actor.CloudControllerClient.UpdateServiceInstance(original.GUID, updates)
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			return actor.CloudControllerClient.PollJobForState(jobURL, constant.JobPolling)
		},
	))

	allWarnings = append(allWarnings, warnings...)

	return allWarnings, err
}

func (actor Actor) UpgradeServiceInstance(serviceInstanceName string, spaceGUID string) (Warnings, error) {
	var serviceInstance resources.ServiceInstance
	var servicePlan resources.ServicePlan
	var jobURL ccv3.JobURL

	return handleServiceInstanceErrors(railway.Sequentially(
		func() (warnings ccv3.Warnings, err error) {
			serviceInstance, _, warnings, err = actor.CloudControllerClient.GetServiceInstanceByNameAndSpace(serviceInstanceName, spaceGUID)
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			if serviceInstance.UpgradeAvailable.Value != true {
				err = actionerror.ServiceInstanceUpgradeNotAvailableError{}
			}
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			servicePlan, warnings, err = actor.CloudControllerClient.GetServicePlanByGUID(serviceInstance.ServicePlanGUID)
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			jobURL, warnings, err = actor.CloudControllerClient.UpdateServiceInstance(serviceInstance.GUID, resources.ServiceInstance{
				MaintenanceInfoVersion: servicePlan.MaintenanceInfoVersion,
			})
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			return actor.CloudControllerClient.PollJobForState(jobURL, constant.JobPolling)
		},
	))
}

func (actor Actor) RenameServiceInstance(currentServiceInstanceName, spaceGUID, newServiceInstanceName string) (Warnings, error) {
	var (
		serviceInstance resources.ServiceInstance
		jobURL          ccv3.JobURL
	)

	return handleServiceInstanceErrors(railway.Sequentially(
		func() (warnings ccv3.Warnings, err error) {
			serviceInstance, _, warnings, err = actor.CloudControllerClient.GetServiceInstanceByNameAndSpace(currentServiceInstanceName, spaceGUID)
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			jobURL, warnings, err = actor.CloudControllerClient.UpdateServiceInstance(
				serviceInstance.GUID,
				resources.ServiceInstance{Name: newServiceInstanceName},
			)
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			return actor.CloudControllerClient.PollJobForState(jobURL, constant.JobPolling)
		},
	))
}

func (actor Actor) fetchServiceInstanceDetails(serviceInstanceName string, spaceGUID string) (resources.ServiceInstance, ccv3.IncludedResources, Warnings, error) {
	query := []ccv3.Query{
		{
			Key:    ccv3.FieldsServicePlan,
			Values: []string{"name", "guid"},
		},
		{
			Key:    ccv3.FieldsServicePlanServiceOffering,
			Values: []string{"name", "guid", "description", "tags", "documentation_url"},
		},
		{
			Key:    ccv3.FieldsServicePlanServiceOfferingServiceBroker,
			Values: []string{"name", "guid"},
		},
	}

	serviceInstance, included, warnings, err := actor.CloudControllerClient.GetServiceInstanceByNameAndSpace(serviceInstanceName, spaceGUID, query...)
	switch err.(type) {
	case nil:
	case ccerror.ServiceInstanceNotFoundError:
		return resources.ServiceInstance{}, ccv3.IncludedResources{}, Warnings(warnings), actionerror.ServiceInstanceNotFoundError{Name: serviceInstanceName}
	default:
		return resources.ServiceInstance{}, ccv3.IncludedResources{}, Warnings(warnings), err
	}

	return serviceInstance, included, Warnings(warnings), nil
}

type ManagedServiceInstanceParams struct {
	ServiceOfferingName string
	ServicePlanName     string
	ServiceInstanceName string
	ServiceBrokerName   string
	SpaceGUID           string
	Tags                types.OptionalStringSlice
	Parameters          types.OptionalObject
}

func (actor Actor) CreateManagedServiceInstance(params ManagedServiceInstanceParams) (Warnings, error) {
	allWarnings := Warnings{}

	servicePlan, warnings, err := actor.GetServicePlanByNameOfferingAndBroker(
		params.ServicePlanName,
		params.ServiceOfferingName,
		params.ServiceBrokerName,
	)
	allWarnings = append(allWarnings, warnings...)
	if err != nil {
		return allWarnings, err
	}

	serviceInstance := resources.ServiceInstance{
		Type:            resources.ManagedServiceInstance,
		Name:            params.ServiceInstanceName,
		ServicePlanGUID: servicePlan.GUID,
		SpaceGUID:       params.SpaceGUID,
		Tags:            params.Tags,
		Parameters:      params.Parameters,
	}

	jobURL, clientWarnings, err := actor.CloudControllerClient.CreateServiceInstance(serviceInstance)
	allWarnings = append(allWarnings, clientWarnings...)
	if err != nil {
		return allWarnings, err
	}

	clientWarnings, err = actor.CloudControllerClient.PollJobForState(jobURL, constant.JobPolling)
	allWarnings = append(allWarnings, clientWarnings...)

	return allWarnings, err

}

type ServiceInstanceDeleteState int

const (
	ServiceInstanceUnknownState ServiceInstanceDeleteState = iota
	ServiceInstanceDidNotExist
	ServiceInstanceGone
	ServiceInstanceDeleteInProgress
)

func (actor Actor) DeleteServiceInstance(serviceInstanceName, spaceGUID string, wait bool) (ServiceInstanceDeleteState, Warnings, error) {
	var (
		serviceInstance resources.ServiceInstance
		jobURL          ccv3.JobURL
	)

	warnings, err := railway.Sequentially(
		func() (warnings ccv3.Warnings, err error) {
			serviceInstance, _, warnings, err = actor.CloudControllerClient.GetServiceInstanceByNameAndSpace(serviceInstanceName, spaceGUID)
			return
		},
		func() (warnings ccv3.Warnings, err error) {
			jobURL, warnings, err = actor.CloudControllerClient.DeleteServiceInstance(serviceInstance.GUID)
			return
		},
		func() (ccv3.Warnings, error) {
			return actor.pollJob(jobURL, wait)
		},
	)

	switch err.(type) {
	case nil:
	case ccerror.ServiceInstanceNotFoundError:
		return ServiceInstanceDidNotExist, Warnings(warnings), nil
	default:
		return ServiceInstanceUnknownState, Warnings(warnings), err
	}

	if jobURL != "" && !wait {
		return ServiceInstanceDeleteInProgress, Warnings(warnings), nil
	}
	return ServiceInstanceGone, Warnings(warnings), nil
}

func (actor Actor) pollJob(jobURL ccv3.JobURL, wait bool) (ccv3.Warnings, error) {
	switch {
	case jobURL == "":
		return ccv3.Warnings{}, nil
	case wait:
		return actor.CloudControllerClient.PollJob(jobURL)
	default:
		return actor.CloudControllerClient.PollJobForState(jobURL, constant.JobPolling)
	}
}

func assertServiceInstanceType(requiredType resources.ServiceInstanceType, instance resources.ServiceInstance) error {
	if instance.Type != requiredType {
		return actionerror.ServiceInstanceTypeError{
			Name:         instance.Name,
			RequiredType: requiredType,
		}
	}

	return nil
}

func handleServiceInstanceErrors(warnings ccv3.Warnings, err error) (Warnings, error) {
	switch e := err.(type) {
	case nil:
		return Warnings(warnings), nil
	case ccerror.ServiceInstanceNotFoundError:
		return Warnings(warnings), actionerror.ServiceInstanceNotFoundError{Name: e.Name}
	default:
		return Warnings(warnings), err
	}
}
