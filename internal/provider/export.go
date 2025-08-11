package provider

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	hat = "#cloud-config"
)

func transform(model CloudConfigResourceModel) ExportModel {
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

	return output
}

func ExportContent(model CloudConfigResourceModel) (string, error) {
	yaml, err := yaml.Marshal(transform(model))

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(fmt.Sprintf(`%s
%s
  `, hat, yaml)), nil
}
