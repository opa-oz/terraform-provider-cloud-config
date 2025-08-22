---
layout: page
page_title: "Example configuration"
subcategory: Guides
description: |-
    Example configuration from real cloud-init
---

# Example configuration

## Original cloud-config

[Original source](https://www.techtutorials.tv/sections/promox/automate-vm-creation-on-proxmox-with-ansible/)

```yaml
#cloud-config
hostname: ubuntu
manage_etc_hosts: true
fqdn: ubuntu.lan
user: ansible
ssh_authorized_keys:
  - let's pretend it's ssh public key  
chpasswd:
  expire: False
users:
  - default
package_update: true
package_upgrade: true
timezone: Asia/Tokyo
packages:
  - qemu-guest-agent
  - ufw
runcmd:
    - systemctl enable qemu-guest-agent
    - systemctl start qemu-guest-agent
    - ufw limit proto tcp to any port 22
    - ufw enable
    - locale-gen en_GB.UTF-8
    - localectl set-locale LANG=en_GB.UTF-8
    - chfn -f Ansible ansible
```

## As terraform configuration

_Kudos to [bpg/terraform-provider-proxmox](https://github.com/bpg/terraform-provider-proxmox/tree/main)_

```terraform
data "local_file" "ssh_public_key" {
  filename = var.ssh_key
}
locals {
  timezone     = "Asia/Tokyo"
  locale       = "en_GB.UTF-8"
  ansible_user = "ansible"
}

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

resource "proxmox_virtual_environment_file" "user_data_cloud_config" {
  provider     = proxmox.api
  content_type = "snippets"
  datastore_id = "local"
  node_name    = var.node_name

  source_raw {
    data      = cloud-config.vm_cloud_config.content
    file_name = "user-data-${var.vmid}-cloud-config.yaml"
  }
}
```

After `terraform apply`, if we check on a host, we will get:

```yaml
#cloud-config
hostname: test-ubuntu
fqdn: test-ubuntu.lan
timezone: Asia/Tokyo
runcmd:
    - systemctl enable qemu-guest-agent
    - systemctl start qemu-guest-agent
    - ufw limit proto tcp to any port 22
    - ufw enable
    - locale-gen en_GB.UTF-8
    - localectl set-locale LANG=en_GB.UTF-8
    - chfn -f Ansible ansible
manage_etc_hosts: true
ssh_authorized_keys:
    - ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCAH3ADjeJ77EEZkvHAfMalLtrX8oSGp0lX1+Nu06j4VlUz89wIoeHBSOUup41E78TRfF5+lKN0HlmJLH6NMpa79+JgSmeBf3WoQ1D2yyJ/lvotKi6Ni7EtDp/stXKoQ1u6zYVJs1iRFbDeTkZT8SeAx9MHqUH2o01D+hR1GfcJoyZhrK5yCwcoiaMnfFZz9atjNDg9D5lONiFNYD30zXYF9UZxIQGhowzcUfVFvoqzXjEN4Vqo3q0BgU2chgNRHFpI/3UcXM5xhEI4n0HgXCYLD0xuTVryHbpaihL0AXRmnfCkJeQYQkPdYjIsKWjiqceIIwgb8lv2atfJDzlljiwZ
chpasswd:
    expire: false
package_update: true
package_upgrade: true
packages:
    - qemu-guest-agent
    - ufw
user:
    name: ansible
users:
    - name: default
```

