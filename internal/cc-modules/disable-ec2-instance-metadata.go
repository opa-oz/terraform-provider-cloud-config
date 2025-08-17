package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DisableEC2InstanceMetadataModel struct {
	DisableEC2Metadata types.Bool `tfsdk:"disable_ec2_metadata"`
}

type DisableEC2InstanceMetadataOutputModel struct {
	DisableEC2Metadata bool `yaml:"disable_ec2_metadata,omitempty"`
}

// DisableEC2InstanceMetadata
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#disable-ec2-instance-metadata-service
func DisableEC2InstanceMetadata() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"disable_ec2_metadata": schema.BoolAttribute{
				MarkdownDescription: "Set `true` to disable IPv4 routes to EC2 metadata. Default: `false`.",
				Optional:            true,
			},
		},
	}
}
