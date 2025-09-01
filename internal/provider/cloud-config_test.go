package provider

import (
	"fmt"
	"math/big"
	"strconv"
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
				if utils.IsNumeric(val) {
					nVal, _ := strconv.ParseFloat(val, 64)
					checks = append(checks, statecheck.ExpectKnownValue(
						resourceName,
						path,
						knownvalue.NumberExact(big.NewFloat(nVal)),
					))
				} else {
					checks = append(checks, statecheck.ExpectKnownValue(
						resourceName,
						path,
						knownvalue.StringExact(val),
					))
				}
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

func TestZypperModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
zypper {
      config = {
        "download.use_deltarpm": "true", 
        "reposdir": "/etc/zypp/repos.dir", 
        "servicesdir": "/etc/zypp/services.d"
      }

      repos {
        baseurl = "http://dl.opensuse.org/dist/leap/v/update"
        id = "opensuse-oss-update"
      }

      repos {
        baseurl = "http://dl.opensuse.org/dist/leap/v/repo/oss/"
        id = "opensuse-oss"
      }
}
      `,
			expectedValues: map[string]string{
				"zypper.config.reposdir":    "/etc/zypp/repos.dir",
				"zypper.config.servicesdir": "/etc/zypp/services.d",
				"zypper.repos.0.id":         "opensuse-oss-update",
				"zypper.repos.0.baseurl":    "http://dl.opensuse.org/dist/leap/v/update",
				"zypper.repos.1.id":         "opensuse-oss",
				"zypper.repos.1.baseurl":    "http://dl.opensuse.org/dist/leap/v/repo/oss/",
			},
			expectedOutput: `
zypper:
    repos:
        - id: opensuse-oss-update
          baseurl: http://dl.opensuse.org/dist/leap/v/update
        - id: opensuse-oss
          baseurl: http://dl.opensuse.org/dist/leap/v/repo/oss/
    config:
        download.use_deltarpm: "true"
        reposdir: /etc/zypp/repos.dir
        servicesdir: /etc/zypp/services.d
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestWireguardModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
wireguard {
      readinessprobe = ["systemctl restart service", "curl https://webhook.endpoint/example", "nc -zv some-service-fqdn 443"] 

      interfaces {
         name = "wg01"
         config_path = "/etc/wireguard/wg0.conf"
        content = "[Interface]"
      }
}
      `,
			expectedValues: map[string]string{
				"wireguard.readinessprobe.0":         "systemctl restart service",
				"wireguard.readinessprobe.1":         "curl https://webhook.endpoint/example",
				"wireguard.readinessprobe.2":         "nc -zv some-service-fqdn 443",
				"wireguard.interfaces.0.name":        "wg01",
				"wireguard.interfaces.0.config_path": "/etc/wireguard/wg0.conf",
				"wireguard.interfaces.0.content":     "[Interface]",
			},
			expectedOutput: `
wireguard:
    interfaces:
        - name: wg01
          config_path: /etc/wireguard/wg0.conf
          content: '[Interface]'
    readinessprobe:
        - systemctl restart service
        - curl https://webhook.endpoint/example
        - nc -zv some-service-fqdn 443
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestSeedRandomModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
random_seed {
      file = "/dev/urandom"
      data = "my random string"
      encoding = "raw"
      command_required = true
      command = ["sh", "-c", "dd if=/dev/urandom of=$RANDOM_SEED_FILE"]
}
      `,
			expectedValues: map[string]string{
				"random_seed.file":             "/dev/urandom",
				"random_seed.data":             "my random string",
				"random_seed.encoding":         "raw",
				"random_seed.command_required": "true",
				"random_seed.command.0":        "sh",
				"random_seed.command.1":        "-c",
				"random_seed.command.2":        "dd if=/dev/urandom of=$RANDOM_SEED_FILE",
			},
			expectedOutput: `
random_seed:
    file: /dev/urandom
    data: my random string
    encoding: raw
    command:
        - sh
        - -c
        - dd if=/dev/urandom of=$RANDOM_SEED_FILE
    command_required: true
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestRPIModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
rpi {
      enable_rpi_connect = true

      interfaces {
        spi = true
        i2c = true
        ssh = true
        onewire = true
        remote_gpio = true

        serial {
          console = true
          hardware = true
        }

      }
}
      `,
			expectedValues: map[string]string{
				"rpi.enable_rpi_connect":         "true",
				"rpi.interfaces.spi":             "true",
				"rpi.interfaces.i2c":             "true",
				"rpi.interfaces.ssh":             "true",
				"rpi.interfaces.onewire":         "true",
				"rpi.interfaces.remote_gpio":     "true",
				"rpi.interfaces.serial.console":  "true",
				"rpi.interfaces.serial.hardware": "true",
			},
			expectedOutput: `
rpi:
    interfaces:
        spi: true
        i2c: true
        ssh: true
        serial:
            console: true
            hardware: true
        onewire: true
        remote_gpio: true
    enable_rpi_connect: true
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestNTPModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
ntp {
  enabled= true
  ntp_client= "myntpclient"
  pools= ["0.int.pool.ntp.org", "1.int.pool.ntp.org", "ntp.myorg.org"]
  servers= [
   "ntp.server.local",
   "ntp.ubuntu.com",
   "192.168.23.2"
  ]
  allow=  [ "192.168.23.0/32" ]
  peers=  [ "km001",  "km002" ]

      config {
        confpath= "/etc/myntpclient/myntpclient.conf"
        check_exe= "myntpclientd"
        packages=  [ "myntpclient" ]
        service_name= "myntpclient"
        template= "template"
      }
}
      `,
			expectedValues: map[string]string{
				"ntp.enabled":    "true",
				"ntp.ntp_client": "myntpclient",
				"ntp.pools.0":    "0.int.pool.ntp.org",
				"ntp.pools.1":    "1.int.pool.ntp.org",
				"ntp.pools.2":    "ntp.myorg.org",
				"ntp.servers.0":  "ntp.server.local",
				"ntp.servers.1":  "ntp.ubuntu.com",
				"ntp.servers.2":  "192.168.23.2",
				"ntp.allow.0":    "192.168.23.0/32",
				"ntp.peers.0":    "km001",
				"ntp.peers.1":    "km002",
			},
			expectedOutput: `
ntp:
    pools:
        - 0.int.pool.ntp.org
        - 1.int.pool.ntp.org
        - ntp.myorg.org
    servers:
        - ntp.server.local
        - ntp.ubuntu.com
        - 192.168.23.2
    peers:
        - km001
        - km002
    allow:
        - 192.168.23.0/32
    ntp_client: myntpclient
    enabled: true
    config:
        confpath: /etc/myntpclient/myntpclient.conf
        check_exe: myntpclientd
        packages:
            - myntpclient
        service_name: myntpclient
        template: template

		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestLandscapeModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
landscape {
      client {
          url= "https://landscape.canonical.com/message-system"
          ping_url= "http://landscape.canonical.com/ping"
          data_path= "/var/lib/landscape/client"
          http_proxy= "http://my.proxy.com/foobar"
          https_proxy= "https://my.proxy.com/foobar"
          tags= "server,cloud"
          computer_title= "footitle"
          registration_key= "fookey"
          account_name= "fooaccount"
      }
}
      `,
			expectedValues: map[string]string{
				"landscape.client.url":              "https://landscape.canonical.com/message-system",
				"landscape.client.ping_url":         "http://landscape.canonical.com/ping",
				"landscape.client.data_path":        "/var/lib/landscape/client",
				"landscape.client.http_proxy":       "http://my.proxy.com/foobar",
				"landscape.client.https_proxy":      "https://my.proxy.com/foobar",
				"landscape.client.tags":             "server,cloud",
				"landscape.client.computer_title":   "footitle",
				"landscape.client.registration_key": "fookey",
				"landscape.client.account_name":     "fooaccount",
			},
			expectedOutput: `
landscape:
    client:
        url: https://landscape.canonical.com/message-system
        ping_url: http://landscape.canonical.com/ping
        data_path: /var/lib/landscape/client
        computer_title: footitle
        account_name: fooaccount
        registration_key: fookey
        tags: server,cloud
        http_proxy: http://my.proxy.com/foobar
        https_proxy: https://my.proxy.com/foobar
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestBootCMDModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic commands",
			input: `
bootcmd = [ "cat /etc/hosts" ]
			`,
			expectedValues: map[string]string{
				"bootcmd.0": "cat /etc/hosts",
			},
			expectedOutput: `
bootcmd:
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
        password = "p12345678"
        type = "text"
      }
}
			`,
			expectedValues: map[string]string{
				"chpasswd.users.0.name":     "ansible",
				"chpasswd.users.0.type":     "RANDOM",
				"chpasswd.users.1.name":     "docker",
				"chpasswd.users.1.password": "p12345678",
				"chpasswd.users.1.type":     "text",
			},
			expectedOutput: `
chpasswd:
    users:
        - name: ansible
          type: RANDOM
        - name: docker
          password: p12345678
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
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestPkgUpdateUpgradeModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic Auth keys",
			input: `
      package_update = true
      package_upgrade = true
      package_reboot_if_required = false
      packages = [
        "qemu-guest-agent",
      "ufw"
      ]
			`,
			expectedValues: map[string]string{
				"package_update":             "true",
				"package_upgrade":            "true",
				"packages.0":                 "qemu-guest-agent",
				"packages.1":                 "ufw",
				"package_reboot_if_required": "false",
			},
			expectedOutput: `
package_update: true
package_upgrade: true
packages:
    - qemu-guest-agent
    - ufw
      `,
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestUsersAndGroupsModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "User configuration",
			input: `
user {
  name = "myname"
  doas = ["do this", "do that"]
  expiredate = "Jan 2025"
  gecos = "This use is jsut for testing"
  homedir = "/home/myname"
  inactive = "90 days" # deactivate in 90 days
  lock_passwd = false
  no_create_home = false
  no_log_init = true
  no_user_group = true
  passwd = "mypwd"
  hashed_passwd = "hashed pwd"
  plain_text_passwd = "mypwd"
  create_groups = true
  primary_group = "default"
  selinux_user = "myname"
  shell = "/bin/bash"
  snapuser = "myname"
  ssh_authorized_keys = ["mykey"]
  ssh_import_id = ["import key"]
  system = true
  sudo = ["idk what is expected here, documentation is silent"]
  uid = 1001
  groups = ["docker"]
}
			`,
			expectedValues: map[string]string{
				"user.name":                  "myname",
				"user.doas.0":                "do this",
				"user.doas.1":                "do that",
				"user.expiredate":            "Jan 2025",
				"user.gecos":                 "This use is jsut for testing",
				"user.homedir":               "/home/myname",
				"user.inactive":              "90 days",
				"user.lock_passwd":           "false",
				"user.no_create_home":        "false",
				"user.no_log_init":           "true",
				"user.no_user_group":         "true",
				"user.passwd":                "mypwd",
				"user.hashed_passwd":         "hashed pwd",
				"user.plain_text_passwd":     "mypwd",
				"user.create_groups":         "true",
				"user.primary_group":         "default",
				"user.selinux_user":          "myname",
				"user.shell":                 "/bin/bash",
				"user.snapuser":              "myname",
				"user.ssh_authorized_keys.0": "mykey",
				"user.ssh_import_id.0":       "import key",
				"user.system":                "true",
				"user.sudo.0":                "idk what is expected here, documentation is silent",
				"user.uid":                   "1001",
				"user.groups.0":              "docker",
			},
			expectedOutput: `
user:
    name: myname
    doas:
        - do this
        - do that
    expiredate: Jan 2025
    gecos: This use is jsut for testing
    homedir: /home/myname
    inactive: 90 days
    lock_passwd: false
    no_log_init: true
    no_user_group: true
    passwd: mypwd
    hashed_passwd: hashed pwd
    plain_text_passwd: mypwd
    create_groups: true
    primary_group: default
    selinux_user: myname
    shell: /bin/bash
    snapuser: myname
    ssh_authorized_keys:
        - mykey
    ssh_import_id:
        - import key
    system: true
    uid: 1001
    sudo:
        - idk what is expected here, documentation is silent
    groups:
        - docker 
      `,
		},
		{
			name: "UserS configuration",
			input: `
users {
  name = "myname"
  doas = ["do this", "do that"]
  expiredate = "Jan 2025"
  gecos = "This use is jsut for testing"
  homedir = "/home/myname"
  inactive = "90 days" # deactivate in 90 days
  lock_passwd = false
  no_create_home = false
  no_log_init = true
  no_user_group = true
  passwd = "mypwd"
  hashed_passwd = "hashed pwd"
  plain_text_passwd = "mypwd"
  create_groups = true
  primary_group = "default"
  selinux_user = "myname"
  shell = "/bin/bash"
  snapuser = "myname"
  ssh_authorized_keys = ["mykey"]
  ssh_import_id = ["import key"]
  system = true
  sudo = ["idk what is expected here, documentation is silent"]
  uid = 1001
}
users {
  name = "myname2"
  doas = ["do this", "do that"]
  expiredate = "Jan 2025"
  gecos = "This use is jsut for testing"
  homedir = "/home/myname"
  inactive = "90 days" # deactivate in 90 days
  lock_passwd = true
  no_create_home = true
  no_log_init = false
  no_user_group = false
  passwd = "mypwd"
  hashed_passwd = "hashed pwd"
  plain_text_passwd = "mypwd"
  create_groups = true
  primary_group = "default"
  selinux_user = "myname"
  shell = "/bin/bash"
  snapuser = "myname"
  ssh_redirect_user = true
  system = true
  sudo = ["idk what is expected here, documentation is silent"]
  uid = 1001
}
			`,
			expectedValues: map[string]string{
				"users.0.name":                  "myname",
				"users.0.doas.0":                "do this",
				"users.0.doas.1":                "do that",
				"users.0.expiredate":            "Jan 2025",
				"users.0.gecos":                 "This use is jsut for testing",
				"users.0.homedir":               "/home/myname",
				"users.0.inactive":              "90 days",
				"users.0.lock_passwd":           "false",
				"users.0.no_create_home":        "false",
				"users.0.no_log_init":           "true",
				"users.0.no_user_group":         "true",
				"users.0.passwd":                "mypwd",
				"users.0.hashed_passwd":         "hashed pwd",
				"users.0.plain_text_passwd":     "mypwd",
				"users.0.create_groups":         "true",
				"users.0.primary_group":         "default",
				"users.0.selinux_user":          "myname",
				"users.0.shell":                 "/bin/bash",
				"users.0.snapuser":              "myname",
				"users.0.ssh_authorized_keys.0": "mykey",
				"users.0.ssh_import_id.0":       "import key",
				"users.0.system":                "true",
				"users.0.sudo.0":                "idk what is expected here, documentation is silent",
				"users.0.uid":                   "1001",

				"users.1.name":              "myname2",
				"users.1.doas.0":            "do this",
				"users.1.doas.1":            "do that",
				"users.1.expiredate":        "Jan 2025",
				"users.1.gecos":             "This use is jsut for testing",
				"users.1.homedir":           "/home/myname",
				"users.1.inactive":          "90 days",
				"users.1.lock_passwd":       "true",
				"users.1.no_create_home":    "true",
				"users.1.no_log_init":       "false",
				"users.1.no_user_group":     "false",
				"users.1.passwd":            "mypwd",
				"users.1.hashed_passwd":     "hashed pwd",
				"users.1.plain_text_passwd": "mypwd",
				"users.1.create_groups":     "true",
				"users.1.primary_group":     "default",
				"users.1.selinux_user":      "myname",
				"users.1.shell":             "/bin/bash",
				"users.1.snapuser":          "myname",
				"users.1.ssh_redirect_user": "true",
				"users.1.system":            "true",
				"users.1.sudo.0":            "idk what is expected here, documentation is silent",
				"users.1.uid":               "1001",
			},
			expectedOutput: `
users:
    - name: myname
      doas:
        - do this
        - do that
      expiredate: Jan 2025
      gecos: This use is jsut for testing
      homedir: /home/myname
      inactive: 90 days
      lock_passwd: false
      no_log_init: true
      no_user_group: true
      passwd: mypwd
      hashed_passwd: hashed pwd
      plain_text_passwd: mypwd
      create_groups: true
      primary_group: default
      selinux_user: myname
      shell: /bin/bash
      snapuser: myname
      ssh_authorized_keys:
        - mykey
      ssh_import_id:
        - import key
      system: true
      uid: 1001
      sudo:
        - idk what is expected here, documentation is silent
    - name: myname2
      doas:
        - do this
        - do that
      expiredate: Jan 2025
      gecos: This use is jsut for testing
      homedir: /home/myname
      inactive: 90 days
      no_create_home: true
      passwd: mypwd
      hashed_passwd: hashed pwd
      plain_text_passwd: mypwd
      create_groups: true
      primary_group: default
      selinux_user: myname
      shell: /bin/bash
      snapuser: myname
      ssh_redirect_user: true
      system: true
      uid: 1001
      sudo:
        - idk what is expected here, documentation is silent

      `,
		},
		{
			name: "Almost empty",
			input: `
user {
      name = "ansible"
}

users {
  name = "docker"
}
			`,
			expectedValues: map[string]string{
				"user.name":    "ansible",
				"users.0.name": "docker",
			},
			expectedOutput: `
user:
    name: ansible
users:
    - name: docker
      `,
		},
		{
			name: "Groups",
			input: `
      groups = ["docker"]
			`,
			expectedValues: map[string]string{
				"groups.0": "docker",
			},
			expectedOutput: `
groups:
    - docker
      `,
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestDisableEC2InstanceMetadataModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic case",
			input: `
disable_ec2_metadata = true
			`,
			expectedValues: map[string]string{
				"disable_ec2_metadata": "true",
			},
			expectedOutput: `
disable_ec2_metadata: true
			`,
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestApkConfigureModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic case",
			input: `
apk_repos {
  preserve_repositories = true
}
			`,
			expectedValues: map[string]string{
				"apk_repos.preserve_repositories": "true",
			},
			expectedOutput: `
apk_repos:
    preserve_repositories: true
			`,
		},
		{
			name: "Nested case",
			input: `
apk_repos {
  local_repo_base_url = "https://my-local-server/local-alpine"

  alpine_repo {
        base_url= "https://some-alpine-mirror/alpine"
        community_enabled= true
        testing_enabled= true
        version= "edge"
      }
}
			`,
			expectedValues: map[string]string{
				"apk_repos.local_repo_base_url":           "https://my-local-server/local-alpine",
				"apk_repos.alpine_repo.base_url":          "https://some-alpine-mirror/alpine",
				"apk_repos.alpine_repo.version":           "edge",
				"apk_repos.alpine_repo.community_enabled": "true",
				"apk_repos.alpine_repo.testing_enabled":   "true",
			},
			expectedOutput: `
apk_repos:
    local_repo_base_url: https://my-local-server/local-alpine
    alpine_repo:
        community_enabled: true
        testing_enabled: true
        base_url: https://some-alpine-mirror/alpine
        version: edge
      `,
		},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestAptPipeliningModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic case",
			input: `
apt_pipelining {
      os = true
}
			`,
			expectedValues: map[string]string{
				"apt_pipelining.os": "true",
			},
			expectedOutput: `
apt_pipelining: os
			`,
		},
		{
			name: "As a number",
			input: `
apt_pipelining {
      depth = 14
}
			`,
			expectedValues: map[string]string{
				"apt_pipelining.depth": "14",
			},
			expectedOutput: `
apt_pipelining: 14
			`,
		},
		{
			name: "disable",
			input: `
apt_pipelining {
      disable = true 
}
			`,
			expectedValues: map[string]string{
				"apt_pipelining.disable": "true",
			},
			expectedOutput: `
apt_pipelining: false
			`,
		},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestByobuModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
byobu_by_default = "system" 
			`,
			expectedValues: map[string]string{
				"byobu_by_default": "system",
			},
			expectedOutput: `
byobu_by_default: system
			`,
		},
		{
			name: "Basic",
			input: `
byobu_by_default = "disable-system" 
			`,
			expectedValues: map[string]string{
				"byobu_by_default": "disable-system",
			},
			expectedOutput: `
byobu_by_default: disable-system
			`,
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestCACertificatesModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
ca_certs {
  remove_defaults = true
  trusted = ["single_line_cert"]
}
			`,
			expectedValues: map[string]string{
				"ca_certs.remove_defaults": "true",
				"ca_certs.trusted.0":       "single_line_cert",
			},
			expectedOutput: `
ca_certs:
    remove_defaults: true
    trusted:
        - single_line_cert
			`,
		},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestFanModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
fan {
  config = "My fan config"
  config_path = "/etc/config/fan"
}
			`,
			expectedValues: map[string]string{
				"fan.config":      "My fan config",
				"fan.config_path": "/etc/config/fan",
			},
			expectedOutput: `
fan:
    config: My fan config
    config_path: /etc/config/fan
			`,
		},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestFinalMessageModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
final_message = "Let's get it cracking"
			`,
			expectedValues: map[string]string{
				"final_message": "Let's get it cracking",
			},
			expectedOutput: "final_message: Let's get it cracking\n",
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestGrowpartModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
growpart {
  mode = "auto"
  ignore_growroot_disabled = true
  devices = ["/"]
}
			`,
			expectedValues: map[string]string{
				"growpart.mode":                     "auto",
				"growpart.ignore_growroot_disabled": "true",
				"growpart.devices.0":                "/",
			},
			expectedOutput: `
growpart:
    mode: auto
    devices:
        - /
    ignore_growroot_disabled: true
			`,
		},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestGRUBDpkgModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
grub_dpkg {
  enabled = true
  grub_pc_install_devices_empty = true
  grub_pc_install_devices = "/boot"
  grub_efi_install_devices = "/boot/efi"
}
			`,
			expectedValues: map[string]string{
				"grub_dpkg.enabled":                       "true",
				"grub_dpkg.grub_pc_install_devices_empty": "true",
				"grub_dpkg.grub_pc_install_devices":       "/boot",
				"grub_dpkg.grub_efi_install_devices":      "/boot/efi",
			},
			expectedOutput: `
grub_dpkg:
    enabled: true
    grub-pc/install_devices: /boot
    grub-pc/install_devices_empty: true
    grub-efi/install_devices: /boot/efi
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestInstallHotplugModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
updates {
    network {
      when = ["boot"]
      }
}
			`,
			expectedValues: map[string]string{
				"updates.network.when.0": "boot",
			},
			expectedOutput: `
updates:
    network:
        when:
            - boot
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestKeyboardModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
keyboard {
  layout = "en_GB"
  model = "pc105"
  variant = "some"
  options = "some"
}
			`,
			expectedValues: map[string]string{
				"keyboard.layout":  "en_GB",
				"keyboard.model":   "pc105",
				"keyboard.variant": "some",
				"keyboard.options": "some",
			},
			expectedOutput: `
keyboard:
    layout: en_GB
    model: pc105
    variant: some
    options: some
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestKeysToConsoleModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
ssh {
  emit_keys_to_console = false
}

ssh_key_console_blacklist = ["rsa"]
ssh_fp_console_blacklist = ["E25451E0221B5773DEBFF178ECDACB160995AA89"]
			`,
			expectedValues: map[string]string{
				"ssh.emit_keys_to_console":    "false",
				"ssh_key_console_blacklist.0": "rsa",
				"ssh_fp_console_blacklist.0":  "E25451E0221B5773DEBFF178ECDACB160995AA89",
			},
			expectedOutput: `
ssh:
    emit_keys_to_console: false
ssh_key_console_blacklist:
    - rsa
ssh_fp_console_blacklist:
    - E25451E0221B5773DEBFF178ECDACB160995AA89
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestResizefsModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
resize_rootfs = false
resize_rootfs_no_block = true
			`,
			expectedValues: map[string]string{
				"resize_rootfs":          "false",
				"resize_rootfs_no_block": "true",
			},
			expectedOutput: `
resize_rootfs: false
    `},
		{
			name: "On boot",
			input: `
resize_rootfs_no_block = true
			`,
			expectedValues: map[string]string{
				"resize_rootfs_no_block": "true",
			},
			expectedOutput: `
resize_rootfs: noblock
			`,
		},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestSaltMinionModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
salt_minion {
  config_dir = "/etc/salt"
  pkg_name = "salt-minion"
  pki_dir = "/etc/salt/pki/minion"
  private_key = "------BEGIN PRIVATE KEY------"
  public_key = "------BEGIN PUBLIC KEY-------"
  service_name = "salt-minion"
}
      `,
			expectedValues: map[string]string{
				"salt_minion.config_dir":   "/etc/salt",
				"salt_minion.pkg_name":     "salt-minion",
				"salt_minion.pki_dir":      "/etc/salt/pki/minion",
				"salt_minion.private_key":  "------BEGIN PRIVATE KEY------",
				"salt_minion.public_key":   "------BEGIN PUBLIC KEY-------",
				"salt_minion.service_name": "salt-minion",
			},
			expectedOutput: `
salt_minion:
    pkg_name: salt-minion
    service_name: salt-minion
    config_dir: /etc/salt
    public_key: '------BEGIN PUBLIC KEY-------'
    private_key: '------BEGIN PRIVATE KEY------'
    pki_dir: /etc/salt/pki/minion
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestUbuntuAutoinstallModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
autoinstall {
  version = 19
}
      `,
			expectedValues: map[string]string{
				"autoinstall.version": "19",
			},
			expectedOutput: `
autoinstall:
    version: 19 
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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

func TestPowerStateChangeModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
power_state {
      delay = 19 
      mode = "halt" 
      message = "Charming message"
      timeout = 100
      condition_cmd = "cat /etc/log"
}
      `,
			expectedValues: map[string]string{
				"power_state.delay":         "19",
				"power_state.mode":          "halt",
				"power_state.message":       "Charming message",
				"power_state.timeout":       "100",
				"power_state.condition_cmd": "cat /etc/log",
			},
			expectedOutput: `
power_state:
    delay: 19
    mode: halt
    message: Charming message
    timeout: 100
    condition: cat /etc/log

		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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
		{
			name: "Complex case",
			input: `
power_state {
      no_delay = true
      mode = "halt" 
      message = "Charming message"
      timeout = 100
      condition = true
}
      `,
			expectedValues: map[string]string{
				"power_state.no_delay":  "true",
				"power_state.mode":      "halt",
				"power_state.message":   "Charming message",
				"power_state.timeout":   "100",
				"power_state.condition": "true",
			},
			expectedOutput: `
power_state:
    delay: now
    mode: halt
    message: Charming message
    timeout: 100
    condition: true 

		`},
	}

	resource.Test(t, assembleTestCase(testCases, t))
}

func TestPhoneHomeModule(t *testing.T) {
	testCases := []testCase{
		{
			name: "Basic",
			input: `
phone_home {
  url = "http://example.com/{{ v1.instance_id }}/"
  tries = 5
  post = ["pub_key_rsa", "pub_key_ecdsa", "pub_key_ed25519", "instance_id", "hostname", "fqdn"]
}
      `,
			expectedValues: map[string]string{
				"phone_home.url":    "http://example.com/{{ v1.instance_id }}/",
				"phone_home.tries":  "5",
				"phone_home.post.0": "pub_key_rsa",
				"phone_home.post.1": "pub_key_ecdsa",
				"phone_home.post.2": "pub_key_ed25519",
				"phone_home.post.3": "instance_id",
				"phone_home.post.4": "hostname",
				"phone_home.post.5": "fqdn",
			},
			expectedOutput: `
phone_home:
    url: http://example.com/{{ v1.instance_id }}/
    tries: 5
    post:
        - pub_key_rsa
        - pub_key_ecdsa
        - pub_key_ed25519
        - instance_id
        - hostname
        - fqdn
		`},
		{
			name: "Fail in older versions because `chapasswd` block needs to be deleted",
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
