package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AptPipeliningConfig struct {
	OS      types.Bool  `tfsdk:"os"`
	Disable types.Bool  `tfsdk:"disable"`
	Depth   types.Int32 `tfsdk:"depth"`
}

type AptPipeliningModel struct {
	AptPipelining *AptPipeliningConfig `tfsdk:"apt_pipelining"`
}

type AptPipeliningOutputModel struct {
	AptPipelining any `yaml:"apt_pipelining,omitempty"`
}

// func (self AptPipeliningOutputModel) MarshalYAML() (any, error) {
// 	if self.AptPipeliningStr == "os" {
// 		return struct {
// 			val string `yaml:"apt_pipelining"`
// 		}{val: "os"}, nil
// 	}
//
// 	if self.AptPipeliningBool {
// 		return "true", nil
// 	}
//
// 	if self.AptPipeliningNumber != nil {
// 		return self.AptPipeliningNumber, nil
// 	}
//
// 	return nil, nil
// }

// AptPipeliningBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#apt-pipelining
func AptPipeliningBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"apt_pipelining": schema.SingleNestedBlock{
				MarkdownDescription: "This module configures APT’s `Acquire::http::Pipeline-Depth` option, which controls how APT handles HTTP pipelining. It may be useful for pipelining to be disabled, because some web servers (such as S3) do not pipeline properly (LP: #948461).",
				Attributes: map[string]schema.Attribute{
					"os": schema.BoolAttribute{
						MarkdownDescription: "Use distro default. This is default behaivor",
						Optional:            true,
						Validators: []validator.Bool{
							boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("disable")),
							boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("depth")),
						},
					},
					"disable": schema.BoolAttribute{
						MarkdownDescription: "Disable pipelining altogether",
						Optional:            true,
						Validators: []validator.Bool{
							boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("os")),
							boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("depth")),
						},
					},
					"depth": schema.Int32Attribute{
						MarkdownDescription: "Manually specify pipeline depth. This is not recommended.",
						Optional:            true,
						Validators: []validator.Int32{
							int32validator.ConflictsWith(path.MatchRelative().AtParent().AtName("os")),
							int32validator.ConflictsWith(path.MatchRelative().AtParent().AtName("disable")),
						},
					},
				},
			},
		},
	}
}
