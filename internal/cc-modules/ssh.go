package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SSHModel struct {
	SSHAuthorizedKeys types.List `tfsdk:"ssh_authorized_keys"`
}

type SSHOutputModel struct {
	SSHAuthorizedKeys *[]string `yaml:"ssh_authorized_keys,omitempty"`
}

// SSH
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#ssh
func SSH() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"ssh_authorized_keys": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "The SSH public keys to add `.ssh/authorized_keys` in the default user’s home directory.",
				Optional:            true,
			},
		},
	}
}
