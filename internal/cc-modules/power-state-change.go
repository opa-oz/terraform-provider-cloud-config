package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type PowerState struct {
	Delay        types.Int64  `tfsdk:"delay"`
	NoDelay      types.Bool   `tfsdk:"no_delay"`
	Mode         types.String `tfsdk:"mode"`
	Message      types.String `tfsdk:"message"`
	Timeout      types.Int64  `tfsdk:"timeout"`
	Condition    types.Bool   `tfsdk:"condition"`
	ConditionCmd types.String `tfsdk:"condition_cmd"`
}

type PowerStateOutput struct {
	Delay     any    `yaml:"delay,omitempty"`
	Mode      string `yaml:"mode,omitempty"`
	Message   string `yaml:"message,omitempty"`
	Timeout   int64  `yaml:"timeout,omitempty"`
	Condition any    `yaml:"condition,omitempty"`
}

type PowerStateModel struct {
	PowerState *PowerState `tfsdk:"power_state"`
}

type PowerStateOutputModel struct {
	PowerState *PowerStateOutput `yaml:"power_state,omitempty"`
}

// PowerStateChangeBlock
func PowerStateChangeBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"power_state": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("power_state")),
				},
				MarkdownDescription: `
This module handles shutdown/reboot after all config modules have been run. By default it will take no action, and the system will keep running unless a package installation/upgrade requires a system reboot (e.g. installing a new kernel) and package_reboot_if_required is true.

Using this module ensures that cloud-init is entirely finished with modules that would be executed. An example to distinguish delay from timeout:

If you delay 5 (5 minutes) and have a timeout of 120 (2 minutes), the max time until shutdown will be 7 minutes, though it could be as soon as 5 minutes. Cloud-init will invoke ‘shutdown +5’ after the process finishes, or when ‘timeout’ seconds have elapsed.

**NOTE**
With Alpine Linux any message value specified is ignored as Alpine’s halt, poweroff, and reboot commands do not support broadcasting a message.
        `,
				Attributes: map[string]schema.Attribute{
					"delay": schema.Int64Attribute{
						MarkdownDescription: "Time in minutes to delay after cloud-init has finished. Can be now or an integer specifying the number of minutes to delay. Default: `now`.",
						Optional:            true,
						Validators: []validator.Int64{
							int64validator.ConflictsWith(path.MatchRelative().AtParent().AtName("no_delay")),
						},
					},
					"no_delay": schema.BoolAttribute{
						MarkdownDescription: "Apply changes right after cloud-init finish. Same as `delay: now`. Default `true`",
						Optional:            true,
						Validators: []validator.Bool{
							boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("delay")),
						},
					},
					"mode": schema.StringAttribute{
						MarkdownDescription: "Must be one of poweroff, halt, or reboot.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"poweroff",
								"halt",
								"reboot",
							),
						},
					},
					"message": schema.StringAttribute{
						MarkdownDescription: "Optional message to display to the user when the system is powering off or rebooting.",
						Optional:            true,
					},
					"timeout": schema.Int64Attribute{
						MarkdownDescription: "Time in seconds to wait for the cloud-init process to finish before executing shutdown. Default: `30`.",
						Optional:            true,
					},
					"condition": schema.BoolAttribute{
						MarkdownDescription: "Apply state change only if condition is met. May be boolean true (always met), false (never met).  Defaults to 'true'",
						Optional:            true,
						Validators: []validator.Bool{
							boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("condition_cmd")),
						},
					},
					"condition_cmd": schema.StringAttribute{
						MarkdownDescription: "For command formatting, see the documentation for `cc_runcmd`. If exit code is 0, condition is met, otherwise not.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("condition")),
						},
					},
				},
			},
		},
	}
}
