package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PkgUpdateUpgradeModel struct {
	PackageUpdate           types.Bool `tfsdk:"package_update"`
	PackageUpgrade          types.Bool `tfsdk:"package_upgrade"`
	PackageRebootIfRequired types.Bool `tfsdk:"package_reboot_if_required"`

	Packages types.List `tfsdk:"packages"`
}

type PkgUpdateUpgradeOutputModel struct {
	PackageUpdate           bool `yaml:"package_update,omitempty"`
	PackageUpgrade          bool `yaml:"package_upgrade,omitempty"`
	PackageRebootIfRequired bool `yaml:"package_reboot_if_required,omitempty"`

	Packages []string `yaml:"packages,omitempty"`
}

// PkgUpdateUpgrade
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#package-update-upgrade-install
func PkgUpdateUpgrade() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"package_update": schema.BoolAttribute{
				MarkdownDescription: "Set `true` to update packages. Happens before upgrade or install. Default: `false`.",
				Optional:            true,
			},
			"package_upgrade": schema.BoolAttribute{
				MarkdownDescription: "Set `true` to upgrade packages. Happens before install. Default: `false`.",
				Optional:            true,
			},
			"package_reboot_if_required": schema.BoolAttribute{
				MarkdownDescription: "Set `true` to reboot the system if required by presence of `/var/run/reboot-required`. Default: `false`.",
				Optional:            true,
			},
			"packages": schema.ListAttribute{ // TODO: Proper type (array of object/array of string/string)
				ElementType:         types.StringType,
				MarkdownDescription: "An array containing a package specification",
				Optional:            true,
			},
		},
	}
}
