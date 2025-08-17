package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type User struct {
	Name              types.String `tfsdk:"name"`
	Doas              types.List   `tfsdk:"doas"`
	ExpireDate        types.String `tfsdk:"expiredate"`
	Gecos             types.String `tfsdk:"gecos"`
	HomeDir           types.String `tfsdk:"homedir"`
	Inactive          types.String `tfsdk:"inactive"`
	LockPassword      types.Bool   `tfsdk:"lock_passwd"`
	NoCreateHome      types.Bool   `tfsdk:"no_create_home"`
	NoLogInit         types.Bool   `tfsdk:"no_log_init"`
	NoUserGroup       types.Bool   `tfsdk:"no_user_group"`
	Passwd            types.String `tfsdk:"passwd"`
	HashedPasswd      types.String `tfsdk:"hashed_passwd"`
	PlainTextPasswd   types.String `tfsdk:"plain_text_passwd"`
	CreateGroups      types.Bool   `tfsdk:"create_groups"`
	PrimaryGroup      types.String `tfsdk:"primary_group"`
	SELinuxUser       types.String `tfsdk:"selinux_user"`
	Shell             types.String `tfsdk:"shell"`
	SnapUser          types.String `tfsdk:"snapuser"`
	SSHAuthorizedKeys types.List   `tfsdk:"ssh_authorized_keys"`
	SSHImportId       types.List   `tfsdk:"ssh_import_id"`
	SSHRedirectUser   types.Bool   `tfsdk:"ssh_redirect_user"`
	System            types.Bool   `tfsdk:"system"`
	UID               types.Int32  `tfsdk:"uid"`
	Sudo              types.List   `tfsdk:"sudo"`
}

type UserOutput struct {
	Name              string   `yaml:"name,omitempty"`
	Doas              []string `yaml:"doas,omitempty"`
	ExpireDate        string   `yaml:"expiredate,omitempty"`
	Gecos             string   `yaml:"gecos,omitempty"`
	HomeDir           string   `yaml:"homedir,omitempty"`
	Inactive          string   `yaml:"inactive,omitempty"`
	LockPassword      *bool    `yaml:"lock_passwd,omitempty"`
	NoCreateHome      bool     `yaml:"no_create_home,omitempty"`
	NoLogInit         bool     `yaml:"no_log_init,omitempty"`
	NoUserGroup       bool     `yaml:"no_user_group,omitempty"`
	Passwd            string   `yaml:"passwd,omitempty"`
	HashedPasswd      string   `yaml:"hashed_passwd,omitempty"`
	PlainTextPasswd   string   `yaml:"plain_text_passwd,omitempty"`
	CreateGroups      bool     `yaml:"create_groups,omitempty"`
	PrimaryGroup      string   `yaml:"primary_group,omitempty"`
	SELinuxUser       string   `yaml:"selinux_user,omitempty"`
	Shell             string   `yaml:"shell,omitempty"`
	SnapUser          string   `yaml:"snapuser,omitempty"`
	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
	SSHImportId       []string `yaml:"ssh_import_id,omitempty"`
	SSHRedirectUser   bool     `yaml:"ssh_redirect_user,omitempty"`
	System            bool     `yaml:"system,omitempty"`
	UID               *int32   `yaml:"uid,omitempty"`
	Sudo              []string `yaml:"sudo,omitempty"`
}

type UsersAndGroupsModel struct {
	Groups types.List `tfsdk:"groups"`
	User   *User      `tfsdk:"user"`
	Users  *[]User    `tfsdk:"users"`
}

type UsersAndGroupsOutputModel struct {
	Groups *[]string     `yaml:"groups,omitempty"`
	User   *UserOutput   `yaml:"user,omitempty"`
	Users  *[]UserOutput `yaml:"users,omitempty"`
}

func UsersAndGroups() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"groups": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "[WIP] List of user groups to create",
				Optional:            true,
			},
		},
	}
}

// UsersAndGroupsBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#users-and-groups
func UsersAndGroupsBlock() CCModuleNested {
	userAttributes := map[string]schema.Attribute{
		"doas": schema.ListAttribute{
			ElementType:         types.StringType,
			MarkdownDescription: "List of doas rules to add for a user. doas or opendoas must be installed for rules to take effect.",
			Optional:            true,
		},
		"ssh_authorized_keys": schema.ListAttribute{
			ElementType:         types.StringType,
			MarkdownDescription: "List of SSH keys to add to user’s authkeys file. Can not be combined with **ssh_redirect_user**.",
			Optional:            true,
		},
		"ssh_import_id": schema.ListAttribute{
			ElementType:         types.StringType,
			MarkdownDescription: "List of ssh ids to import for user. Can not be combined with ssh_redirect_user.",
			Optional:            true,
		},
		"sudo": schema.ListAttribute{
			ElementType:         types.StringType,
			MarkdownDescription: "Changed in version 22.2.The value ``false`` is deprecated for this key, use ``null`` instead.",
			Optional:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The user’s login name. Required otherwise user creation will be skipped for this user.",
			Optional:            true,
		},
		"expiredate": schema.StringAttribute{
			MarkdownDescription: "Optional. Date on which the user’s account will be disabled. Default: `null`.",
			Optional:            true,
		},
		"gecos": schema.StringAttribute{
			MarkdownDescription: "Optional comment about the user, usually a comma-separated string of real name and contact information.",
			Optional:            true,
		},
		"homedir": schema.StringAttribute{
			MarkdownDescription: "Optional home dir for user. Default: `/home/<username>`.",
			Optional:            true,
		},
		"inactive": schema.StringAttribute{
			MarkdownDescription: "Optional string representing the number of days until the user is disabled.",
			Optional:            true,
		},
		"lock_passwd": schema.BoolAttribute{ // TODO: True is default value
			MarkdownDescription: "Disable password login. Default: `true`.",
			Optional:            true,
		},
		"no_create_home": schema.BoolAttribute{
			MarkdownDescription: "Do not create home directory. Default: `false`.",
			Optional:            true,
		},
		"no_log_init": schema.BoolAttribute{
			MarkdownDescription: "Do not initialize lastlog and faillog for user. Default: `false`.",
			Optional:            true,
		},
		"no_user_group": schema.BoolAttribute{
			MarkdownDescription: " Do not create group named after user. Default: `false`.",
			Optional:            true,
		},
		"passwd": schema.StringAttribute{
			MarkdownDescription: "Hash of user password applied when user does not exist. This will NOT be applied if the user already exists. To generate this hash, run: `mkpasswd --method=SHA-512 --rounds=500000` **Note**: Your password might possibly be visible to unprivileged users on your system, depending on your cloud’s security model. Check if your cloud’s IMDS server is visible from an unprivileged user to evaluate risk.",
			Optional:            true,
			Sensitive:           true,
		},
		"hashed_passwd": schema.StringAttribute{
			MarkdownDescription: "Hash of user password to be applied. This will be applied even if the user is preexisting. To generate this hash, run: `mkpasswd --method=SHA-512 --rounds=500000`. **Note**: Your password might possibly be visible to unprivileged users on your system, depending on your cloud’s security model. Check if your cloud’s IMDS server is visible from an unprivileged user to evaluate risk.",
			Optional:            true,
			Sensitive:           true,
		},
		"plain_text_passwd": schema.StringAttribute{
			MarkdownDescription: "Clear text of user password to be applied. This will be applied even if the user is preexisting. **Note**: SSH keys or certificates are a safer choice for logging in to your system. For local escalation, supplying a hashed password is a safer choice than plain text. Your password might possibly be visible to unprivileged users on your system, depending on your cloud’s security model. An exposed plain text password is an immediate security concern. Check if your cloud’s IMDS server is visible from an unprivileged user to evaluate risk.",
			Optional:            true,
			Sensitive:           true,
		},
		"create_groups": schema.BoolAttribute{ // TODO: True is default value
			MarkdownDescription: "Boolean set `false` to disable creation of specified user groups. Default: `true`.",
			Optional:            true,
		},
		"primary_group": schema.StringAttribute{
			MarkdownDescription: "Primary group for user. Default: `<username>`.",
			Optional:            true,
		},
		"selinux_user": schema.StringAttribute{
			MarkdownDescription: "SELinux user for user’s login. Default: the default SELinux user.",
			Optional:            true,
		},
		"shell": schema.StringAttribute{
			MarkdownDescription: "Path to the user’s login shell. Default: the host system’s default shell.",
			Optional:            true,
		},
		"snapuser": schema.StringAttribute{
			MarkdownDescription: "Specify an email address to create the user as a Snappy user through snap `create-user`. If an Ubuntu SSO account is associated with the address, username and SSH keys will be requested from there.",
			Optional:            true,
		},
		"ssh_redirect_user": schema.BoolAttribute{
			MarkdownDescription: "Boolean set to true to disable SSH logins for this user. When specified, all cloud-provided public SSH keys will be set up in a disabled state for this username. Any SSH login as this username will timeout and prompt with a message to login instead as the **default_username** for this instance. Default: `false`. This key can not be combined with **ssh_import_id** or **ssh_authorized_keys**.",
			Optional:            true,
			Validators: []validator.Bool{
				boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_authorized_keys")),
				boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("ssh_import_id")),
			},
		},
		"system": schema.BoolAttribute{
			MarkdownDescription: "Create user as system user with no home directory. Default: `false`.",
			Optional:            true,
		},
		"uid": schema.Int32Attribute{
			MarkdownDescription: "The user’s ID. Default value [system default].",
			Optional:            true,
		},
	}

	return CCModuleNested{
		block: map[string]schema.Block{
			"user": schema.SingleNestedBlock{
				MarkdownDescription: "The user dictionary values override the `default_user` configuration from `/etc/cloud/cloud.cfg`. The user dictionary keys supported for the `default_user` are the same as the users schema.",
				Attributes:          userAttributes,
			},
			"users": schema.ListNestedBlock{
				MarkdownDescription: "List of users",
				NestedObject: schema.NestedBlockObject{
					Attributes: userAttributes,
				},
			},
		},
	}
}
