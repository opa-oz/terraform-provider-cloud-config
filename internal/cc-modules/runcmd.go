package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RunCMDModule struct {
	RunCMD types.List `tfsdk:"runcmd"`
}

type RunCMDOutputModule struct {
	RunCMD []string `yaml:"runcmd"`
}

// RunCMD
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#runcmd
// TODO: Support `array of strings`
// > If the item is a list, the items will be executed as if passed to execve(3) (with the first argument as the command).
func RunCMD() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"runcmd": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Run arbitrary commands at a rc.local-like time-frame with output to the console. Each item will be interpreted by `sh`.",
				Optional:            true,
			},
		},
	}
}
