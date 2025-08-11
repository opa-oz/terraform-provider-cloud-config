package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SetHostnameModel struct {
	Hostname types.String `tfsdk:"hostname"`
	FQDN     types.String `tfsdk:"fqdn"`

	PreferFQDNOverHostname types.Bool `tfsdk:"prefer_fqdn_over_hostname"`
	PreserveHostname       types.Bool `tfsdk:"preserve_hostname"`
	CreateHostnameFile     types.Bool `tfsdk:"create_hostname_file"`
}

type SetHostnameOutputModel struct {
	Hostname string `yaml:"hostname,omitempty"`
	FQDN     string `yaml:"fqdn,omitempty"`

	PreferFQDNOverHostname bool  `yaml:"prefer_fqdn_over_hostname,omitempty"`
	PreserveHostname       bool  `yaml:"preserve_hostname,omitempty"`
	CreateHostnameFile     *bool `yaml:"create_hostname_file,omitempty"` // WARN: Pointer, because default value is `true`
}

// SetHostname
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#set-hostname
func SetHostname() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				MarkdownDescription: "The hostname to set.",
				Optional:            true,
			},
			"fqdn": schema.StringAttribute{
				MarkdownDescription: "The fully qualified domain name to set.",
				Optional:            true,
			},
			"prefer_fqdn_over_hostname": schema.BoolAttribute{
				MarkdownDescription: "If true, the fqdn will be used if it is set. If false, the hostname will be used. If unset, the result is distro-dependent.",
				Optional:            true,
			},
			"preserve_hostname": schema.BoolAttribute{
				MarkdownDescription: "If true, the hostname will not be changed. Default: `false`.",
				// Default:             booldefault.StaticBool(false),
				Optional: true,
			},
			"create_hostname_file": schema.BoolAttribute{
				MarkdownDescription: "If `false`, the hostname file (e.g. `/etc/hostname`) will not be created if it does not exist. On systems that use systemd, setting `create_hostname_file` to `false` will set the hostname transiently. If true, the hostname file will always be created and the hostname will be set statically on systemd systems. Default: `true`.",
				// Default:             booldefault.StaticBool(true),
				Optional: true,
			},
		},
	}
}
