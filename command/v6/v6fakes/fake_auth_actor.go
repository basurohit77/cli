// Code generated by counterfeiter. DO NOT EDIT.
package v6fakes

import (
	"sync"

	"code.cloudfoundry.org/cli/api/uaa/constant"
	v6 "code.cloudfoundry.org/cli/command/v6"
)

type FakeAuthActor struct {
	AuthenticateStub        func(string, string, string, constant.GrantType) error
	authenticateMutex       sync.RWMutex
	authenticateArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 constant.GrantType
	}
	authenticateReturns struct {
		result1 error
	}
	authenticateReturnsOnCall map[int]struct {
		result1 error
	}
	CloudControllerAPIVersionStub        func() string
	cloudControllerAPIVersionMutex       sync.RWMutex
	cloudControllerAPIVersionArgsForCall []struct {
	}
	cloudControllerAPIVersionReturns struct {
		result1 string
	}
	cloudControllerAPIVersionReturnsOnCall map[int]struct {
		result1 string
	}
	RevokeStub        func() error
	revokeMutex       sync.RWMutex
	revokeArgsForCall []struct {
	}
	revokeReturns struct {
		result1 error
	}
	revokeReturnsOnCall map[int]struct {
		result1 error
	}
	UAAAPIVersionStub        func() string
	uAAAPIVersionMutex       sync.RWMutex
	uAAAPIVersionArgsForCall []struct {
	}
	uAAAPIVersionReturns struct {
		result1 string
	}
	uAAAPIVersionReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeAuthActor) Authenticate(arg1 string, arg2 string, arg3 string, arg4 constant.GrantType) error {
	fake.authenticateMutex.Lock()
	ret, specificReturn := fake.authenticateReturnsOnCall[len(fake.authenticateArgsForCall)]
	fake.authenticateArgsForCall = append(fake.authenticateArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 constant.GrantType
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("Authenticate", []interface{}{arg1, arg2, arg3, arg4})
	fake.authenticateMutex.Unlock()
	if fake.AuthenticateStub != nil {
		return fake.AuthenticateStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.authenticateReturns
	return fakeReturns.result1
}

func (fake *FakeAuthActor) AuthenticateCallCount() int {
	fake.authenticateMutex.RLock()
	defer fake.authenticateMutex.RUnlock()
	return len(fake.authenticateArgsForCall)
}

func (fake *FakeAuthActor) AuthenticateCalls(stub func(string, string, string, constant.GrantType) error) {
	fake.authenticateMutex.Lock()
	defer fake.authenticateMutex.Unlock()
	fake.AuthenticateStub = stub
}

func (fake *FakeAuthActor) AuthenticateArgsForCall(i int) (string, string, string, constant.GrantType) {
	fake.authenticateMutex.RLock()
	defer fake.authenticateMutex.RUnlock()
	argsForCall := fake.authenticateArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeAuthActor) AuthenticateReturns(result1 error) {
	fake.authenticateMutex.Lock()
	defer fake.authenticateMutex.Unlock()
	fake.AuthenticateStub = nil
	fake.authenticateReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeAuthActor) AuthenticateReturnsOnCall(i int, result1 error) {
	fake.authenticateMutex.Lock()
	defer fake.authenticateMutex.Unlock()
	fake.AuthenticateStub = nil
	if fake.authenticateReturnsOnCall == nil {
		fake.authenticateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.authenticateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeAuthActor) CloudControllerAPIVersion() string {
	fake.cloudControllerAPIVersionMutex.Lock()
	ret, specificReturn := fake.cloudControllerAPIVersionReturnsOnCall[len(fake.cloudControllerAPIVersionArgsForCall)]
	fake.cloudControllerAPIVersionArgsForCall = append(fake.cloudControllerAPIVersionArgsForCall, struct {
	}{})
	fake.recordInvocation("CloudControllerAPIVersion", []interface{}{})
	fake.cloudControllerAPIVersionMutex.Unlock()
	if fake.CloudControllerAPIVersionStub != nil {
		return fake.CloudControllerAPIVersionStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.cloudControllerAPIVersionReturns
	return fakeReturns.result1
}

func (fake *FakeAuthActor) CloudControllerAPIVersionCallCount() int {
	fake.cloudControllerAPIVersionMutex.RLock()
	defer fake.cloudControllerAPIVersionMutex.RUnlock()
	return len(fake.cloudControllerAPIVersionArgsForCall)
}

func (fake *FakeAuthActor) CloudControllerAPIVersionCalls(stub func() string) {
	fake.cloudControllerAPIVersionMutex.Lock()
	defer fake.cloudControllerAPIVersionMutex.Unlock()
	fake.CloudControllerAPIVersionStub = stub
}

func (fake *FakeAuthActor) CloudControllerAPIVersionReturns(result1 string) {
	fake.cloudControllerAPIVersionMutex.Lock()
	defer fake.cloudControllerAPIVersionMutex.Unlock()
	fake.CloudControllerAPIVersionStub = nil
	fake.cloudControllerAPIVersionReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeAuthActor) CloudControllerAPIVersionReturnsOnCall(i int, result1 string) {
	fake.cloudControllerAPIVersionMutex.Lock()
	defer fake.cloudControllerAPIVersionMutex.Unlock()
	fake.CloudControllerAPIVersionStub = nil
	if fake.cloudControllerAPIVersionReturnsOnCall == nil {
		fake.cloudControllerAPIVersionReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.cloudControllerAPIVersionReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeAuthActor) Revoke() error {
	fake.revokeMutex.Lock()
	ret, specificReturn := fake.revokeReturnsOnCall[len(fake.revokeArgsForCall)]
	fake.revokeArgsForCall = append(fake.revokeArgsForCall, struct {
	}{})
	fake.recordInvocation("Revoke", []interface{}{})
	fake.revokeMutex.Unlock()
	if fake.RevokeStub != nil {
		return fake.RevokeStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.revokeReturns
	return fakeReturns.result1
}

func (fake *FakeAuthActor) RevokeCallCount() int {
	fake.revokeMutex.RLock()
	defer fake.revokeMutex.RUnlock()
	return len(fake.revokeArgsForCall)
}

func (fake *FakeAuthActor) RevokeCalls(stub func() error) {
	fake.revokeMutex.Lock()
	defer fake.revokeMutex.Unlock()
	fake.RevokeStub = stub
}

func (fake *FakeAuthActor) RevokeReturns(result1 error) {
	fake.revokeMutex.Lock()
	defer fake.revokeMutex.Unlock()
	fake.RevokeStub = nil
	fake.revokeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeAuthActor) RevokeReturnsOnCall(i int, result1 error) {
	fake.revokeMutex.Lock()
	defer fake.revokeMutex.Unlock()
	fake.RevokeStub = nil
	if fake.revokeReturnsOnCall == nil {
		fake.revokeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.revokeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeAuthActor) UAAAPIVersion() string {
	fake.uAAAPIVersionMutex.Lock()
	ret, specificReturn := fake.uAAAPIVersionReturnsOnCall[len(fake.uAAAPIVersionArgsForCall)]
	fake.uAAAPIVersionArgsForCall = append(fake.uAAAPIVersionArgsForCall, struct {
	}{})
	fake.recordInvocation("UAAAPIVersion", []interface{}{})
	fake.uAAAPIVersionMutex.Unlock()
	if fake.UAAAPIVersionStub != nil {
		return fake.UAAAPIVersionStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.uAAAPIVersionReturns
	return fakeReturns.result1
}

func (fake *FakeAuthActor) UAAAPIVersionCallCount() int {
	fake.uAAAPIVersionMutex.RLock()
	defer fake.uAAAPIVersionMutex.RUnlock()
	return len(fake.uAAAPIVersionArgsForCall)
}

func (fake *FakeAuthActor) UAAAPIVersionCalls(stub func() string) {
	fake.uAAAPIVersionMutex.Lock()
	defer fake.uAAAPIVersionMutex.Unlock()
	fake.UAAAPIVersionStub = stub
}

func (fake *FakeAuthActor) UAAAPIVersionReturns(result1 string) {
	fake.uAAAPIVersionMutex.Lock()
	defer fake.uAAAPIVersionMutex.Unlock()
	fake.UAAAPIVersionStub = nil
	fake.uAAAPIVersionReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeAuthActor) UAAAPIVersionReturnsOnCall(i int, result1 string) {
	fake.uAAAPIVersionMutex.Lock()
	defer fake.uAAAPIVersionMutex.Unlock()
	fake.UAAAPIVersionStub = nil
	if fake.uAAAPIVersionReturnsOnCall == nil {
		fake.uAAAPIVersionReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.uAAAPIVersionReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeAuthActor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.authenticateMutex.RLock()
	defer fake.authenticateMutex.RUnlock()
	fake.cloudControllerAPIVersionMutex.RLock()
	defer fake.cloudControllerAPIVersionMutex.RUnlock()
	fake.revokeMutex.RLock()
	defer fake.revokeMutex.RUnlock()
	fake.uAAAPIVersionMutex.RLock()
	defer fake.uAAAPIVersionMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeAuthActor) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ v6.AuthActor = new(FakeAuthActor)
