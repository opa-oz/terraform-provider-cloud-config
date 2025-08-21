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

func transformBootCMD(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if !model.BootCMD.IsUnknown() {
		elems := model.BootCMD.Elements()

		if len(elems) > 0 {
			cmds := make([]string, len(elems))
			diagnostics := model.BootCMD.ElementsAs(ctx, &cmds, false)

			if diagnostics.HasError() {
				return diagnostics
			}
			output.BootCMD = cmds
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

func transformApkConfigure(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	transformApkRepo := func(repo *ccmodules.ApkRepo) (ccmodules.ApkRepoOutput, diag.Diagnostics) {
		apkRepo := ccmodules.ApkRepoOutput{}

		apkRepo.PreserveRepositories = repo.PreserveRepositories.ValueBool()
		apkRepo.LocalRepoBaseUrl = repo.LocalRepoBaseUrl.ValueString()

		if repo.AlpineRepo != nil {
			alpineRepo := ccmodules.AlpineRepoOutput{}
			alpineRepo.CommunityEnabled = repo.AlpineRepo.CommunityEnabled.ValueBool()
			alpineRepo.TestingEnabled = repo.AlpineRepo.TestingEnabled.ValueBool()

			alpineRepo.BaseUrl = repo.AlpineRepo.BaseUrl.ValueString()
			alpineRepo.Version = repo.AlpineRepo.Version.ValueString()

			apkRepo.AlpineRepo = &alpineRepo
		}

		return apkRepo, nil
	}

	if model.ApkRepos != nil {
		repo, diagnostics := transformApkRepo(model.ApkRepos)
		if diagnostics.HasError() {
			return diagnostics
		}
		output.ApkRepos = &repo
	}
	return nil
}

func transformAptPipelining(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.AptPipelining == nil {
		return nil
	}

	if model.AptPipelining.OS.ValueBool() {
		// if `os` is true, use value would be string
		output.AptPipelining = "os"
	} else if model.AptPipelining.Disable.ValueBool() {
		// if `disable` is true, it should be false
		output.AptPipelining = false
	} else if !model.AptPipelining.Depth.IsNull() {
		// if `depth` is configured, use as a number
		output.AptPipelining = model.AptPipelining.Depth.ValueInt32Pointer()
	}

	return nil
}

func transformCACertificatesHosts(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.CACerts == nil {
		return nil
	}

	caCerts := ccmodules.CACertsOutput{}

	caCerts.RemoveDefaults = model.CACerts.RemoveDefaults.ValueBool()

	if !model.CACerts.Trusted.IsUnknown() {
		elems := model.CACerts.Trusted.Elements()

		if len(elems) > 0 {
			certs := make([]string, len(elems))
			diagnostics := model.CACerts.Trusted.ElementsAs(ctx, &certs, false)

			if diagnostics.HasError() {
				return diagnostics
			}

			caCerts.Trusted = &certs
		}
	}

	output.CACerts = &caCerts

	return nil
}

func transformFan(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.Fan == nil {
		return nil
	}

	fan := ccmodules.FanOutput{}

	fan.Config = model.Fan.Config.ValueString()
	fan.ConfigPath = model.Fan.ConfigPath.ValueString()

	output.Fan = &fan

	return nil
}

func transformGrowpart(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.Growpart == nil {
		return nil
	}

	growpart := ccmodules.GrowpartOutput{}

	growpart.IgnoreGrowrootDisabled = model.Growpart.IgnoreGrowrootDisabled.ValueBool()
	growpart.Mode = model.Growpart.Mode.ValueString()

	if !model.Growpart.Devices.IsUnknown() {
		elems := model.Growpart.Devices.Elements()

		if len(elems) > 0 {
			certs := make([]string, len(elems))
			diagnostics := model.Growpart.Devices.ElementsAs(ctx, &certs, false)

			if diagnostics.HasError() {
				return diagnostics
			}

			growpart.Devices = &certs
		}
	}

	output.Growpart = &growpart

	return nil
}

func transformGRUBDpkg(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.GRUBDpkg == nil {
		return nil
	}

	config := ccmodules.GRUBDpkgOutput{}

	config.Enabled = model.GRUBDpkg.Enabled.ValueBool()
	config.GRUBPC_InstallDevicesEmpty = model.GRUBDpkg.GRUBPC_InstallDevicesEmpty.ValueBool()
	config.GRUBPC_InstallDevices = model.GRUBDpkg.GRUBPC_InstallDevices.ValueString()
	config.GRUBEFI_InstallDevices = model.GRUBDpkg.GRUBEFI_InstallDevices.ValueString()

	output.GRUBDpkg = &config

	return nil
}

func transformInstallHotplug(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.Updates == nil || model.Updates.Network == nil {
		return nil
	}

	config := ccmodules.UpdatesOutput{
		Network: &ccmodules.NetworkOutput{},
	}

	if !model.Updates.Network.When.IsUnknown() {
		elems := model.Updates.Network.When.Elements()

		if len(elems) > 0 {
			whens := make([]string, len(elems))
			diagnostics := model.Updates.Network.When.ElementsAs(ctx, &whens, false)

			if diagnostics.HasError() {
				return diagnostics
			}

			config.Network.When = &whens
			output.Updates = &config
		}
	}

	return nil
}

func transformKeyboard(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.Keyboard == nil {
		return nil
	}

	config := ccmodules.KeyboardOutput{
		Layout:  model.Keyboard.Layout.ValueString(),
		Model:   model.Keyboard.Model.ValueString(),
		Variant: model.Keyboard.Variant.ValueString(),
		Options: model.Keyboard.Options.ValueString(),
	}

	output.Keyboard = &config

	return nil
}

func transformKeysToConsole(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.SSH != nil {
		// NOTE: default value is anyway `true`
		if !model.SSH.EmitKeysToConsole.IsNull() && !model.SSH.EmitKeysToConsole.ValueBool() {
			output.SSH = &ccmodules.SSHOutput{
				EmitKeysToConsole: model.SSH.EmitKeysToConsole.ValueBoolPointer(),
			}
		}
	}

	if !model.SSHKeyConsoleBlacklist.IsUnknown() {
		elems := model.SSHKeyConsoleBlacklist.Elements()

		if len(elems) > 0 {
			whens := make([]string, len(elems))
			diagnostics := model.SSHKeyConsoleBlacklist.ElementsAs(ctx, &whens, false)

			if diagnostics.HasError() {
				return diagnostics
			}

			output.SSHKeyConsoleBlacklist = &whens
		}
	}

	if !model.SSHFPConsoleBlacklist.IsUnknown() {
		elems := model.SSHFPConsoleBlacklist.Elements()

		if len(elems) > 0 {
			whens := make([]string, len(elems))
			diagnostics := model.SSHFPConsoleBlacklist.ElementsAs(ctx, &whens, false)

			if diagnostics.HasError() {
				return diagnostics
			}

			output.SSHFPConsoleBlacklist = &whens
		}
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

	diagnostics = transformBootCMD(ctx, &output, model)
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

	diagnostics = transformApkConfigure(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformAptPipelining(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	output.ByobuByDefault = model.ByobuByDefault.ValueString()

	diagnostics = transformCACertificatesHosts(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformFan(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	output.FinalMessage = model.FinalMessage.ValueString()

	diagnostics = transformGrowpart(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformGRUBDpkg(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformInstallHotplug(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformKeyboard(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformKeysToConsole(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

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
