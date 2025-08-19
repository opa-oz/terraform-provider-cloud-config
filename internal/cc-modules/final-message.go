package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FinalMessageModel struct {
	FinalMessage types.String `tfsdk:"final_message"`
}

type FinalMessageOutputModel struct {
	FinalMessage string `yaml:"final_message,omitempty"`
}

// FinalMessage
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#final-message
func FinalMessage() CCModuleFlat {
	return CCModuleFlat{
		attributes: map[string]schema.Attribute{
			"final_message": schema.StringAttribute{
				MarkdownDescription: `
This module configures the final message that cloud-init writes. The message is specified as a Jinja template with the following variables set:

 - **version**: cloud-init version
 - **timestamp**: time at cloud-init finish
 - **datasource**: cloud-init data source
 - **uptime**: system uptime

This message is written to the cloud-init log (usually /var/log/cloud-init.log) as well as stderr (which usually redirects to /var/log/cloud-init-output.log).

Upon exit, this module writes the system uptime, timestamp, and cloud-init version to /var/lib/cloud/instance/boot-finished independent of any user data specified for this module.
        `,
				Optional: true,
			},
		},
	}
}
