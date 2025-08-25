package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type Client struct {
	URL             types.String `tfsdk:"url"`
	PingURL         types.String `tfsdk:"ping_url"`
	DataPath        types.String `tfsdk:"data_path"`
	LogLevel        types.String `tfsdk:"log_level"`
	ComputerTitle   types.String `tfsdk:"computer_title"`
	AccountName     types.String `tfsdk:"account_name"`
	RegistrationKey types.String `tfsdk:"registration_key"`
	Tags            types.String `tfsdk:"tags"`
	HTTPProxy       types.String `tfsdk:"http_proxy"`
	HTTPSProxy      types.String `tfsdk:"https_proxy"`
}

type ClientOutput struct {
	URL             string `yaml:"url,omitempty"`
	PingURL         string `yaml:"ping_url,omitempty"`
	DataPath        string `yaml:"data_path,omitempty"`
	LogLevel        string `yaml:"log_level,omitempty"`
	ComputerTitle   string `yaml:"computer_title,omitempty"`
	AccountName     string `yaml:"account_name,omitempty"`
	RegistrationKey string `yaml:"registration_key,omitempty"`
	Tags            string `yaml:"tags,omitempty"`
	HTTPProxy       string `yaml:"http_proxy,omitempty"`
	HTTPSProxy      string `yaml:"https_proxy,omitempty"`
}

type Landscape struct {
	Client *Client `tfsdk:"client"`
}

type LandscapeOutput struct {
	Client *ClientOutput `yaml:"client"`
}

type LandscapeModel struct {
	Landscape *Landscape `tfsdk:"landscape"`
}

type LandscapeOutputModel struct {
	Landscape *LandscapeOutput `yaml:"landscape,omitempty"`
}

// LandscapeBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#landscape
func LandscapeBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"landscape": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("landscape")),
				},
				MarkdownDescription: `
This module installs and configures landscape-client. The Landscape client will only be installed if the key landscape is present in config.

Landscape client configuration is given under the client key under the main landscape config key. The config parameters are not interpreted by cloud-init, but rather are converted into a ConfigObj-formatted file and written out to the [client] section in /etc/landscape/client.conf. The following default client config is provided, but can be overridden

If tags is defined, its contents should be a string delimited with a comma (“,”) rather than a list.
        `,
				Blocks: map[string]schema.Block{
					"client": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"url":       schema.StringAttribute{MarkdownDescription: "The Landscape server URL to connect to. Default: `https://landscape.canonical.com/message-system`.", Optional: true},
							"ping_url":  schema.StringAttribute{MarkdownDescription: "The URL to perform lightweight exchange initiation with. Default: `https://landscape.canonical.com/ping`", Optional: true},
							"data_path": schema.StringAttribute{MarkdownDescription: "The directory to store data files in. Default: `/var/lib/land‐scape/client/.`", Optional: true},
							"log_level": schema.StringAttribute{MarkdownDescription: "The log level for the client. Default: `info`.", Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf(
										"debug",
										"info",
										"warning",
										"error",
										"critical",
									),
								},
							},
							"computer_title":   schema.StringAttribute{MarkdownDescription: "The title of this computer.", Optional: true},
							"account_name":     schema.StringAttribute{MarkdownDescription: "The account this computer belongs to.", Optional: true},
							"registration_key": schema.StringAttribute{MarkdownDescription: "The account-wide key used for registering clients.", Optional: true},
							"tags":             schema.StringAttribute{MarkdownDescription: "Comma separated list of tag names to be sent to the server.", Optional: true},
							"http_proxy":       schema.StringAttribute{MarkdownDescription: "The URL of the HTTP proxy, if one is needed.", Optional: true},
							"https_proxy":      schema.StringAttribute{MarkdownDescription: "The URL of the HTTPS proxy, if one is needed.", Optional: true},
						},
					},
				},
			},
		},
	}
}
