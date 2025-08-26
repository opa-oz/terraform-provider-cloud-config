package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BootCMDModule struct {
	BootCMD types.List `tfsdk:"bootcmd"`
}

type BootCMDOutputModule struct {
	BootCMD *[]string `yaml:"bootcmd,omitempty"`
}

// BootCMD
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#bootcmd
// TODO: Support `array of strings`
// > If the item is a list, the items will be executed as if passed to execve(3) (with the first argument as the command).
func BootCMD() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"bootcmd": schema.ListAttribute{
				ElementType: types.StringType,
				MarkdownDescription: `
This module runs arbitrary commands very early in the boot process, only slightly after a boothook would run. This is very similar to a boothook, but more user friendly. Commands can be specified as strings.

bootcmd should only be used for things that could not be done later in the boot process.

When writing files, do not use /tmp dir as it races with systemd-tmpfiles-clean (LP: #1707222). Use /run/somedir instead.

Use of INSTANCE_ID variable within this module is deprecated. Use [jinja templates](https://cloudinit.readthedocs.io/en/latest/explanation/format.html#user-data-formats-jinja) with [ v1.instance_id ](https://cloudinit.readthedocs.io/en/latest/explanation/instancedata.html#v1-instance-id) instead.
        `,
				Optional: true,
			},
		},
	}
}
