package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type RPISerial struct {
	Console  types.Bool `tfsdk:"console"`
	Hardware types.Bool `tfsdk:"hardware"`
}

type RPISerialOutput struct {
	Console  bool `yaml:"console,omitempty"`
	Hardware bool `yaml:"hardware,omitempty"`
}

type RPIInterface struct {
	SPI        types.Bool `tfsdk:"spi"`
	I2C        types.Bool `tfsdk:"i2c"`
	SSH        types.Bool `tfsdk:"ssh"`
	Serial     *RPISerial `tfsdk:"serial"`
	Onewire    types.Bool `tfsdk:"onewire"`
	RemoteGPIO types.Bool `tfsdk:"remote_gpio"`
}

type RPIInterfaceOutput struct {
	SPI        bool             `yaml:"spi,omitempty"`
	I2C        bool             `yaml:"i2c,omitempty"`
	SSH        bool             `yaml:"ssh,omitempty"`
	Serial     *RPISerialOutput `yaml:"serial,omitempty"`
	Onewire    bool             `yaml:"onewire,omitempty"`
	RemoteGPIO bool             `yaml:"remote_gpio,omitempty"`
}

type RPI struct {
	Interfaces       *RPIInterface `tfsdk:"interfaces"`
	EnableRPIConnect types.Bool    `tfsdk:"enable_rpi_connect"`
}

type RPIOutput struct {
	Interfaces       *RPIInterfaceOutput `yaml:"interfaces,omitempty"`
	EnableRPIConnect bool                `yaml:"enable_rpi_connect,omitempty"`
}

type RPIModel struct {
	RPI *RPI `tfsdk:"rpi"`
}

type RPIOutputModel struct {
	RPI *RPIOutput `yaml:"rpi,omitempty"`
}

// RPIBlock Raspberry Pi Configuration
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#raspberry-pi-configuration
func RPIBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"rpi": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("rpi")),
				},
				MarkdownDescription: `
This module handles ARM interface configuration for Raspberry Pi.

It also handles Raspberry Pi Connect installation and enablement. Raspberry Pi Connect service will be installed and enabled to auto start on boot.

This only works on Raspberry Pi OS (bookworm and later).
        `,
				Attributes: map[string]schema.Attribute{
					"enable_rpi_connect": schema.BoolAttribute{
						MarkdownDescription: "Install and enable Raspberry Pi Connect. Default: `false`.",
						Optional:            true,
					},
				},
				Blocks: map[string]schema.Block{
					"interfaces": schema.SingleNestedBlock{
						PlanModifiers: []planmodifier.Object{
							utils.NullWhen(path.Root("interfaces")),
						},
						Blocks: map[string]schema.Block{
							"serial": schema.SingleNestedBlock{
								PlanModifiers: []planmodifier.Object{
									utils.NullWhen(path.Root("serial")),
								},
								Attributes: map[string]schema.Attribute{
									"console": schema.BoolAttribute{
										MarkdownDescription: "Enable serial console. Default: `false`.",
										Optional:            true,
									},
									"hardware": schema.BoolAttribute{
										MarkdownDescription: "Enable UART hardware. Default: `false`.",
										Optional:            true,
									},
								},
							},
						},
						Attributes: map[string]schema.Attribute{
							"spi": schema.BoolAttribute{
								MarkdownDescription: "Enable SPI interface. Default: `false`.",
								Optional:            true,
							},
							"i2c": schema.BoolAttribute{
								MarkdownDescription: "Enable I2C interface. Default: `false`.",
								Optional:            true,
							},
							"ssh": schema.BoolAttribute{
								MarkdownDescription: "**NOTE** This is not in `cloud-init` documentation, but in examples. Just for compatibility sake let it be here. Default: `false`.",
								Optional:            true,
							},
							"onewire": schema.BoolAttribute{
								MarkdownDescription: "Enable 1-Wire interface. Default: `false`.",
								Optional:            true,
							},
							"remote_gpio": schema.BoolAttribute{
								MarkdownDescription: "Enable remote GPIO interface. Default: `false`.",
								Optional:            true,
							},
						},
					},
				},
			},
		},
	}
}
