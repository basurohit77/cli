package isolated

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	. "code.cloudfoundry.org/cli/cf/util/testhelpers/matchers"
	"code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = FDescribe("rollback command", func() {
	var (
		orgName   string
		spaceName string
		appName   string
	)

	BeforeEach(func() {
		orgName = helpers.NewOrgName()
		spaceName = helpers.NewSpaceName()
		appName = helpers.PrefixedRandomName("app")
	})

	Describe("help", func() {
		When("--help flag is set", func() {
			It("appears in cf help -a", func() {
				session := helpers.CF("help", "-a")
				Eventually(session).Should(Exit(0))
				Expect(session).To(HaveCommandInCategoryWithDescription("rollback", "EXPERIMENTAL COMMANDS", "Rollback to the specified revision of an app"))
			})

			It("Displays rollback command usage to output", func() {
				session := helpers.CF("rollback", "--help")

				Eventually(session).Should(Exit(0))

				Expect(session).To(Say("NAME:"))
				Expect(session).To(Say("rollback - Rollback to the specified revision of an app"))
				Expect(session).To(Say("USAGE:"))
				Expect(session).To(Say(`cf rollback APP_NAME \[--revision REVISION_NUMBER\]`))
				Expect(session).To(Say("OPTIONS:"))
				Expect(session).To(Say("--revision      Roll back to the given app revision"))
				Expect(session).To(Say("SEE ALSO:"))
				Expect(session).To(Say("revisions"))
			})
		})
	})

	When("the environment is set up correctly", func() {
		BeforeEach(func() {
			helpers.SetupCF(orgName, spaceName)
		})

		AfterEach(func() {
			helpers.QuickDeleteOrg(orgName)
		})

		Describe("version dependent display", func() {

			var domainName string

			BeforeEach(func() {
				domainName = helpers.DefaultSharedDomain()
			})

			When("the app is started and has 2 instances", func() {
				BeforeEach(func() {
					helpers.WithHelloWorldApp(func(appDir string) {
						manifestContents := []byte(fmt.Sprintf(`
---
applications:
- name: %s
  memory: 128M
  instances: 2
  disk_quota: 128M
  routes:
  - route: %s.%s
`, appName, appName, domainName))
						manifestPath := filepath.Join(appDir, "manifest.yml")
						err := ioutil.WriteFile(manifestPath, manifestContents, 0666)
						Expect(err).ToNot(HaveOccurred())

						// Create manifest
						Eventually(helpers.CF("push", appName, "-p", appDir, "-f", manifestPath, "-b", "staticfile_buildpack")).Should(Exit(0))
						Eventually(helpers.CF("push", appName, "-p", appDir, "-f", manifestPath, "-b", "staticfile_buildpack")).Should(Exit(0))
					})
				})

				When("the -f flag is not provided", func() {
					var buffer *Buffer

					BeforeEach(func() {
						buffer = NewBuffer()
					})
					When("the user enters y", func() {

						BeforeEach(func() {
							_, err := buffer.Write([]byte("y\n"))
							Expect(err).ToNot(HaveOccurred())
						})

						It("prompts the user to rollback, then successfully rolls back", func() {
							userName, _ := helpers.GetCredentials()

							session := helpers.CFWithStdin(buffer, "rollback", appName, "--revision", "1")
							Eventually(session).Should(Exit(0))

							Expect(session).To(Say(`Rolling '%s' back to revision '1' will create a new revision. The new revision '3' will use the settings from revision '1'.`, appName))
							Expect(session).To(Say(`Are you sure you want to continue? [y/N]>`))
							Expect(session).To(Say(`Rolling back to revision 1 for app %s in org %s / space %s as %s...`, appName, orgName, spaceName, userName))
							Expect(session).To(Say(`OK`))

							session = helpers.CF("revisions", appName)

							Eventually(session).Should(Exit(0))
							Expect(session).To(Say(`3\s+[\w\-]+\s+New droplet deployed.`))
						})
					})
				})
			})
		})
	})
})
