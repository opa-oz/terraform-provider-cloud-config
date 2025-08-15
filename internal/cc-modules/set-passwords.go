package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ChPasswdUser struct {
	Name     types.String `tfsdk:"name"`
	Password types.String `tfsdk:"password"`
	Type     types.String `tfsdk:"type"`
}

type ChPasswdUserOutput struct {
	Name     string `yaml:"name,omitempty"`
	Password string `yaml:"password,omitempty"`
	Type     string `yaml:"type,omitempty"`
}

type ChPasswd struct {
	Users  *[]ChPasswdUser `tfsdk:"users"`
	Expire types.Bool      `tfsdk:"expire"`
}

type ChPasswdOutput struct {
	Users  *[]ChPasswdUserOutput `yaml:"users,omitempty"`
	Expire *bool                 `yaml:"expire,omitempty"`
}

type SetPasswordsModel struct {
	SSHPwauth types.Bool `tfsdk:"ssh_pwauth"`

	ChPasswd *ChPasswd `tfsdk:"chpasswd"`
}

type SetPasswordsOutputModel struct {
	SSHPwauth *bool           `yaml:"ssh_pwauth,omitempty"` // NOTE: `false` needs to be set explicitly
	ChPasswd  *ChPasswdOutput `yaml:"chpasswd,omitempty"`
}

// SetPasswords
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#set-passwords
func SetPasswords() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"ssh_pwauth": schema.BoolAttribute{
				MarkdownDescription: `
Sets whether or not to accept password authentication. true will enable password auth. false will disable. *Default*: leave the value unchanged. 

In order for this config to be applied, SSH may need to be restarted. On systemd systems, this restart will only happen if the SSH service has already been started. On non-systemd systems, a restart will be attempted regardless of the service state.

_Changed in version 22.3. Use of non-boolean values for this field is deprecated._
        `,
				Optional: true,
			},
		},
	}
}

func SetPasswordsBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"chpasswd": schema.SingleNestedBlock{
				MarkdownDescription: `
The chpasswd config key accepts a dictionary containing either (or both) of users and expire.

 - The users key is used to assign a password to a corresponding pre-existing user.
 - The expire key is used to set whether to expire all user passwords specified by this module, such that a password will need to be reset on the user’s next login.
        `,
				Attributes: map[string]schema.Attribute{
					"expire": schema.BoolAttribute{
						MarkdownDescription: "Whether to expire all user passwords such that a password will need to be reset on the user’s next login. Default: `true`.",
						Optional:            true,
					},
				},
				Blocks: map[string]schema.Block{
					"users": schema.ListNestedBlock{
						MarkdownDescription: "This key represents a list of existing users to set passwords for. Each item under users contains the following required keys: *name* and *password* or in the case of a randomly generated password, *name* and *type*. The *type* key has a default value of 'hash', and may alternatively be set to 'text' or 'RANDOM'. Randomly generated passwords may be insecure, use at your own risk.",
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "User's name",
									Required:            true,
								},
								"password": schema.StringAttribute{ // TODO: Do proper validation
									MarkdownDescription: "User's password",
									Optional:            true,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "The *type* key has a default value of 'hash', and may alternatively be set to 'text' or 'RANDOM'.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}
