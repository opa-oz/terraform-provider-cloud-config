package ccmodules

import (
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/opa-oz/terraform-provider-cloud-config/internal/utils"
)

type Interface struct {
	Name       types.String `tfsdk:"name"`
	ConfigPath types.String `tfsdk:"config_path"`
	Content    types.String `tfsdk:"content"`
}

type Wireguard struct {
	Interfaces     types.List `tfsdk:"interfaces"`
	ReadinessProbe types.List `tfsdk:"readinessprobe"`
}

type InterfaceOutput struct {
	Name       string `yaml:"name,omitempty"`
	ConfigPath string `yaml:"config_path,omitempty"`
	Content    string `yaml:"content,omitempty"`
}

type WireguardOutput struct {
	Interfaces     *[]InterfaceOutput `yaml:"interfaces,omitempty"`
	ReadinessProbe *[]string          `yaml:"readinessprobe,omitempty"`
}

type WireguardModel struct {
	Wireguard *Wireguard `tfsdk:"wireguard"`
}

type WireguardOutputModel struct {
	Wireguard *WireguardOutput `yaml:"wireguard,omitempty"`
}

// WireguardBlock
// @see https://cloudinit.readthedocs.io/en/latest/reference/modules.html#wireguard
func WireguardBlock() CCModuleNested {
	return CCModuleNested{
		block: map[string]schema.Block{
			"wireguard": schema.SingleNestedBlock{
				PlanModifiers: []planmodifier.Object{
					utils.NullWhen(path.Root("wireguard")),
				},
				MarkdownDescription: `
The WireGuard module provides a dynamic interface for configuring WireGuard (as a peer or server) in a straightforward way.

This module takes care of:
 - writing interface configuration files
 - enabling and starting interfaces
 - installing wireguard-tools package
 - loading WireGuard kernel module
 - executing readiness probes

**What is a readiness probe?**
The idea behind readiness probes is to ensure WireGuard connectivity before continuing the cloud-init process. This could be useful if you need access to specific services like an internal APT Repository Server (e.g., Landscape) to install/update packages.

**Example**
An edge device can’t access the internet but uses cloud-init modules which will install packages (e.g. landscape, packages, ubuntu_advantage). Those modules will fail due to missing internet connection. The wireguard module fixes that problem as it waits until all readiness probes (which can be arbitrary commands, e.g. checking if a proxy server is reachable over WireGuard network) are finished, before continuing the cloud-init config stage.

>In order to use DNS with WireGuard you have to install the resolvconf package or symlink it to systemd’s resolvectl, otherwise wg-quick commands will throw an error message that executable resolvconf is missing, which leads the wireguard module to fail.
        `,
				Attributes: map[string]schema.Attribute{
					"readinessprobe": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "List of shell commands to be executed as probes.",
						Optional:            true,
					},
				},
				Blocks: map[string]schema.Block{
					"interfaces": schema.ListNestedBlock{
						NestedObject: schema.NestedBlockObject{
							Attributes: map[string]schema.Attribute{
								"name": schema.StringAttribute{
									MarkdownDescription: "Name of the interface. Typically wgx (example: wg0).",
									Optional:            true,
								},
								"config_path": schema.StringAttribute{
									MarkdownDescription: "Path to configuration file of Wireguard interface.",
									Optional:            true,
								},
								"content": schema.StringAttribute{
									MarkdownDescription: "Wireguard interface configuration. Contains key, peer, …",
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
