package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ccmodules "github.com/opa-oz/terraform-provider-cloud-config/internal/cc-modules"

	"gopkg.in/yaml.v3"
)

const (
	hat = "#cloud-config"
)

func castArray[T any](ctx context.Context, arr types.List) (*[]T, diag.Diagnostics) {
	elems := arr.Elements()

	if len(elems) > 0 {
		cmds := make([]T, len(elems))
		diagnostics := arr.ElementsAs(ctx, &cmds, false)

		if diagnostics.HasError() {
			return nil, diagnostics
		}

		return &cmds, nil
	}

	return nil, nil
}

func transformWriteFiles(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.WriteFiles.IsUnknown() {
		return nil
	}
	length := len(model.WriteFiles.Elements())
	if length == 0 {
		return nil
	}

	writeFiles := make([]ccmodules.WriteFileOutput, length)
	res, diagnostics := castArray[ccmodules.WriteFile](ctx, model.WriteFiles)

	if diagnostics.HasError() {
		return diagnostics
	}

	for k, v := range *res {
		item := ccmodules.WriteFileOutput{
			Path:        v.Path.ValueString(),
			Content:     v.Content.ValueString(),
			Owner:       v.Owner.ValueString(),
			Permissions: v.Permissions.ValueString(),
			Encoding:    v.Encoding.ValueString(),
			Append:      v.Append.ValueBool(),
			Defer:       v.Defer.ValueBool(),
		}

		if v.Source != nil {
			src := ccmodules.WriteFileSourceOutput{
				URI: v.Source.URI.ValueString(),
			}

			if !v.Source.Headers.IsUnknown() {
				config := make(map[string]string)

				diagnostics := v.Source.Headers.ElementsAs(ctx, &config, false)
				if diagnostics.HasError() {
					return diagnostics
				}

				src.Headers = &config
			}

			item.Source = &src
		}

		writeFiles[k] = item
	}

	output.WriteFiles = &writeFiles
	return nil
}

func transformZypper(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.Zypper == nil {
		return nil
	}

	zypper := ccmodules.ZypperOutput{}

	if !model.Zypper.Repos.IsUnknown() {
		res, diagnostics := castArray[ccmodules.ZypperRepository](ctx, model.Zypper.Repos)

		if diagnostics.HasError() {
			return diagnostics
		}

		repos := make([]ccmodules.ZypperRepositoryOutput, len(*res))
		for k, v := range *res {
			repos[k] = ccmodules.ZypperRepositoryOutput{
				ID:      v.ID.ValueString(),
				BaseURL: v.BaseURL.ValueString(),
			}
		}

		zypper.Repos = &repos
	}

	if !model.Zypper.Config.IsUnknown() {
		config := make(map[string]string)

		diagnostics := model.Zypper.Config.ElementsAs(ctx, &config, false)
		if diagnostics.HasError() {
			return diagnostics
		}

		zypper.Config = &config
	}

	output.Zypper = &zypper

	return nil
}

func transformWireguard(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.Wireguard == nil {
		return nil
	}

	wireguard := ccmodules.WireguardOutput{}

	if !model.Wireguard.ReadinessProbe.IsUnknown() {
		res, diagnostics := castArray[string](ctx, model.Wireguard.ReadinessProbe)
		if diagnostics.HasError() {
			return diagnostics
		}
		wireguard.ReadinessProbe = res
	}

	if !model.Wireguard.Interfaces.IsUnknown() {
		res, diagnostics := castArray[ccmodules.Interface](ctx, model.Wireguard.Interfaces)
		if diagnostics.HasError() {
			return diagnostics
		}
		interfaces := make([]ccmodules.InterfaceOutput, len(*res))
		for k, v := range *res {
			interfaces[k] = ccmodules.InterfaceOutput{
				Name:       v.Name.ValueString(),
				ConfigPath: v.ConfigPath.ValueString(),
				Content:    v.Content.ValueString(),
			}
		}

		wireguard.Interfaces = &interfaces
	}

	output.Wireguard = &wireguard

	return nil
}

func transformRPI(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.RPI == nil {
		return nil
	}

	rpi := ccmodules.RPIOutput{}

	rpi.EnableRPIConnect = model.RPI.EnableRPIConnect.ValueBool()

	if model.RPI.Interfaces != nil {
		interfaces := ccmodules.RPIInterfaceOutput{}

		interfaces.SPI = model.RPI.Interfaces.SPI.ValueBool()
		interfaces.I2C = model.RPI.Interfaces.I2C.ValueBool()
		interfaces.SSH = model.RPI.Interfaces.SSH.ValueBool()
		interfaces.Onewire = model.RPI.Interfaces.Onewire.ValueBool()
		interfaces.RemoteGPIO = model.RPI.Interfaces.RemoteGPIO.ValueBool()

		if model.RPI.Interfaces.Serial != nil {
			serial := ccmodules.RPISerialOutput{}

			serial.Console = model.RPI.Interfaces.Serial.Console.ValueBool()
			serial.Hardware = model.RPI.Interfaces.Serial.Hardware.ValueBool()

			interfaces.Serial = &serial
		}

		rpi.Interfaces = &interfaces
	}

	output.RPI = &rpi

	return nil
}

func transformSeedRandom(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.RandomSeed == nil {
		return nil
	}

	seed := ccmodules.RandomSeedOutput{}

	if !model.RandomSeed.Command.IsUnknown() {
		res, diagnostics := castArray[string](ctx, model.RandomSeed.Command)
		if diagnostics.HasError() {
			return diagnostics
		}
		seed.Command = res
	}

	seed.File = model.RandomSeed.File.ValueString()
	seed.Data = model.RandomSeed.Data.ValueString()
	seed.Encoding = model.RandomSeed.Encoding.ValueString()

	seed.CommandRequired = model.RandomSeed.CommandRequired.ValueBool()

	output.RandomSeed = &seed

	return nil
}

func transformNTP(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.NTP == nil {
		return nil
	}

	ntp := ccmodules.NTPOutput{}

	if !model.NTP.Pools.IsUnknown() {
		res, diagnostics := castArray[string](ctx, model.NTP.Pools)
		if diagnostics.HasError() {
			return diagnostics
		}
		ntp.Pools = res
	}

	if !model.NTP.Servers.IsUnknown() {
		res, diagnostics := castArray[string](ctx, model.NTP.Servers)
		if diagnostics.HasError() {
			return diagnostics
		}
		ntp.Servers = res
	}

	if !model.NTP.Peers.IsUnknown() {
		res, diagnostics := castArray[string](ctx, model.NTP.Peers)
		if diagnostics.HasError() {
			return diagnostics
		}
		ntp.Peers = res
	}

	if !model.NTP.Allow.IsUnknown() {
		res, diagnostics := castArray[string](ctx, model.NTP.Allow)
		if diagnostics.HasError() {
			return diagnostics
		}
		ntp.Allow = res
	}

	ntp.NTPClient = model.NTP.NTPClient.ValueString()
	ntp.Enabled = model.NTP.Enabled.ValueBoolPointer()

	if model.NTP.Config != nil {
		config := ccmodules.NTPConfigOutput{}

		if !model.NTP.Config.Packages.IsUnknown() {
			res, diagnostics := castArray[string](ctx, model.NTP.Config.Packages)
			if diagnostics.HasError() {
				return diagnostics
			}
			config.Packages = res
		}

		config.Confpath = model.NTP.Config.Confpath.ValueString()
		config.CheckExe = model.NTP.Config.CheckExe.ValueString()
		config.ServiceName = model.NTP.Config.ServiceName.ValueString()
		config.Template = model.NTP.Config.Template.ValueString()

		ntp.Config = &config
	}

	output.NTP = &ntp

	return nil
}

func transformLandscape(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.Landscape == nil {
		return nil
	}

	config := ccmodules.LandscapeOutput{}

	if model.Landscape.Client == nil {
		return nil
	}

	client := ccmodules.ClientOutput{}

	client.URL = model.Landscape.Client.URL.ValueString()
	client.PingURL = model.Landscape.Client.PingURL.ValueString()
	client.DataPath = model.Landscape.Client.DataPath.ValueString()
	client.LogLevel = model.Landscape.Client.LogLevel.ValueString()
	client.ComputerTitle = model.Landscape.Client.ComputerTitle.ValueString()
	client.AccountName = model.Landscape.Client.AccountName.ValueString()
	client.RegistrationKey = model.Landscape.Client.RegistrationKey.ValueString()
	client.Tags = model.Landscape.Client.Tags.ValueString()
	client.HTTPProxy = model.Landscape.Client.HTTPProxy.ValueString()
	client.HTTPSProxy = model.Landscape.Client.HTTPSProxy.ValueString()

	config.Client = &client
	output.Landscape = &config

	return nil
}

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
		res, diagnostics := castArray[string](ctx, model.RunCMD)
		if diagnostics.HasError() {
			return diagnostics
		}
		output.RunCMD = res
	}
	return nil
}

func transformBootCMD(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if !model.BootCMD.IsUnknown() {
		res, diagnostics := castArray[string](ctx, model.BootCMD)
		if diagnostics.HasError() {
			return diagnostics
		}
		output.BootCMD = res
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
		res, diagnostics := castArray[string](ctx, model.SSHAuthorizedKeys)
		if diagnostics.HasError() {
			return diagnostics
		}
		output.SSHAuthorizedKeys = res
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
		res, diagnostics := castArray[string](ctx, model.Packages)
		if diagnostics.HasError() {
			return diagnostics
		}
		output.Packages = res
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
			res, diagnostics := castArray[string](ctx, user.Doas)
			if diagnostics.HasError() {
				return out, diagnostics
			}
			out.Doas = res
		}

		if !user.SSHAuthorizedKeys.IsUnknown() {
			res, diagnostics := castArray[string](ctx, user.SSHAuthorizedKeys)
			if diagnostics.HasError() {
				return out, diagnostics
			}
			out.SSHAuthorizedKeys = res
		}

		if !user.SSHImportId.IsUnknown() {
			res, diagnostics := castArray[string](ctx, user.SSHImportId)
			if diagnostics.HasError() {
				return out, diagnostics
			}
			out.SSHImportId = res
		}

		if !user.Sudo.IsUnknown() {
			res, diagnostics := castArray[string](ctx, user.Sudo)
			if diagnostics.HasError() {
				return out, diagnostics
			}
			out.Sudo = res
		}

		if !user.Groups.IsUnknown() {
			res, diagnostics := castArray[string](ctx, user.Groups)
			if diagnostics.HasError() {
				return out, diagnostics
			}
			out.Groups = res
		}

		return out, nil
	}

	if !model.Groups.IsUnknown() {
		//}
		res, diagnostics := castArray[string](ctx, model.Groups)
		if diagnostics.HasError() {
			return diagnostics
		}
		output.Groups = res
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
		res, diagnostics := castArray[string](ctx, model.CACerts.Trusted)
		if diagnostics.HasError() {
			return diagnostics
		}

		caCerts.Trusted = res
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
		res, diagnostics := castArray[string](ctx, model.Growpart.Devices)
		if diagnostics.HasError() {
			return diagnostics
		}
		growpart.Devices = res
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
		// elems := model.Updates.Network.When.Elements()
		//
		// if len(elems) > 0 {
		// 	whens := make([]string, len(elems))
		// 	diagnostics := model.Updates.Network.When.ElementsAs(ctx, &whens, false)
		//
		// 	if diagnostics.HasError() {
		// 		return diagnostics
		// 	}
		//
		// 	config.Network.When = &whens
		// 	output.Updates = &config
		// }

		res, diagnostics := castArray[string](ctx, model.Updates.Network.When)
		if diagnostics.HasError() {
			return diagnostics
		}
		config.Network.When = res
		output.Updates = &config
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
		// elems := model.SSHKeyConsoleBlacklist.Elements()
		//
		// if len(elems) > 0 {
		// 	whens := make([]string, len(elems))
		// 	diagnostics := model.SSHKeyConsoleBlacklist.ElementsAs(ctx, &whens, false)
		//
		// 	if diagnostics.HasError() {
		// 		return diagnostics
		// 	}
		//
		// 	output.SSHKeyConsoleBlacklist = &whens
		// }

		res, diagnostics := castArray[string](ctx, model.SSHKeyConsoleBlacklist)
		if diagnostics.HasError() {
			return diagnostics
		}
		output.SSHKeyConsoleBlacklist = res
	}

	if !model.SSHFPConsoleBlacklist.IsUnknown() {
		// elems := model.SSHFPConsoleBlacklist.Elements()
		//
		// if len(elems) > 0 {
		// 	whens := make([]string, len(elems))
		// 	diagnostics := model.SSHFPConsoleBlacklist.ElementsAs(ctx, &whens, false)
		//
		// 	if diagnostics.HasError() {
		// 		return diagnostics
		// 	}
		//
		// 	output.SSHFPConsoleBlacklist = &whens
		// }

		res, diagnostics := castArray[string](ctx, model.SSHFPConsoleBlacklist)
		if diagnostics.HasError() {
			return diagnostics
		}
		output.SSHFPConsoleBlacklist = res
	}

	return nil
}

func transformResizefs(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	// NOTE: default value is anyway `true`
	if !model.Resizefs.IsNull() && !model.Resizefs.ValueBool() {
		output.Resizefs = model.Resizefs.ValueBoolPointer()
	} else if model.ResizefsNoBlock.ValueBool() {
		output.Resizefs = "noblock"
	} else {
		output.Resizefs = nil
	}

	return nil
}

func transformSaltMinion(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.SaltMinion == nil {
		return nil
	}

	config := ccmodules.SaltMinionOutput{}

	config.PkgName = model.SaltMinion.PkgName.ValueString()
	config.ServiceName = model.SaltMinion.ServiceName.ValueString()
	config.ConfigDir = model.SaltMinion.ConfigDir.ValueString()
	config.PublicKey = model.SaltMinion.PublicKey.ValueString()
	config.PrivateKey = model.SaltMinion.PrivateKey.ValueString()
	config.PkiDir = model.SaltMinion.PkiDir.ValueString()

	output.SaltMinion = &config

	return nil
}

func transformUbuntuAutoinstall(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.Autoinstall == nil {
		return nil
	}

	config := ccmodules.AutoinstallOutput{}

	config.Version = model.Autoinstall.Version.ValueInt32()

	output.Autoinstall = &config

	return nil
}

func transformPowerStateChange(_ context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.PowerState == nil {
		return nil
	}

	config := ccmodules.PowerStateOutput{}

	config.Mode = model.PowerState.Mode.ValueString()
	config.Message = model.PowerState.Message.ValueString()

	config.Timeout = model.PowerState.Timeout.ValueInt64()

	if model.PowerState.NoDelay.ValueBool() {
		// If `no_delay` is true - value is `now`
		config.Delay = "now"
	} else if !model.PowerState.Delay.IsUnknown() && !model.PowerState.Delay.IsNull() {
		config.Delay = model.PowerState.Delay.ValueInt64()
	}

	if !model.PowerState.ConditionCmd.IsUnknown() && !model.PowerState.ConditionCmd.IsNull() {
		config.Condition = model.PowerState.ConditionCmd.ValueString()
	} else if !model.PowerState.Condition.IsUnknown() {
		config.Condition = model.PowerState.Condition.ValueBool()
	}

	output.PowerState = &config

	return nil
}

func transformPhoneHome(ctx context.Context, output *ExportModel, model CloudConfigResourceModel) diag.Diagnostics {
	if model.PhoneHome == nil {
		return nil
	}

	config := ccmodules.PhoneHomeOutput{}

	config.URL = model.PhoneHome.URL.ValueString()
	config.Tries = int(model.PhoneHome.Tries.ValueInt64())

	if !model.PhoneHome.Post.IsUnknown() {
		res, diagnostics := castArray[string](ctx, model.PhoneHome.Post)
		if diagnostics.HasError() {
			return diagnostics
		}

		config.Post = res
	}

	output.PhoneHome = &config

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

	diagnostics = transformResizefs(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformSaltMinion(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformUbuntuAutoinstall(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformPowerStateChange(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformPhoneHome(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformLandscape(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformNTP(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformRPI(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformSeedRandom(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformWireguard(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformZypper(ctx, &output, model)
	if diagnostics.HasError() {
		return output, diagnostics
	}

	diagnostics = transformWriteFiles(ctx, &output, model)
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
