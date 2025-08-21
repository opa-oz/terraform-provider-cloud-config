package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type Keyboard struct {
	Layout  types.String `tfsdk:"layout"`
	Model   types.String `tfsdk:"model"`
	Variant types.String `tfsdk:"variant"`
	Options types.String `tfsdk:"options"`
}

type KeyboardOutput struct {
	Layout  string `yaml:"layout,omitempty"`
	Model   string `yaml:"model,omitempty"`
	Variant string `yaml:"variant,omitempty"`
	Options string `yaml:"options,omitempty"`
}

type KeyboardModel struct {
	Keyboard *Keyboard `tfsdk:"keyboard"`
}

type KeyboardOutputModel struct {
	Keyboard *KeyboardOutput `yaml:"keyboard,omitempty"`
}

// KeyboardBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#keyboard
func KeyboardBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"keyboard": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("keyboard")),
				},
				MarkdownDescription: "Handle keyboard configuration.",
				Attributes: map[string]schema.Attribute{
					"layout": schema.StringAttribute{
						MarkdownDescription: "Required. Keyboard layout. Corresponds to XKBLAYOUT.",
						Optional:            true,
					},
					"model": schema.StringAttribute{
						MarkdownDescription: "Keyboard model. Corresponds to XKBMODEL. Default: `pc105`.",
						Optional:            true,
					},
					"variant": schema.StringAttribute{
						MarkdownDescription: "Required for Alpine Linux, optional otherwise. Keyboard variant. Corresponds to `XKBVARIANT`.",
						Optional:            true,
					},
					"options": schema.StringAttribute{
						MarkdownDescription: "Keyboard options. Corresponds to `XKBOPTIONS`.",
						Optional:            true,
					},
				},
			},
		},
	}
}
