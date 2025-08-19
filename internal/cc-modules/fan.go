package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type Fan struct {
	Config     types.String `tfsdk:"config"`
	ConfigPath types.String `tfsdk:"config_path"`
}

type FanOutput struct {
	Config     string `yaml:"config,omitempty"`
	ConfigPath string `yaml:"config_path,omitempty"`
}

type FanModel struct {
	Fan *Fan `tfsdk:"fan"`
}

type FanOutputModel struct {
	Fan *FanOutput `yaml:"fan,omitempty"`
}

// FanBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#fan
func FanBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"fan": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("fan")),
				},
				MarkdownDescription: `
This module installs, configures and starts the Ubuntu fan network system ([Read more about Ubuntu Fan](https://wiki.ubuntu.com/FanNetworking)).

If cloud-init sees a fan entry in cloud-config it will:

 - Write config_path with the contents of the config key
 - Install the package ubuntu-fan if it is not installed
 - Ensure the service is started (or restarted if was previously running)

Additionally, the ubuntu-fan package will be automatically installed if not present.
        `,
				Attributes: map[string]schema.Attribute{
					"config": schema.StringAttribute{
						MarkdownDescription: "The fan configuration to use as a single multi-line string.",
						Optional:            true,
					},
					"config_path": schema.StringAttribute{
						MarkdownDescription: "The path to write the fan configuration to. Default: `/etc/network/fan`.",
						Optional:            true,
					},
				},
			},
		},
	}
}
