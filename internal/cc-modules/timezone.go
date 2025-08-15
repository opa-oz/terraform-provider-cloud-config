package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TimezoneModel struct {
	Timezone types.String `tfsdk:"timezone"`
}

type TimezoneOutputModel struct {
	Timezone string `yaml:"timezone,omitempty"`
}

// Timezone
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#timezone
func Timezone() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"timezone": schema.StringAttribute{
				MarkdownDescription: "The timezone to use as represented in /usr/share/zoneinfo.",
				Optional:            true,
			},
		},
	}
}
