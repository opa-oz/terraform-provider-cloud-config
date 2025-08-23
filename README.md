# terraform-provider-cloud-config

[![Tests](https://github.com/opa-oz/terraform-provider-cloud-config/actions/workflows/test.yaml/badge.svg)](https://github.com/opa-oz/terraform-provider-cloud-config/actions/workflows/test.yaml)
[![Release](https://img.shields.io/github/v/release/opa-oz/terraform-provider-cloud-config?style=flat-square)](https://github.com/opa-oz/terraform-provider-cloud-config/releases)
[![Terraform registry](https://img.shields.io/badge/terraform-registry-623CE4.svg?style=flat-square&logo=terraform)](https://registry.terraform.io/providers/opa-oz/cloud-config)
[![OpenTofu Registry](https://img.shields.io/badge/opentofu-registry-yellow.svg)](https://search.opentofu.org/provider/opa-oz/cloud-config)

Better than `yamlencode` (IMO)

## Motivation

“I just wanted to build a `#cloud-config` for my Proxmox VMs using Terraform, but `yamlencode` failed me. 

I find it more comfortable (and it’s actually working!) if I do it as a Terraform provider.

Whether `ephemeral` or not is open for debate; let me know in Issues.”

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

The configuration alike:

```hcl
resource "cloud-config" "vm_cloud_config" {
  provider            = cloudconfig
  hostname            = var.name
  manage_etc_hosts    = true
  fqdn                = "${var.name}.lan"
  ssh_authorized_keys = [trimspace(data.local_file.ssh_public_key.content)]
  package_update      = true
  package_upgrade     = true
  timezone            = local.timezone
  packages = [
    "qemu-guest-agent",
    "ufw",
  ]
  runcmd = [
    "systemctl enable qemu-guest-agent",
    "systemctl start qemu-guest-agent",
    "ufw limit proto tcp to any port 22",
    "ufw enable",
    "locale-gen ${local.locale}",
    "localectl set-locale LANG=${local.locale}",
    "chfn -f Ansible ${local.ansible_user}",
  ]

  chpasswd { expire = false }

  user {
    name = local.ansible_user
  }

  users {
    name = "default"
  }
}

# Example with Proxmox provider
resource "proxmox_virtual_environment_file" "user_data_cloud_config" {
  provider     = proxmox.api
  content_type = "snippets"
  datastore_id = "local"
  node_name    = var.node_name

  source_raw {
    data      = cloud-config.vm_cloud_config.content # << `content` is computed output value
    file_name = "user-data-${var.vmid}-cloud-config.yaml"
  }
}

# Example with output
output "cloud-config" {
    value = cloud-config.vm_cloud_config.content # << `content` is computed output value
}
```


Result would look somethign like (templated to match terraform desciprion):
```yaml

#cloud-config
hostname: {{item.name}}
manage_etc_hosts: true
fqdn: {{item.name}}.lan
user: {{ansible_user}}
ssh_authorized_keys:
  - "ssh ..." 
chpasswd:
  expire: False
users:
  - default
package_update: true
package_upgrade: true
timezone: {{timezone}}
packages:
  - qemu-guest-agent
  - ufw
runcmd:
  - systemctl enable qemu-guest-agent
  - systemctl start qemu-guest-agent
  - ufw limit proto tcp to any port 22
  - ufw enable
  - locale-gen {{locale}}
  - localectl set-locale LANG={{locale}}
  - chfn -f Ansible {{ansible_user}}

```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Module support

CloudInit has a lot of modules ([https://cloudinit.readthedocs.io/en/latest/reference/modules.html#module-reference](https://cloudinit.readthedocs.io/en/latest/reference/modules.html#module-reference)).


**Progress: 42%**
`[█████████████——————————————]`

| Module             | Support        | Notes                |
|--------------------|----------------|-----------------------|
| Ansible            | TBD            |                       |
| APK Configure      | **Full**            |                       |
| Apt Configure      | TBD            |                       |
| Apt Pipelining      | **Full**            | Funny work-around is involved    |
| Bootcmd             | _Partial_            |          For now only "array of strings" supported, "array of array of strings" TBD     |
| Byobu              |  **Full**           |                       |
| CA Certificates     | **Full**            |                       |
| Chef               | TBD            |                       |
| Disable EC2 Instance Metadata Service| **Full**            |                       |
| Disk Setup          | TBD            |                       |
| Fan                | **Full**            |                       |
| Final Message      | **Full**            |                       |
| Growpart           | **Full**            |                       |
| GRUB dpkg           | **Full**            |                       |
| Install Hotplug     | **Full**            |                       |
| Keyboard           | **Full**            |                       |
| Keys to Console    | **Full**            |                       |
| Landscape          | TBD            |                       |
| Locale             | **Full**            |                       |
| LXD                | TBD            |                       |
| MCollective        | TBD            |                       |
| Mounts             | TBD            |                       |
| NTP                | TBD            |                       |
| Package Update Upgrade Install| _Partial_            |    Attributes `apt` and `snap` aren't supported                  |
| Phone Home          | TBD            |                       | 
| Power State Change  | TBD            |                       |
| Puppet             | TBD            |                       |
| Raspberry Pi Configuration| TBD            |                       |
| Resizefs           | **Full**            |                       |
| Resolv Conf        | TBD            |                       |
| Red Hat Subscription| TBD            |                       |
| Rsyslog            | TBD            |                       |
| Runcmd             | _Partial_            |     For now only "array of strings" supported, "array of array of strings" TBD                |
| Salt Minion        | _Partial_            |   `conf` and `grains` is yet to be supported    |
| Scripts Per Boot   | **Full**            |  This module actually doesn't have any configuration   |
| Scripts Per Instance| **Full**            |  This module actually doesn't have any configuration   |
| Scripts Per Once   | **Full**            |  This module actually doesn't have any configuration   |
| Scripts User       | **Full**            |  This module actually doesn't have any configuration   |
| Scripts Vendor     | TBD            |                       |
| Seed Random        | TBD            |                       |
| Set Hostname        | **Full**            |                       |
| Set Passwords      | _Partial_            |    Mostly supported, lacking validation                   |
| Snap               | TBD            |                       |
| Spacewalk          | TBD            |                       |
| SSH                | _Partial_            |         Only `ssh_authorized_keys` is supported              |
| SSH AuthKey Fingerprints | TBD            |                       |
| SSH Import ID      | TBD            |                       |
| Timezone           | **Full**            |           No brainer            |
| Ubuntu Drivers     | TBD            |                       |
| Ubuntu Autoinstall | **Full**            |                       |
| Ubuntu Pro         | TBD            |                       |
| Update Etc Hosts  | **Full**            |      Supported as two fields, because in original it's `true/false/'localhost'`                 |
| Update Hostname    | TBD            |                       |
| Users and Groups   | _Partial_            |     No support for deprecated fields, no support for nested `groups` object                  |
| Wireguard          | TBD            |                       |
| Write Files        | TBD            |                       |
| Yum Add Repo       | TBD            |                       |
| Zypper Add Repo    | TBD            |                       |

----

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/S6S1UZ9P7)
