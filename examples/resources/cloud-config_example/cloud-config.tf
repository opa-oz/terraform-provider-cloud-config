resource "cloud-config" "main" {
  hostname                  = "domain"
  fqdn                      = "domain.lan"
  prefer_fqdn_over_hostname = true
  preserve_hostname         = true
  create_hostname_file      = false

  locale            = "en_GB"
  locale_configfile = "/etc/locale"

  timezone = "Asia/Tokyo"

  runcmd = ["echo '11'", "cat /etc/hosts"]

  manage_etc_hosts_localhost = true
}
