package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UpdateEtcHostsModule struct {
	ManageEtcHosts          types.Bool `tfsdk:"manage_etc_hosts"`
	ManageEtcHostsLocalhost types.Bool `tfsdk:"manage_etc_hosts_localhost"`
}

type UpdateEtcHostsOutputModule struct {
	ManageEtcHosts any `yaml:"manage_etc_hosts,omitempty"`
}

// UpdateEtcHosts
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#update-etc-hosts
func UpdateEtcHosts() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"manage_etc_hosts": schema.BoolAttribute{
				MarkdownDescription: "Whether to manage `/etc/hosts` on the system. If true, render the hosts file using `/etc/cloud/templates/hosts.tmpl` replacing `$hostname` and `$fdqn`.",
				Optional:            true,
			},
			"manage_etc_hosts_localhost": schema.BoolAttribute{
				MarkdownDescription: "Append a 127.0.1.1 entry that resolves from FQDN and hostname every boot.",
				Optional:            true,
			},
		},
	}
}
