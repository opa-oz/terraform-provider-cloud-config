package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type ZypperRepository struct {
	ID      types.String `tfsdk:"id"`
	BaseURL types.String `tfsdk:"baseurl"`
}

type ZypperRepositoryOutput struct {
	ID      string `yaml:"id,omitempty"`
	BaseURL string `yaml:"baseurl,omitempty"`
}

type Zypper struct {
	Repos  types.List `tfsdk:"repos"`
	Config types.Map  `tfsdk:"config"`
}

type ZypperOutput struct {
	Repos  *[]ZypperRepositoryOutput `yaml:"repos,omitempty"`
	Config *map[string]string        `yaml:"config,omitempty"`
}

type ZypperModel struct {
	Zypper *Zypper `tfsdk:"zypper"`
}

type ZypperOutputModel struct {
	Zypper *ZypperOutput `yaml:"zypper,omitempty"`
}

// ZypperBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#zypper-add-repo
func ZypperBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"zypper": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("zypper")),
				},
				MarkdownDescription: `
Zypper behavior can be configured using the config key, which will modify /etc/zypp/zypp.conf. The configuration writer will only append the provided configuration options to the configuration file. Any duplicate options will be resolved by the way the zypp.conf INI file is parsed.

> Setting configdir is not supported and will be skipped.

The repos key may be used to add repositories to the system. Beyond the required id and baseurl attributions, no validation is performed on the repos entries.

It is assumed the user is familiar with the Zypper repository file format. This configuration is also applicable for systems with transactional-updates.
        `,
				Attributes: map[string]schema.Attribute{
					"config": schema.MapAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "Any supported zypo.conf key is written to `/etc/zypp/zypp.conf`.",
						Optional:            true,
					},
				},
				Blocks: map[string]schema.Block{
					"repos": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "The unique id of the repo, used when writing `/etc/zypp/repos.d/<id>.repo`.",
									Optional:            true,
								},
								"baseurl": schema.StringAttribute{
									MarkdownDescription: "The base repositoy URL.",
									Optional:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}
