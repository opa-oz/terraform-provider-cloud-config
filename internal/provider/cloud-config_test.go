package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

const resourceName = "cloud-config.test"

func TestAccExampleResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccExampleResourceConfig("one"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("one"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("fqdn"),
						knownvalue.StringExact("one.lan"),
					),
				},
			},
			// Update and Read testing
			{
				Config: testAccExampleResourceConfig("two"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("two"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("fqdn"),
						knownvalue.StringExact("two.lan"),
					),
					statecheck.ExpectKnownValue(
						resourceName,
						tfjsonpath.New("content"),
						knownvalue.StringExact(strings.TrimSpace(`
#cloud-config
hostname: two
fqdn: two.lan
prefer_fqdn_over_hostname: true
preserve_hostname: true
create_hostname_file: false
locale: en_two
locale_configfile: /etc/locale
timezone: Asia/two
runcmd:
    - echo '11'
    - cat two
manage_etc_hosts: localhost
ssh_authorized_keys:
    - two
              `)),
					),
				},
			},
		},
	})
}

func testAccExampleResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "cloud-config" "test" {
  hostname = %[1]q
  fqdn = "%[1]s.lan"
  prefer_fqdn_over_hostname = true
  preserve_hostname = true
  create_hostname_file = false

  locale = "en_%[1]s"
  locale_configfile = "/etc/locale"

  timezone = "Asia/%[1]s"

	runcmd = ["echo '11'", "cat %[1]s"]

	manage_etc_hosts_localhost = true

	ssh_authorized_keys = [%[1]q]
}
`, configurableAttribute)
}

func expectedOutput(s string) string {
	return strings.TrimSpace(fmt.Sprintf(`
#cloud-config
%s
		`, strings.TrimSpace(s)))
}

func wrapInput(s string) string {
	return fmt.Sprintf(`
resource "cloud-config" "test" {
		%s
}
		`, s)
}

type testCase struct {
	name           string
	input          string
	expectedValues map[string]string // key: attribute name, value: expected string or "null"
	expectedOutput string
}

func assembleTestCase(testCases []testCase, t *testing.T) resource.TestCase {
	var steps []resource.TestStep
	for _, tc := range testCases {
		var checks []statecheck.StateCheck
		for attr, val := range tc.expectedValues {
			path := tfjsonpath.New(attr)

			if strings.Contains(attr, ".") {
				steps := strings.Split(attr, ".")
				path = tfjsonpath.New(steps[0])
				for i, step := range steps {
					if i == 0 {
						continue
					}

					if utils.IsNumeric(step) {
						stepVal, _ := utils.ToInt(step)
						path = path.AtSliceIndex(stepVal)
					} else {
						path = path.AtMapKey(step)
					}
				}
			}

			switch val {
			case "null":
				checks = append(checks, statecheck.ExpectKnownValue(
					resourceName,
					path,
					knownvalue.Null(),
				))
			case "true":
				checks = append(checks, statecheck.ExpectKnownValue(
					resourceName,
					path,
					knownvalue.Bool(true),
				))
			case "false":
				checks = append(checks, statecheck.ExpectKnownValue(
					resourceName,
					path,
					knownvalue.Bool(false),
				))
			default:
				checks = append(checks, statecheck.ExpectKnownValue(
					resourceName,
					path,
					knownvalue.StringExact(val),
				))
			}
		}

		checks = append(checks, statecheck.ExpectKnownValue(
			resourceName,
			tfjsonpath.New("content"),
			knownvalue.StringExact(expectedOutput(tc.expectedOutput)),
		))

		steps = append(steps, resource.TestStep{
			Config:            wrapInput(tc.input),
			ConfigStateChecks: checks,
		})
	}

	return resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	}
}

func TestSetHostnameModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Create with hostname and fqdn",
			input: `
hostname = "myspace"
fqdn = "myspace.domain.com"
			`,
			expectedValues: map[string]string{
				"hostname": "myspace",
				"fqdn":     "myspace.domain.com",
			},
			expectedOutput: "hostname: myspace\nfqdn: myspace.domain.com\n",
		},
		{
			name: "Update with hostname only",
			input: `
hostname = "another-one"
			`,
			expectedValues: map[string]string{
				"hostname": "another-one",
				"fqdn":     "null",
			},
			expectedOutput: "hostname: another-one\n",
		},
		{
			name: "Other configurable values",
			input: `
prefer_fqdn_over_hostname= true
preserve_hostname= true
create_hostname_file= false
			`,
			expectedValues: map[string]string{
				"prefer_fqdn_over_hostname": "true",
				"preserve_hostname":         "true",
				"create_hostname_file":      "false",
				"hostname":                  "null",
			},
			expectedOutput: "prefer_fqdn_over_hostname: true\npreserve_hostname: true\ncreate_hostname_file: false\n",
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestUpdateEtcHostsModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Don't manage hosts",
			input: `
manage_etc_hosts = false
			`,
			expectedValues: map[string]string{
				"manage_etc_hosts": "false",
			},
			expectedOutput: "manage_etc_hosts: false\n",
		},
		{
			name: "Manage hosts",
			input: `
manage_etc_hosts = true
			`,
			expectedValues: map[string]string{
				"manage_etc_hosts": "true",
			},
			expectedOutput: "manage_etc_hosts: true\n",
		},
		{
			name: "Delegate to localhost",
			input: `
manage_etc_hosts_localhost = true 
			`,
			expectedValues: map[string]string{
				"manage_etc_hosts":           "null",
				"manage_etc_hosts_localhost": "true",
			},
			expectedOutput: "manage_etc_hosts: localhost\n",
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestLocaleModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic locale",
			input: `
locale = "en_GB"
			`,
			expectedValues: map[string]string{
				"locale": "en_GB",
			},
			expectedOutput: "locale: en_GB\n",
		},
		{
			name: "Localse config file",
			input: `
locale_configfile = "/tmp/locale"
			`,
			expectedValues: map[string]string{
				"locale": "null",
			},
			expectedOutput: "locale_configfile: /tmp/locale\n",
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestTimezoneModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic timezone",
			input: `
timezone = "Asia/Tokyo"
			`,
			expectedValues: map[string]string{
				"timezone": "Asia/Tokyo",
			},
			expectedOutput: "timezone: Asia/Tokyo\n",
		},
		{
			name: "Emptying",
			input: `
fqdn = "domain.com"
			`,
			expectedValues: map[string]string{
				"timezone": "null",
				"fqdn":     "domain.com",
			},
			expectedOutput: "fqdn: domain.com\n",
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestRunCMDModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic commands",
			input: `
runcmd = [ "cat /etc/hosts" ]
			`,
			expectedValues: map[string]string{
				"runcmd.0": "cat /etc/hosts",
			},
			expectedOutput: `
runcmd:
    - cat /etc/hosts
			`,
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestSSHModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic Auth keys",
			input: `
ssh_authorized_keys = [ "ssh key" ]
			`,
			expectedValues: map[string]string{
				"ssh_authorized_keys.0": "ssh key",
			},
			expectedOutput: `
ssh_authorized_keys:
    - ssh key 
			`,
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestSetPasswordModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Password auth disabled",
			input: `
ssh_pwauth = false 
			`,
			expectedValues: map[string]string{
				"ssh_pwauth": "false",
			},
			expectedOutput: "ssh_pwauth: false\n",
		},
		{
			name: "Password auth",
			input: `
ssh_pwauth = true 
			`,
			expectedValues: map[string]string{
				"ssh_pwauth": "true",
			},
			expectedOutput: "ssh_pwauth: true\n",
		},
		{
			name: "chpasswd block",
			input: `
chpasswd {
			expire = false

      users {
        name = "ansible"
        type = "RANDOM"
      }

      users {
        name = "docker"
        password = "12345678"
        type = "text"
      }
}
			`,
			expectedValues: map[string]string{
				"chpasswd.users.0.name":     "ansible",
				"chpasswd.users.0.type":     "RANDOM",
				"chpasswd.users.1.name":     "docker",
				"chpasswd.users.1.password": "12345678",
				"chpasswd.users.1.type":     "text",
			},
			expectedOutput: `
chpasswd:
    users:
        - name: ansible
          type: RANDOM
        - name: docker
          password: "12345678"
          type: text
    expire: false
      `,
		},
		{
			name: "remove user",
			input: `
chpasswd {
      users {
        name = "ansible"
        type = "RANDOM"
      }
}
			`,
			expectedValues: map[string]string{
				"chpasswd.users.0.name": "ansible",
				"chpasswd.users.0.type": "RANDOM",
			},
			expectedOutput: `
chpasswd:
    users:
        - name: ansible
          type: RANDOM
      `,
		},
		{
			name: "only expiry",
			input: `
chpasswd {
    expire = false
}
			`,
			expectedValues: map[string]string{
				"chpasswd.expire": "false",
			},
			expectedOutput: `
chpasswd:
    expire: false
      `,
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}
