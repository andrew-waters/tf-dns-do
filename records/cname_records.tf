variable "cname_records" {
    type = "map"
}

resource "digitalocean_record" "cname_records" {
    count    = "${length(var.cname_records)}"
    domain   = "${lookup(var.cname_records[count.index], "domain")}"
    type     = "CNAME"
    name     = "${lookup(var.cname_records[count.index], "name")}"
    value    = "${lookup(var.cname_records[count.index], "value")}"
}
