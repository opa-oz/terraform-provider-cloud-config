package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type Growpart struct {
	Mode                   types.String `tfsdk:"mode"`
	Devices                types.List   `tfsdk:"devices"`
	IgnoreGrowrootDisabled types.Bool   `tfsdk:"ignore_growroot_disabled"`
}

type GrowpartOutput struct {
	Mode                   string    `yaml:"mode"`
	Devices                *[]string `yaml:"devices"`
	IgnoreGrowrootDisabled bool      `yaml:"ignore_growroot_disabled"`
}

type GrowpartModel struct {
	Growpart *Growpart `tfsdk:"growpart"`
}

type GrowpartOutputModel struct {
	Growpart *GrowpartOutput `yaml:"growpart,omitempty"`
}

// GrowpartBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#growpart
func GrowpartBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"growpart": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("growpart")),
				},
				MarkdownDescription: `
Growpart resizes partitions to fill the available disk space. This is useful for cloud instances with a larger amount of disk space available than the pristine image uses, as it allows the instance to automatically make use of the extra space.

Note that this only works if the partition to be resized is the last one on a disk with classic partitioning scheme (MBR, BSD, GPT). LVM, Btrfs and ZFS have no such restrictions.

The devices on which to run growpart are specified as a list under the devices key.

There is some functionality overlap between this module and the growroot functionality of cloud-initramfs-tools. However, there are some situations where one tool is able to function and the other is not. The default configuration for both should work for most cloud instances. To explicitly prevent cloud-initramfs-tools from running growroot, the file /etc/growroot-disabled can be created.

By default, both growroot and cc_growpart will check for the existence of this file and will not run if it is present. However, this file can be ignored for cc_growpart by setting ignore_growroot_disabled to true. Read more about cloud-initramfs-tools.

On FreeBSD, there is also the growfs service, which has a lot of overlap with cc_growpart and cc_resizefs, but only works on the root partition. In that configuration, we use it, otherwise, we fall back to gpart.

**Note**: growfs may insert a swap partition, if none is present, unless instructed not to via growfs_swap_size=0 in either kenv(1), or rc.conf(5).

Growpart is enabled by default on the root partition.
        `,
				Attributes: map[string]schema.Attribute{
					"mode": schema.StringAttribute{
						MarkdownDescription: `
The utility to use for resizing. Default: auto

Possible options:
 - auto - Use any available utility
 - growpart - Use growpart utility
 - gpart - Use BSD gpart utility
 - 'off' - Take no action.
            `,
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"auto",
								"growpart",
								"gpart",
								"off",
							),
						},
					},
					"devices": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "The devices to resize. Each entry can either be the path to the device’s mountpoint in the filesystem or a path to the block device in ‘/dev’. Default: `[/]`",
						Optional:            true,
					},
					"ignore_growroot_disabled": schema.BoolAttribute{
						MarkdownDescription: "If true, ignore the presence of `/etc/growroot-disabled`. If false and the file exists, then don’t resize. Default: `false`.",
						Optional:            true,
					},
				},
			},
		},
	}
}
