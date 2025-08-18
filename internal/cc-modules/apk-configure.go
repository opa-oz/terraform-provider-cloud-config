package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type AlpineRepo struct {
	CommunityEnabled types.Bool   `tfsdk:"community_enabled"`
	TestingEnabled   types.Bool   `tfsdk:"testing_enabled"`
	BaseUrl          types.String `tfsdk:"base_url"`
	Version          types.String `tfsdk:"version"`
}

type AlpineRepoOutput struct {
	CommunityEnabled bool   `yaml:"community_enabled,omitempty"`
	TestingEnabled   bool   `yaml:"testing_enabled,omitempty"`
	BaseUrl          string `yaml:"base_url,omitempty"`
	Version          string `yaml:"version,omitempty"`
}

type ApkRepo struct {
	PreserveRepositories types.Bool   `tfsdk:"preserve_repositories"`
	LocalRepoBaseUrl     types.String `tfsdk:"local_repo_base_url"`
	AlpineRepo           *AlpineRepo  `tfsdk:"alpine_repo"`
}

type ApkRepoOutput struct {
	PreserveRepositories bool              `yaml:"preserve_repositories,omitempty"`
	LocalRepoBaseUrl     string            `yaml:"local_repo_base_url,omitempty"`
	AlpineRepo           *AlpineRepoOutput `yaml:"alpine_repo,omitempty"`
}

type ApkConfigureModel struct {
	ApkRepos *ApkRepo `tfsdk:"apk_repos"`
}

type ApkConfigureOutputModel struct {
	ApkRepos *ApkRepoOutput `yaml:"apk_repos,omitempty"`
}

// ApkConfigureBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#apk-configure
func ApkConfigureBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"apk_repos": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("apk_repos")),
				},
				MarkdownDescription: "This module handles configuration of the Alpine Package Keeper (APK) /etc/apk/repositories file.",
				Attributes: map[string]schema.Attribute{
					"preserve_repositories": schema.BoolAttribute{
						MarkdownDescription: `
By default, cloud-init will generate a new repositories file /etc/apk/repositories based on any valid configuration settings specified within a apk_repos section of cloud config. To disable this behavior and preserve the repositories file from the pristine image, set **preserve_repositories** to true.
The **preserve_repositories** option overrides all other config keys that would alter /etc/apk/repositories.
            `,
						Optional: true,
					},
					"local_repo_base_url": schema.StringAttribute{
						MarkdownDescription: "The base URL of an Alpine repository containing unofficial packages.",
						Optional:            true,
					},
				},
				Blocks: map[string]schema.Block{
					"alpine_repo": schema.SingleNestedBlock{
						Attributes: map[string]schema.Attribute{
							"community_enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether to add the Community repo to the repositories file. By default the Community repo is not included.",
								Optional:            true,
							},
							"testing_enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether to add the Testing repo to the repositories file. By default the Testing repo is not included. It is only recommended to use the Testing repo on a machine running the Edge version of Alpine as packages installed from Testing may have dependencies that conflict with those in non-Edge Main or Community repos.",
								Optional:            true,
							},
							"base_url": schema.StringAttribute{
								MarkdownDescription: "The base URL of an Alpine repository, or mirror, to download official packages from. If not specified then it defaults to https://alpine.global.ssl.fastly.net/alpine.",
								Optional:            true,
							},
							"version": schema.StringAttribute{
								MarkdownDescription: "version",
								Optional:            true,
							},
						},
					},
				},
			},
		},
	}
}
