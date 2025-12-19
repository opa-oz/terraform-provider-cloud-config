package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type Spacewalk struct {
	Server        types.String `tfsdk:"server"`
	Proxy         types.String `tfsdk:"proxy"`
	ActivationKey types.String `tfsdk:"activation_key"`
}

type SpacewalkOutput struct {
	Server        string `yaml:"server,omitempty"`
	Proxy         string `yaml:"proxy,omitempty"`
	ActivationKey string `yaml:"activation_key,omitempty"`
}

type SpacewalkModel struct {
	Spacewalk *Spacewalk `tfsdk:"spacewalk"`
}

type SpacewalkOutputModel struct {
	Spacewalk *SpacewalkOutput `yaml:"spacewalk,omitempty"`
}

// SpacewalkBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#spacewalk
func SpacewalkBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"spacewalk": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("spacewalk")),
				},
				MarkdownDescription: `
This module installs Spacewalk and applies basic configuration. 
If the Spacewalk config key is present, Spacewalk will be installed. The server to connect to after installation must be provided in the server in Spacewalk configuration. A proxy to connect through and an activation key may optionally be specified.
				`,
				Attributes: map[string]schema.Attribute{
					"server": schema.StringAttribute{
						MarkdownDescription: "The Spacewalk server to use.",
						Optional:            true,
					},
					"proxy": schema.StringAttribute{
						MarkdownDescription: "The proxy to use when connecting to Spacewalk.",
						Optional:            true,
					},
					"activation_key": schema.StringAttribute{
						MarkdownDescription: "The activation key to use when registering with Spacewalk.",
						Optional:            true,
					},
				},
			},
		},
	}
}
