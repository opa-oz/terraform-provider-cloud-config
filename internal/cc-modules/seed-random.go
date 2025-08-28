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

type RandomSeed struct {
	File            types.String `tfsdk:"file"`
	Data            types.String `tfsdk:"data"`
	Encoding        types.String `tfsdk:"encoding"`
	Command         types.List   `tfsdk:"command"`
	CommandRequired types.Bool   `tfsdk:"command_required"`
}

type RandomSeedOutput struct {
	File            string    `yaml:"file,omitempty"`
	Data            string    `yaml:"data,omitempty"`
	Encoding        string    `yaml:"encoding,omitempty"`
	Command         *[]string `yaml:"command,omitempty"`
	CommandRequired bool      `yaml:"command_required,omitempty"`
}

type SeedRandomModel struct {
	RandomSeed *RandomSeed `tfsdk:"random_seed"`
}

type SeedRandomOutputModel struct {
	RandomSeed *RandomSeedOutput `yaml:"random_seed,omitempty"`
}

// SeedRandomBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#seed-random
func SeedRandomBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"random_seed": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("random_seed")),
				},
				MarkdownDescription: `
All cloud instances started from the same image will produce similar data when they are first booted as they are all starting with the same seed for the kernel’s entropy keyring. To avoid this, random seed data can be provided to the instance, either as a string or by specifying a command to run to generate the data.

Configuration for this module is under the random_seed config key. If the cloud provides its own random seed data, it will be appended to data before it is written to file.

If the command key is specified, the given command will be executed. This will happen after file has been populated. That command’s environment will contain the value of the file key as RANDOM_SEED_FILE. If a command is specified that cannot be run, no error will be reported unless command_required is set to true.
        `,
				Attributes: map[string]schema.Attribute{
					"file": schema.StringAttribute{
						MarkdownDescription: "File to write random data to. Default: `/dev/urandom`.",
						Optional:            true,
					},
					"data": schema.StringAttribute{
						MarkdownDescription: "This data will be written to file before data from the datasource. When using a multi-line value or specifying binary data, be sure to follow YAML syntax and use the | and `!binary` YAML format specifiers when appropriate.",
						Optional:            true,
					},
					"encoding": schema.StringAttribute{
						MarkdownDescription: "Used to decode data provided. Allowed values are `raw`, `base64`, `b64`, `gzip`, or `gz`. Default: `raw`.",
						Optional:            true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"raw",
								"base64",
								"b64",
								"gzip",
								"gz",
							),
						},
					},
					"command": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "Execute this command to seed random. The command will have RANDOM_SEED_FILE in its environment set to the value of file above.",
						Optional:            true,
					},
					"command_required": schema.BoolAttribute{
						MarkdownDescription: "If true, and command is not available to be run then an exception is raised and cloud-init will record failure. Otherwise, only debug error is mentioned. Default: `false`.",
						Optional:            true,
					},
				},
			},
		},
	}
}
