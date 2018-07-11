variable "DIGITALOCEAN_API_TOKEN" {}

provider "digitalocean" {
  token = "${var.DIGITALOCEAN_API_TOKEN}"
}

variable "domains_flattened" {
    type = "map"
}

variable "a_records_flattened" {
    type = "map"
}

variable "mx_records_flattened" {
    type = "map"
}

variable "cname_records_flattened" {
    type = "map"
}

module "domains" {
  source  = "./domains"
  domains = "${var.domains_flattened}"
}

module "records" {
  source     = "./records"
  a_records = "${var.a_records_flattened}"
  mx_records = "${var.mx_records_flattened}"
  cname_records = "${var.cname_records_flattened}"
}
