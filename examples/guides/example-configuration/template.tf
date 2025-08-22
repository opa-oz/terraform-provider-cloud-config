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


