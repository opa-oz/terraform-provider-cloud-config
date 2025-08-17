package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	ccmodules "github.com/opa-oz/terraform-provider-cloud-config/internal/cc-modules"

	"gopkg.in/yaml.v3"
)

const (
	hat = "#cloud-config"
)

func transformSetHostname(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	output.Hostname = model.Hostname.ValueString()
	output.FQDN = model.FQDN.ValueString()
	output.PreserveHostname = model.PreserveHostname.ValueBool()
	output.PreferFQDNOverHostname = model.PreferFQDNOverHostname.ValueBool()

	// NOTE: default value is anyway `true`
	if !model.CreateHostnameFile.IsNull() && !model.CreateHostnameFile.ValueBool() {
		output.CreateHostnameFile = model.CreateHostnameFile.ValueBoolPointer()
	}

	return nil
}

func transformLocale(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if !model.Locale.IsUnknown() {
		output.Locale = model.Locale.ValueString()
	}
	if !model.LocaleConfigfile.IsUnknown() {
		output.LocaleConfigfile = model.LocaleConfigfile.ValueString()
	}
	return nil
}

func transformRunCMD(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if !model.RunCMD.IsUnknown() {
		elems := model.RunCMD.Elements()

		if len(elems) > 0 {
			cmds := make([]string, len(elems))
			diagnostics := model.RunCMD.ElementsAs(ctx, &cmds, false)

			if diagnostics.HasError() {
				return diagnostics
			}
			output.RunCMD = cmds
		}
	}
	return nil
}

func transformManageEtcHosts(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if !model.ManageEtcHostsLocalhost.IsNull() && model.ManageEtcHostsLocalhost.ValueBool() {
		output.ManageEtcHosts = "localhost"
	} else if !model.ManageEtcHosts.IsNull() {
		output.ManageEtcHosts = model.ManageEtcHosts.ValueBool()
	}

	return nil
}

func transformSSH(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if !model.SSHAuthorizedKeys.IsUnknown() {
		elems := model.SSHAuthorizedKeys.Elements()

		if len(elems) > 0 {
			sshAuthKeys := make([]string, len(elems))
			diagnostics := model.SSHAuthorizedKeys.ElementsAs(ctx, &sshAuthKeys, false)

			if diagnostics.HasError() {
				return diagnostics
			}
			output.SSHAuthorizedKeys = sshAuthKeys
		}
	}
	return nil
}

func transformSetPasswords(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if !model.SSHPwauth.IsNull() {
		output.SSHPwauth = model.SSHPwauth.ValueBoolPointer()
	}

	if model.ChPasswd != nil {
		output.ChPasswd = &ccmodules.ChangePasswordOutput{}

		md := model.ChPasswd

		if !md.Expire.IsNull() {
			output.ChPasswd.Expire = md.Expire.ValueBoolPointer()
		}

		if md.Users != nil {
			usrs := make([]ccmodules.ChangePasswordUserOutput, len(*md.Users))
			for i, usr := range *md.Users {
				newUsr := ccmodules.ChangePasswordUserOutput{}

				if !usr.Name.IsNull() {
					newUsr.Name = usr.Name.ValueString()
				}
				if !usr.Password.IsNull() {
					newUsr.Password = usr.Password.ValueString()
				}
				if !usr.Type.IsNull() {
					newUsr.Type = usr.Type.ValueString()
				}

				usrs[i] = newUsr
			}

			output.ChPasswd.Users = &usrs
		}
	}

	return nil
}

func transformPkgUpdateUpgrade(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	output.PackageUpdate = model.PackageUpdate.ValueBool()
	output.PackageUpgrade = model.PackageUpgrade.ValueBool()
	output.PackageRebootIfRequired = model.PackageRebootIfRequired.ValueBool()

	if !model.Packages.IsUnknown() {
		elems := model.Packages.Elements()

		if len(elems) > 0 {
			cmds := make([]string, len(elems))
			diagnostics := model.Packages.ElementsAs(ctx, &cmds, false)

			if diagnostics.HasError() {
				return diagnostics
			}
			output.Packages = cmds
		}
	}
	return nil
}

func transformUsersAndGroups(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	transformUser := func(user *ccmodules.User) (ccmodules.UserOutput, diag.Diagnostics) {
		out := ccmodules.UserOutput{}

		out.Name = user.Name.ValueString()
		out.ExpireDate = user.ExpireDate.ValueString()
		out.Gecos = user.Gecos.ValueString()
		out.HomeDir = user.HomeDir.ValueString()
		out.Inactive = user.Inactive.ValueString()
		out.Passwd = user.Passwd.ValueString()
		out.HashedPasswd = user.HashedPasswd.ValueString()
		out.PlainTextPasswd = user.PlainTextPasswd.ValueString()
		out.PrimaryGroup = user.PrimaryGroup.ValueString()
		out.SELinuxUser = user.SELinuxUser.ValueString()
		out.Shell = user.Shell.ValueString()
		out.SnapUser = user.SnapUser.ValueString()

		// NOTE: default value is anyway `true`
		if !user.LockPassword.IsNull() && !user.LockPassword.ValueBool() {
			out.LockPassword = user.LockPassword.ValueBoolPointer()
		}

		out.NoCreateHome = user.NoCreateHome.ValueBool()
		out.NoLogInit = user.NoLogInit.ValueBool()
		out.NoUserGroup = user.NoUserGroup.ValueBool()
		out.CreateGroups = user.CreateGroups.ValueBool()
		out.SSHRedirectUser = user.SSHRedirectUser.ValueBool()
		out.System = user.System.ValueBool()

		if !user.UID.IsNull() {
			out.UID = user.UID.ValueInt32Pointer()
		}

		if !user.Doas.IsUnknown() {
			elems := user.Doas.Elements()

			if len(elems) > 0 {
				cmds := make([]string, len(elems))
				diagnostics := user.Doas.ElementsAs(ctx, &cmds, false)

				if diagnostics.HasError() {
					return out, diagnostics
				}
				out.Doas = cmds
			}
		}

		if !user.SSHAuthorizedKeys.IsUnknown() {
			elems := user.SSHAuthorizedKeys.Elements()

			if len(elems) > 0 {
				cmds := make([]string, len(elems))
				diagnostics := user.SSHAuthorizedKeys.ElementsAs(ctx, &cmds, false)

				if diagnostics.HasError() {
					return out, diagnostics
				}
				out.SSHAuthorizedKeys = cmds
			}
		}

		if !user.SSHImportId.IsUnknown() {
			elems := user.SSHImportId.Elements()

			if len(elems) > 0 {
				cmds := make([]string, len(elems))
				diagnostics := user.SSHImportId.ElementsAs(ctx, &cmds, false)

				if diagnostics.HasError() {
					return out, diagnostics
				}
				out.SSHImportId = cmds
			}
		}

		if !user.Sudo.IsUnknown() {
			elems := user.Sudo.Elements()

			if len(elems) > 0 {
				cmds := make([]string, len(elems))
				diagnostics := user.Sudo.ElementsAs(ctx, &cmds, false)

				if diagnostics.HasError() {
					return out, diagnostics
				}
				out.Sudo = cmds
			}
		}

		if !user.Groups.IsUnknown() {
			elems := user.Groups.Elements()

			if len(elems) > 0 {
				cmds := make([]string, len(elems))
				diagnostics := user.Groups.ElementsAs(ctx, &cmds, false)

				if diagnostics.HasError() {
					return out, diagnostics
				}
				out.Groups = cmds
			}
		}

		return out, nil
	}

	if !model.Groups.IsUnknown() {
		elems := model.Groups.Elements()

		if len(elems) > 0 {
			cmds := make([]string, len(elems))
			diagnostics := model.Groups.ElementsAs(ctx, &cmds, false)

			if diagnostics.HasError() {
				return diagnostics
			}
			output.Groups = &cmds
		}
	}

	if model.User != nil {
		user, diagnostics := transformUser(model.User)
		if diagnostics.HasError() {
			return diagnostics
		}
		output.User = &user
	}

	if model.Users != nil {
		usrs := make([]ccmodules.UserOutput, len(*model.Users))
		for i, usr := range *model.Users {
			user, diagnostics := transformUser(&usr)
			if diagnostics.HasError() {
				return diagnostics
			}
			usrs[i] = user
		}

		output.Users = &usrs
	}

	return nil
}

func transform(ctx context.Context, model CloudConfigResourceModel) (ExportModel, diag.Diagnostics) {
	output := ExportModel{}

	diagnostics := transformSetHostname(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformLocale(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	output.Timezone = model.Timezone.ValueString()

	diagnostics = transformRunCMD(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformManageEtcHosts(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformSSH(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformSetPasswords(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformPkgUpdateUpgrade(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformUsersAndGroups(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	output.DisableEC2Metadata = model.DisableEC2Metadata.ValueBool()

	return output, nil
}

func ExportContent(ctx context.Context, model CloudConfigResourceModel) (string, diag.Diagnostics) {
	output, diagnostics := transform(ctx, model)
	if diagnostics != nil {
		return "", diagnostics
	}
	yaml, err := yaml.Marshal(output)

	if err != nil {
		return "", diag.Diagnostics{
			diag.NewErrorDiagnostic("Cannot marshal YAML", err.Error()),
		}
	}

	return strings.TrimSpace(fmt.Sprintf(`%s
%s
  `, hat, yaml)), nil
}
