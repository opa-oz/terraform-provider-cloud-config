package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type SSHObj struct {
	EmitKeysToConsole types.Bool `tfsdk:"emit_keys_to_console"`
}

type SSHOutput struct {
	EmitKeysToConsole *bool `yaml:"emit_keys_to_console,omitempty"`
}

type KeysToConsoleModel struct {
	SSH                    *SSHObj    `tfsdk:"ssh"`
	SSHKeyConsoleBlacklist types.List `tfsdk:"ssh_key_console_blacklist"`
	SSHFPConsoleBlacklist  types.List `tfsdk:"ssh_fp_console_blacklist"`
}

type KeysToConsoleOutputModel struct {
	SSH                    *SSHOutput `yaml:"ssh,omitempty"`
	SSHKeyConsoleBlacklist *[]string  `yaml:"ssh_key_console_blacklist,omitempty"`
	SSHFPConsoleBlacklist  *[]string  `yaml:"ssh_fp_console_blacklist,omitempty"`
}

// KeysToConsole
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#keys-to-console
func KeysToConsole() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"ssh_key_console_blacklist": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Avoid printing matching SSH key types to the system console.",
				Optional:            true,
			},
			"ssh_fp_console_blacklist": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Avoid printing matching SSH fingerprints to the system console.",
				Optional:            true,
			},
		},
	}
}

// KeysToConsoleBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#keys-to-console
func KeysToConsoleBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"ssh": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("ssh")),
				},
				MarkdownDescription: `
For security reasons it may be desirable not to write SSH host keys and their fingerprints to the console. To avoid either of them being written to the console, the emit_keys_to_console config key under the main ssh config key can be used.

To avoid the fingerprint of types of SSH host keys being written to console the ssh_fp_console_blacklist config key can be used. By default, all types of keys will have their fingerprints written to console.

To avoid host keys of a key type being written to console the ssh_key_console_blacklist config key can be used. By default, all supported host keys are written to console.
        `,
				Attributes: map[string]schema.Attribute{
					"emit_keys_to_console": schema.BoolAttribute{
						MarkdownDescription: "Set false to avoid printing SSH keys to system console. Default: `true`.",
						Optional:            true,
					},
				},
			},
		},
	}
}
