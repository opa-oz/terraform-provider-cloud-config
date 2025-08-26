package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type NTPConfig struct {
	Confpath    types.String `tfsdk:"confpath"`
	CheckExe    types.String `tfsdk:"check_exe"`
	Packages    types.List   `tfsdk:"packages"`
	ServiceName types.String `tfsdk:"service_name"`
	Template    types.String `tfsdk:"template"`
}

type NTPConfigOutput struct {
	Confpath    string    `yaml:"confpath,omitempty"`
	CheckExe    string    `yaml:"check_exe,omitempty"`
	Packages    *[]string `yaml:"packages,omitempty"`
	ServiceName string    `yaml:"service_name,omitempty"`
	Template    string    `yaml:"template,omitempty"`
}

type NTP struct {
	Pools     types.List   `tfsdk:"pools"`
	Servers   types.List   `tfsdk:"servers"`
	Peers     types.List   `tfsdk:"peers"`
	Allow     types.List   `tfsdk:"allow"`
	NTPClient types.String `tfsdk:"ntp_client"`
	Enabled   types.Bool   `tfsdk:"enabled"`
	Config    *NTPConfig   `tfsdk:"config"`
}

type NTPOutput struct {
	Pools     *[]string        `yaml:"pools,omitempty"`
	Servers   *[]string        `yaml:"servers,omitempty"`
	Peers     *[]string        `yaml:"peers,omitempty"`
	Allow     *[]string        `yaml:"allow,omitempty"`
	NTPClient string           `yaml:"ntp_client,omitempty"`
	Enabled   *bool            `yaml:"enabled,omitempty"`
	Config    *NTPConfigOutput `yaml:"config,omitempty"`
}

type NTPModel struct {
	NTP *NTP `tfsdk:"ntp"`
}

type NTPOutputModel struct {
	NTP *NTPOutput `yaml:"ntp,omitempty"`
}

// NTPBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#ntp
func NTPBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"ntp": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("ntp")),
				},
				MarkdownDescription: `
Handle Network Time Protocol (NTP) configuration. If ntp is not installed on the system and NTP configuration is specified, ntp will be installed.

If there is a default NTP config file in the image or one is present in the distro’s ntp package, it will be copied to a file with .dist appended to the filename before any changes are made.

A list of NTP pools and NTP servers can be provided under the ntp config key.

If no NTP servers or pools are provided, 4 pools will be used in the format:

{0-3}.{distro}.pool.ntp.org
        `,
				Attributes: map[string]schema.Attribute{
					"pools": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "List of ntp pools. If both pools and servers are empty, 4 default pool servers will be provided of the format `{0-3}.{distro}.pool.ntp.org`. NOTE: for Alpine Linux when using the Busybox NTP client this setting will be ignored due to the limited functionality of Busybox’s ntpd.",
						Optional:            true,
					},
					"servers": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "List of ntp servers. If both pools and servers are empty, 4 default pool servers will be provided with the format `{0-3}.{distro}.pool.ntp.org`.",
						Optional:            true,
					},
					"peers": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "List of ntp peers.",
						Optional:            true,
					},
					"allow": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "List of CIDRs to allow.",
						Optional:            true,
					},
					"ntp_client": schema.StringAttribute{
						MarkdownDescription: "Name of an NTP client to use to configure system NTP. When unprovided or ‘auto’ the default client preferred by the distribution will be used. The following built-in client names can be used to override existing configuration defaults: chrony, ntp, openntpd, ntpdate, systemd-timesyncd.",
						Optional:            true,
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Attempt to enable ntp clients if set to True. If set to false, ntp client will not be configured or installed.",
						Optional:            true,
					},
				},
				Blocks: map[string]schema.Block{
					"config": schema.SingleNestedBlock{
						PlanModifiers: []planmodifier.Object{
							utils.NullWhen(path.Root("config")),
						},
						Attributes: map[string]schema.Attribute{
							"confpath": schema.StringAttribute{
								MarkdownDescription: "Configuration settings or overrides for the ntp_client specified.",
								Optional:            true,
							},
							"check_exe": schema.StringAttribute{
								MarkdownDescription: "The executable name for the `ntp_client`. For example, ntp service `check_exe` is ‘ntpd’ because it runs the ntpd binary.",
								Optional:            true,
							},
							"packages": schema.ListAttribute{
								ElementType:         types.StringType,
								MarkdownDescription: "List of packages needed to be installed for the selected **ntp_client**.",
								Optional:            true,
							},
							"service_name": schema.StringAttribute{
								MarkdownDescription: "The systemd or sysvinit service name used to start and stop the ntp_client service.",
								Optional:            true,
							},
							"template": schema.StringAttribute{
								MarkdownDescription: "Inline template allowing users to customize their ntp_client configuration with the use of the Jinja templating engine. The template content should start with `## template:jinja`. Within the template, you can utilize any of the following ntp module config keys: servers, pools, allow, and peers. Each cc_ntp schema config key and expected value type is defined above.",
								Optional:            true,
							},
						},
					},
				},
			},
		},
	}
}
