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

func transform(ctx context.Context, model CloudConfigResourceModel) (ExportModel, diag.Diagnostics) {
	output := ExportModel{}

	output.Hostname = model.Hostname.ValueString()
	output.FQDN = model.FQDN.ValueString()
	output.PreserveHostname = model.PreserveHostname.ValueBool()
	output.PreferFQDNOverHostname = model.PreferFQDNOverHostname.ValueBool()

	// NOTE: default value is anyway `true`
	if !model.CreateHostnameFile.IsNull() && !model.CreateHostnameFile.ValueBool() {
		output.CreateHostnameFile = model.CreateHostnameFile.ValueBoolPointer()
	}

	if !model.Locale.IsUnknown() {
		output.Locale = model.Locale.ValueString()
	}
	if !model.LocaleConfigfile.IsUnknown() {
		output.LocaleConfigfile = model.LocaleConfigfile.ValueString()
	}

	output.Timezone = model.Timezone.ValueString()

	if !model.RunCMD.IsUnknown() {
		elems := model.RunCMD.Elements()

		if len(elems) > 0 {
			cmds := make([]string, len(elems))
			diagnostics := model.RunCMD.ElementsAs(ctx, &cmds, false)

			if diagnostics.HasError() {
				return output, diagnostics
			}
			output.RunCMD = cmds
		}
	}

	if !model.ManageEtcHostsLocalhost.IsNull() && model.ManageEtcHostsLocalhost.ValueBool() {
		output.ManageEtcHosts = "localhost"
	} else if !model.ManageEtcHosts.IsNull() {
		output.ManageEtcHosts = model.ManageEtcHosts.ValueBool()
	}

	if !model.SSHAuthorizedKeys.IsUnknown() {
		elems := model.SSHAuthorizedKeys.Elements()

		if len(elems) > 0 {
			sshAuthKeys := make([]string, len(elems))
			diagnostics := model.SSHAuthorizedKeys.ElementsAs(ctx, &sshAuthKeys, false)

			if diagnostics.HasError() {
				return output, diagnostics
			}
			output.SSHAuthorizedKeys = sshAuthKeys
		}
	}

	if !model.SSHPwauth.IsNull() {
		output.SSHPwauth = model.SSHPwauth.ValueBoolPointer()
	}

	if model.ChPasswd != nil {
		output.ChPasswd = &ccmodules.ChPasswdOutput{}

		md := model.ChPasswd

		if !md.Expire.IsNull() {
			output.ChPasswd.Expire = md.Expire.ValueBoolPointer()
		}

		if md.Users != nil {
			usrs := make([]ccmodules.ChPasswdUserOutput, len(*md.Users))
			for i, usr := range *md.Users {
				newUsr := ccmodules.ChPasswdUserOutput{}

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
