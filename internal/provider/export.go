package provider

import (
	"fmt"

	ccmodules "github.com/opa-oz/terraform-provider-cloud-config/internal/cc-modules"
	"gopkg.in/yaml.v3"
)

const (
	hat = "#cloud-config"
)

type ExportModel struct {
	ccmodules.SetHostnameOutputModel `yaml:",inline"`
}

func transform(model CloudConfigResourceModel) ExportModel {
	output := ExportModel{}

	output.Hostname = model.Hostname.ValueString()
	output.FQDN = model.FQDN.ValueString()
	output.PreserveHostname = model.PreserveHostname.ValueBool()
	output.PreferFQDNOverHostname = model.PreferFQDNOverHostname.ValueBool()

	// NOTE: default value is anyway `true`
	if !model.CreateHostnameFile.ValueBool() {
		output.CreateHostnameFile = model.CreateHostnameFile.ValueBool()
	}

	return output
}

func ExportContent(model CloudConfigResourceModel) (string, error) {
	yaml, err := yaml.Marshal(transform(model))

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(`%s
%s
  `, hat, yaml), nil
}
