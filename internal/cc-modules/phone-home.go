package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type PhoneHome struct {
	URL   types.String `tfsdk:"url"`
	Tries types.Int64  `tfsdk:"tries"`
	Post  types.List   `tfsdk:"post"`
}

type PhoneHomeOutput struct {
	URL   string    `yaml:"url,omitempty"`
	Tries int       `yaml:"tries,omitempty"`
	Post  *[]string `yaml:"post,omitempty"`
}

type PhoneHomeModel struct {
	PhoneHome *PhoneHome `tfsdk:"phone_home"`
}

type PhoneHomeOutputModel struct {
	PhoneHome *PhoneHomeOutput `yaml:"phone_home,omitempty"`
}

// PhoneHomeBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#phone-home
func PhoneHomeBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"phone_home": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("phone_home")),
				},
				MarkdownDescription: `
This module can be used to post data to a remote host after boot is complete.

Either all data can be posted, or a list of keys to post.

Available keys are:

- pub_key_rsa
- pub_key_ecdsa
- pub_key_ed25519
- instance_id
- hostname
- fqdn

Data is sent as x-www-form-urlencoded arguments.
        `,
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						MarkdownDescription: "The URL to send the phone home data to.",
						Optional:            true,
					},
					"tries": schema.Int64Attribute{
						MarkdownDescription: "The number of times to try sending the phone home data. Default: `10`.",
						Optional:            true,
					},
					"post": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "A list of keys to post or all. Default: `all`.",
						Optional:            true,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(
								stringvalidator.OneOf(
									"pub_key_rsa",
									"pub_key_ecdsa",
									"pub_key_ed25519",
									"instance_id",
									"hostname",
									"fqdn",
									"all",
								),
							),
						},
					},
				},
			},
		},
	}
}
