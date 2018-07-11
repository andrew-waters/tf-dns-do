variable "mx_records" {
    type = "map"
}

resource "digitalocean_record" "mx_records" {
    count    = "${length(var.mx_records)}"
    domain   = "${lookup(var.mx_records[count.index], "domain")}"
    type     = "MX"
    priority = "${lookup(var.mx_records[count.index], "priority")}"
    name     = "${lookup(var.mx_records[count.index], "name")}"
    value    = "${lookup(var.mx_records[count.index], "value")}"
}
