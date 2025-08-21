package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ResizefsModel struct {
	Resizefs        types.Bool `tfsdk:"resize_rootfs"`
	ResizefsNoBlock types.Bool `tfsdk:"resize_rootfs_no_block"`
}

type ResizefsOutputModel struct {
	Resizefs any `yaml:"resize_rootfs,omitempty"`
}

// Resizefs
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#final-message
func Resizefs() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"resize_rootfs": schema.BoolAttribute{
				MarkdownDescription: `
Resize a filesystem to use all available space on partition. This module is useful along with cc_growpart and will ensure that if the root partition has been resized, the root filesystem will be resized along with it.

By default, cc_resizefs will resize the root partition and will block the boot process while the resize command is running.

This module can be disabled altogether by setting resize_rootfs to false.

Default: true
        `,
				Optional: true,
			},
			"resize_rootfs_no_block": schema.BoolAttribute{
				MarkdownDescription: `
Optionally, the resize operation can be performed in the background while cloud-init continues running modules. This can be enabled by setting resize_rootfs to noblock.
        `,
				Optional: true,
			},
		},
	}
}
