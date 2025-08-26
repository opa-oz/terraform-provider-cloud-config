package provider

import (
	"context"
	"maps"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	ccmodules "github.com/opa-oz/terraform-provider-cloud-config/internal/cc-modules"
)

var _ resource.Resource = &CloudConfigResource{}

// var _ resource.ResourceWithImportState = &CloudConfigResource{}

func NewCloudConfigResource() resource.Resource {
	return &CloudConfigResource{}
}

type CloudConfigResource struct {
}

type CloudConfigResourceModel struct {
	Content types.String `tfsdk:"content"`

	ccmodules.SetHostnameModel
	ccmodules.LocaleModel
	ccmodules.TimezoneModel
	ccmodules.RunCMDModule
	ccmodules.BootCMDModule
	ccmodules.UpdateEtcHostsModule
	ccmodules.SSHModel
	ccmodules.SetPasswordsModel
	ccmodules.PkgUpdateUpgradeModel
	ccmodules.UsersAndGroupsModel
	ccmodules.DisableEC2InstanceMetadataModel
	ccmodules.ApkConfigureModel
	ccmodules.AptPipeliningModel
	ccmodules.ByobuModel
	ccmodules.CACertificatesModel
	ccmodules.FanModel
	ccmodules.FinalMessageModel
	ccmodules.GrowpartModel
	ccmodules.GRUBDpkgModel
	ccmodules.InstallHotplugModel
	ccmodules.KeyboardModel
	ccmodules.KeysToConsoleModel
	ccmodules.ResizefsModel
	ccmodules.SaltMinionModel
	ccmodules.UbuntuAutoinstallModel
	ccmodules.PowerStateModel
	ccmodules.PhoneHomeModel
	ccmodules.LandscapeModel
	ccmodules.NTPModel
}

type ExportModel struct {
	ccmodules.SetHostnameOutputModel                `yaml:",inline"`
	ccmodules.LocaleOutputModel                     `yaml:",inline"`
	ccmodules.TimezoneOutputModel                   `yaml:",inline"`
	ccmodules.RunCMDOutputModule                    `yaml:",inline"`
	ccmodules.BootCMDOutputModule                   `yaml:",inline"`
	ccmodules.UpdateEtcHostsOutputModule            `yaml:",inline"`
	ccmodules.SSHOutputModel                        `yaml:",inline"`
	ccmodules.SetPasswordsOutputModel               `yaml:",inline"`
	ccmodules.PkgUpdateUpgradeOutputModel           `yaml:",inline"`
	ccmodules.UsersAndGroupsOutputModel             `yaml:",inline"`
	ccmodules.DisableEC2InstanceMetadataOutputModel `yaml:",inline"`
	ccmodules.ApkConfigureOutputModel               `yaml:",inline"`
	ccmodules.AptPipeliningOutputModel              `yaml:",inline"`
	ccmodules.ByobuOutputModel                      `yaml:",inline"`
	ccmodules.CACertificatesOutputModel             `yaml:",inline"`
	ccmodules.FanOutputModel                        `yaml:",inline"`
	ccmodules.FinalMessageOutputModel               `yaml:",inline"`
	ccmodules.GrowpartOutputModel                   `yaml:",inline"`
	ccmodules.GRUBDpkgOutputModel                   `yaml:",inline"`
	ccmodules.InstallHotplugOutputModel             `yaml:",inline"`
	ccmodules.KeyboardOutputModel                   `yaml:",inline"`
	ccmodules.KeysToConsoleOutputModel              `yaml:",inline"`
	ccmodules.ResizefsOutputModel                   `yaml:",inline"`
	ccmodules.SaltMinionOutputModel                 `yaml:",inline"`
	ccmodules.UbuntuAutoinstallOutputModel          `yaml:",inline"`
	ccmodules.PowerStateOutputModel                 `yaml:",inline"`
	ccmodules.PhoneHomeOutputModel                  `yaml:",inline"`
	ccmodules.LandscapeOutputModel                  `yaml:",inline"`
	ccmodules.NTPOutputModel                        `yaml:",inline"`
}

func (r *CloudConfigResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName
}

func (r *CloudConfigResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		MarkdownDescription: "Cloud-config file in-memory representation", // NOTE: https://github.com/nobbs/terraform-provider-sops/blob/main/internal/provider/file_function.go

		Attributes: map[string]schema.Attribute{
			"content": schema.StringAttribute{
				Computed:            true,
				Sensitive:           true,
				MarkdownDescription: "YAML content of cloud-init file",
			},
		},
		Blocks: map[string]schema.Block{},
	}

	flat_modules := []ccmodules.CCModuleFlat{
		ccmodules.SetHostname(),
		ccmodules.Locale(),
		ccmodules.Timezone(),
		ccmodules.RunCMD(),
		ccmodules.BootCMD(),
		ccmodules.UpdateEtcHosts(),
		ccmodules.SSH(),
		ccmodules.SetPasswords(),
		ccmodules.PkgUpdateUpgrade(),
		ccmodules.UsersAndGroups(),
		ccmodules.DisableEC2InstanceMetadata(),
		ccmodules.Byobu(),
		ccmodules.FinalMessage(),
		ccmodules.KeysToConsole(),
		ccmodules.Resizefs(),
	}

	for _, module := range flat_modules {
		maps.Insert(schema.Attributes, maps.All(module.Attributes()))
	}

	block_modules := []ccmodules.CCModuleNested{
		ccmodules.SetPasswordsBlock(),
		ccmodules.UsersAndGroupsBlock(),
		ccmodules.ApkConfigureBlock(),
		ccmodules.AptPipeliningBlock(),
		ccmodules.CACertificatesBlock(),
		ccmodules.FanBlock(),
		ccmodules.GrowpartBlock(),
		ccmodules.GRUBDpkgBlock(),
		ccmodules.InstallHotplugBlock(),
		ccmodules.KeyboardBlock(),
		ccmodules.KeysToConsoleBlock(),
		ccmodules.SaltMinionBlock(),
		ccmodules.UbuntuAutoinstallBlock(),
		ccmodules.PowerStateChangeBlock(),
		ccmodules.PhoneHomeBlock(),
		ccmodules.LandscapeBlock(),
		ccmodules.NTPBlock(),
	}

	for _, module := range block_modules {
		maps.Insert(schema.Blocks, maps.All(module.Block()))
	}

	resp.Schema = schema
}

func (r *CloudConfigResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
}

func (r *CloudConfigResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudConfigResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	content, err := ExportContent(ctx, data)
	if err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	data.Content = types.StringValue(content)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudConfigResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudConfigResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudConfigResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data CloudConfigResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	content, err := ExportContent(ctx, data)
	if err != nil {
		resp.Diagnostics.Append(err...)
		return
	}

	data.Content = types.StringValue(content)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CloudConfigResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudConfigResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *CloudConfigResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("manage_etc_hosts"),
			path.MatchRoot("manage_etc_hosts_localhost"),
		),
	}
}

// func (r *CloudConfigResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
// }
