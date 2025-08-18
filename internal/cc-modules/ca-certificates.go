package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type CACerts struct {
	RemoveDefaults types.Bool `tfsdk:"remove_defaults"`
	Trusted        types.List `tfsdk:"trusted"`
}

type CACertsOutput struct {
	RemoveDefaults bool      `yaml:"remove_defaults,omitempty"`
	Trusted        *[]string `yaml:"trusted,omitempty"`
}

type CACertificatesModel struct {
	CACerts *CACerts `tfsdk:"ca_certs"`
}

type CACertificatesOutputModel struct {
	CACerts *CACertsOutput `yaml:"ca_certs,omitempty"`
}

// CACertificatesBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#ca-certificates
func CACertificatesBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"ca_certs": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("ca_certs")),
				},
				MarkdownDescription: `
This module adds CA certificates to the system’s CA store and updates any related files using the appropriate OS-specific utility. The default CA certificates can be disabled/deleted from use by the system with the configuration option remove_defaults.

Certificates must be specified using valid YAML. To specify a multi-line certificate, the YAML multi-line list syntax must be used.

Alpine Linux requires the ca-certificates package to be installed in order to provide the update-ca-certificates command.
        `,
				Attributes: map[string]schema.Attribute{
					"remove_defaults": schema.BoolAttribute{
						MarkdownDescription: "Remove default CA certificates if true. Default: `false`.",
						Optional:            true,
					},
					"trusted": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "List of trusted CA certificates to add.",
						Optional:            true,
					},
				},
			},
		},
	}
}
