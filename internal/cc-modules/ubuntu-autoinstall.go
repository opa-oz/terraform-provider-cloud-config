package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type Autoinstall struct {
	Version types.Int32 `tfsdk:"version"`
}

type AutoinstallOutput struct {
	Version int32 `yaml:"version,omitempty"`
}

type UbuntuAutoinstallModel struct {
	Autoinstall *Autoinstall `tfsdk:"autoinstall"`
}

type UbuntuAutoinstallOutputModel struct {
	Autoinstall *AutoinstallOutput `yaml:"autoinstall,omitempty"`
}

// UbuntuAutoinstallBlock
func UbuntuAutoinstallBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"autoinstall": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("autoinstall")),
				},
				MarkdownDescription: `
**Cloud-init ignores this key and its values. It is used by Subiquity, the Ubuntu Autoinstaller. See: https://ubuntu.com/server/docs/install/autoinstall-reference.**

Cloud-init is used by the Ubuntu installer in two stages. The autoinstall key may contain a configuration for the Ubuntu installer.

Cloud-init verifies that an autoinstall key contains a version key and that the installer package is present on the system.

The Ubuntu installer might pass part of this configuration to cloud-init during a later boot as part of the install process. See [the Ubuntu installer documentation](https://canonical-subiquity.readthedocs-hosted.com/en/latest/reference/autoinstall-reference.html#user-data) for more information. Please direct Ubuntu installer questions to their IRC channel (#ubuntu-server on Libera).
        `,
				Attributes: map[string]schema.Attribute{
					"version": schema.Int32Attribute{
						Optional: true,
					},
				},
			},
		},
	}
}
