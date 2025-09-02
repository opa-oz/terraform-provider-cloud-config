package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WriteFileSource struct {
	URI     types.String `tfsdk:"uri"`
	Headers types.Map    `tfsdk:"headers"`
}

type WriteFileSourceOutput struct {
	URI     string             `yaml:"uri,omitempty"`
	Headers *map[string]string `yaml:"headers,omitempty"`
}

type WriteFile struct {
	Path        types.String     `tfsdk:"path"`
	Content     types.String     `tfsdk:"content"`
	Owner       types.String     `tfsdk:"owner"`
	Permissions types.String     `tfsdk:"permissions"`
	Encoding    types.String     `tfsdk:"encoding"`
	Append      types.Bool       `tfsdk:"append"`
	Defer       types.Bool       `tfsdk:"defer"`
	Source      *WriteFileSource `tfsdk:"source"`
}

type WriteFileOutput struct {
	Path        string                 `yaml:"path,omitempty"`
	Content     string                 `yaml:"content,omitempty"`
	Owner       string                 `yaml:"owner,omitempty"`
	Permissions string                 `yaml:"permissions,omitempty"`
	Encoding    string                 `yaml:"encoding,omitempty"`
	Append      bool                   `yaml:"append,omitempty"`
	Defer       bool                   `yaml:"defer,omitempty"`
	Source      *WriteFileSourceOutput `yaml:"source,omitempty"`
}

type WriteFileModel struct {
	WriteFiles types.List `tfsdk:"write_files"`
}

type WriteFileOutputModel struct {
	WriteFiles *[]WriteFileOutput `yaml:"write_files,omitempty"`
}

// WriteFileBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#write-files
func WriteFileBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"write_files": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"source": schema.SingleNestedBlock{
							// PlanModifiers: []planmodifier.Object{
							// 	utils.NullWhen(path.Root("write_files").AtName("source")),
							// },
							MarkdownDescription: "Optional specification for content loading from an arbitrary URI",
							Attributes: map[string]schema.Attribute{
								"uri": schema.StringAttribute{
									MarkdownDescription: "URI from which to load file content. If loading fails repeatedly, content is used instead.",
									Optional:            true,
								},
								"headers": schema.MapAttribute{
									ElementType:         types.StringType,
									MarkdownDescription: "Optional HTTP headers to accompany load request, if applicable.",
									Optional:            true,
								},
							},
						},
					},
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							MarkdownDescription: "Path of the file to which **content** is decoded and written.",
							Optional:            true,
						},
						"content": schema.StringAttribute{
							MarkdownDescription: "Optional content to write to the provided path. When content is present and encoding is not ‘text/plain’, decode the content prior to writing. Default: ''.",
							Optional:            true,
						},
						"owner": schema.StringAttribute{
							MarkdownDescription: "Optional owner:group to chown on the file and new directories. Default: `root:root`.",
							Optional:            true,
						},
						"permissions": schema.StringAttribute{
							MarkdownDescription: "Optional file permissions to set on path represented as an octal string ‘0###’. Default: `0o644`.",
							Optional:            true,
						},
						"encoding": schema.StringAttribute{
							MarkdownDescription: "Optional encoding type of the content. Default: text/plain. No decoding is performed by default. Supported encoding types are: `gz, gzip, gz+base64, gzip+base64, gz+b64, gzip+b64, b64, base64`.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf(
									"gz",
									"gzip",
									"gz+base64",
									"gzip+base64",
									"gz+b64",
									"gzip+b64",
									"b64",
									"base64",
									"text/plain",
								),
							},
						},
						"append": schema.BoolAttribute{
							MarkdownDescription: "Whether to append content to existing file if path exists. Default: false.",
							Optional:            true,
						},
						"defer": schema.BoolAttribute{
							MarkdownDescription: "Defer writing the file until ‘final’ stage, after users were created, and packages were installed. Default: false.",
							Optional:            true,
						},
					},
				},
			},
		},
	}
}
