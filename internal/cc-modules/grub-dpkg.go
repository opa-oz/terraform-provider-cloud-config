package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type GRUBDpkg struct {
	Enabled                    types.Bool   `tfsdk:"enabled"`
	GRUBPC_InstallDevices      types.String `tfsdk:"grub_pc_install_devices"`
	GRUBPC_InstallDevicesEmpty types.Bool   `tfsdk:"grub_pc_install_devices_empty"`
	GRUBEFI_InstallDevices     types.String `tfsdk:"grub_efi_install_devices"`
}

type GRUBDpkgOutput struct {
	Enabled                    bool   `yaml:"enabled,omitempty"`
	GRUBPC_InstallDevices      string `yaml:"grub-pc/install_devices,omitempty"`
	GRUBPC_InstallDevicesEmpty bool   `yaml:"grub-pc/install_devices_empty,omitempty"`
	GRUBEFI_InstallDevices     string `yaml:"grub-efi/install_devices,omitempty"`
}

type GRUBDpkgModel struct {
	GRUBDpkg *GRUBDpkg `tfsdk:"grub_dpkg"`
}

type GRUBDpkgOutputModel struct {
	GRUBDpkg *GRUBDpkgOutput `yaml:"grub_dpkg,omitempty"`
}

// GRUBDpkgBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#grub-dpkg
func GRUBDpkgBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"grub_dpkg": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("grub_dpkg")),
				},
				MarkdownDescription: `
Configure which device is used as the target for GRUB installation. This module can be enabled/disabled using the enabled config key in the grub_dpkg config dict. This module automatically selects a disk using grub-probe if no installation device is specified.

The value placed into the debconf database is in the format expected by the GRUB post-install script expects. Normally, this is a /dev/disk/by-id/ value, but we do fallback to the plain disk name if a by-id name is not present.

If this module is executed inside a container, then the debconf database is seeded with empty values, and install_devices_empty is set to true.
        `,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether to configure which device is used as the target for grub installation. Default: `false`.",
						Optional:            true,
					},
					"grub_pc_install_devices": schema.StringAttribute{
						MarkdownDescription: "Device to use as target for grub installation. If unspecified, grub-probe of `/boot` will be used to find the device.",
						Optional:            true,
					},
					"grub_pc_install_devices_empty": schema.BoolAttribute{
						MarkdownDescription: "Sets values for `grub-pc/install_devices_empty`. If unspecified, will be set to true if `grub-pc/install_devices` is empty, otherwise false.",
						Optional:            true,
					},
					"grub_efi_install_devices": schema.StringAttribute{
						MarkdownDescription: "Partition to use as target for grub installation. If unspecified, grub-probe of `/boot/efi` will be used to find the partition.",
						Optional:            true,
					},
				},
			},
		},
	}
}
