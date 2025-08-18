package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ByobuModel struct {
	ByobuByDefault types.String `tfsdk:"byobu_by_default"`
}

type ByobuOutputModel struct {
	ByobuByDefault string `yaml:"byobu_by_default,omitempty"`
}

// Byobu
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#byobu
func Byobu() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"byobu_by_default": schema.StringAttribute{
				MarkdownDescription: `
This module controls whether Byobu is enabled or disabled system-wide and for the default system user. If Byobu is to be enabled, this module will ensure it is installed. Likewise, if Byobu is to be disabled, it will be removed (if installed).

Valid configuration options for this module are:
 - enable-system: enable Byobu system-wide
 - enable-user: enable Byobu for the default user
 - disable-system: disable Byobu system-wide
 - disable-user: disable Byobu for the default user
 - enable: enable Byobu both system-wide and for the default user
 - disable: disable Byobu for all users
 - user: alias for enable-user
 - system: alias for enable-system
        `,
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"enable-system",
						"enable-user",
						"disable-system",
						"disable-user",
						"enable",
						"disable",
						"user",
						"system",
					),
				},
			},
		},
	}
}
