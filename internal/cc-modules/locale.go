package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type LocaleModel struct {
	Locale           types.String `tfsdk:"locale"`
	LocaleConfigfile types.String `tfsdk:"locale_configfile"`
}

type LocaleOutputModel struct {
	Locale           string `yaml:"locale,omitempty"`
	LocaleConfigfile string `yaml:"locale_configfile,omitempty"`
}

// Locale
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#locale
func Locale() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"locale": schema.StringAttribute{ // TODO: Actually it supports `boolean` as well, but idk how to do it smoothly rn
				MarkdownDescription: "The locale to set as the system’s locale (e.g. ar_PS).",
				Optional:            true,
			},
			"locale_configfile": schema.StringAttribute{
				MarkdownDescription: "The file in which to write the locale configuration (defaults to the distro’s default location).",
				Optional:            true,
			},
		},
	}
}
