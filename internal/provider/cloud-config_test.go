package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

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
						"cloud-config.test",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("one"),
					),
					statecheck.ExpectKnownValue(
						"cloud-config.test",
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
						"cloud-config.test",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact("two"),
					),
					statecheck.ExpectKnownValue(
						"cloud-config.test",
						tfjsonpath.New("fqdn"),
						knownvalue.StringExact("two.lan"),
					),
					statecheck.ExpectKnownValue(
						"cloud-config.test",
						tfjsonpath.New("content"),
						knownvalue.StringExact(strings.TrimSpace(`#cloud-config
hostname: two
fqdn: two.lan
prefer_fqdn_over_hostname: true
preserve_hostname: true
create_hostname_file: false
locale: en_two
locale_configfile: /etc/locale
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
}
`, configurableAttribute)
}
