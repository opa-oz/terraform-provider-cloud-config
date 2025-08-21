package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type Network struct {
	When types.List `tfsdk:"when"`
}

type NetworkOutput struct {
	When *[]string `yaml:"when,omitempty"`
}

type Updates struct {
	Network *Network `tfsdk:"network"`
}

type UpdatesOutput struct {
	Network *NetworkOutput `yaml:"network,omitempty"`
}

type InstallHotplugModel struct {
	Updates *Updates `tfsdk:"updates"`
}

type InstallHotplugOutputModel struct {
	Updates *UpdatesOutput `yaml:"updates,omitempty"`
}

// InstallHotplugBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#install-hotplug
func InstallHotplugBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"updates": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("updates")),
				},
				MarkdownDescription: `
This module will install the udev rules to enable hotplug if supported by the datasource and enabled in the user-data. The udev rules will be installed as /etc/udev/rules.d/90-cloud-init-hook-hotplug.rules.

When hotplug is enabled, newly added network devices will be added to the system by cloud-init. After udev detects the event, cloud-init will refresh the instance metadata from the datasource, detect the device in the updated metadata, then apply the updated network configuration.

Udev rules are installed while cloud-init is running, which means that devices which are added during boot might not be configured. To work around this limitation, one can wait until cloud-init has completed before hotplugging devices.

Currently supported datasources: **Openstack, EC2**
        `,
				Blocks: map[string]schema.Block{
					"network": schema.SingleNestedBlock{
						PlanModifiers: []planmodifier.Object{
							utils.NullWhen(path.Root("network")),
						},
						Attributes: map[string]schema.Attribute{
							"when": schema.ListAttribute{
								ElementType:         types.StringType,
								MarkdownDescription: "array of boot-new-instance/boot-legacy/boot/hotplug",
								Optional:            true,
								Validators: []validator.List{
									listvalidator.ValueStringsAre(
										stringvalidator.OneOf(
											"boot",
											"hotplug",
											"boot-legacy",
											"boot-new-instance",
										),
									),
								},
							},
						},
					},
				},
			},
		},
	}
}
