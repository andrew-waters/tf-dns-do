variable "a_records" {
    type = "map"
}

resource "digitalocean_record" "a_records" {
    count    = "${length(var.a_records)}"
    domain   = "${lookup(var.a_records[count.index], "domain")}"
    type     = "A"
    name     = "${lookup(var.a_records[count.index], "name")}"
    value    = "${lookup(var.a_records[count.index], "value")}"
}
