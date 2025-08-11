package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"

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

	output.Locale = model.Locale.ValueString()
	output.LocaleConfigfile = model.LocaleConfigfile.ValueString()

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

	if model.ManageEtcHostsLocalhost.ValueBool() {
		output.ManageEtcHosts = "localhost"
	} else {
		output.ManageEtcHosts = model.ManageEtcHosts.ValueBool()
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
