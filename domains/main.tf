variable "domains" {
    type = "map"
}

resource "digitalocean_domain" "domains" {
    count      = "${length(var.domains)}"
    name       = "${lookup(var.domains[count.index], "name")}"
    ip_address = "${lookup(var.domains[count.index], "ip")}"
}
