# Digital Ocean DNS Management

This project allows you to manage Domains and DNS Records on DigitalOcean using Terraform and config files. It takes a yaml input file (`./config/domains.yaml`), transforms the spec to something that terraform understands and then creates the resources you need.

DNS Management is [free in DigitalOcean](https://www.digitalocean.com/docs/networking/dns/overview/) which makes it an ideal candidate to dip your toes into IaC (Infrastructure as Code) if you are unfamiliar with it.

There are two types of resources which are created: [domains](https://www.terraform.io/docs/providers/do/r/domain.html) and [records](https://www.terraform.io/docs/providers/do/r/record.html).

It becomes useful to script these resources when you have a significant number of domains to manage regularly so if you're just managing a few, this may be overkill but will still work.

## IMPORTANT: Applying Plans

Instead of simply running `terraform apply` which would result in a race condition between the domain creation and the record creation (and fail) **you need to run** `make apply` (warning: auto approves after destroy).

This one is a major PITA. The upstream DO provider for terraform wraps around the DOGO SDK and doesn't [implement proper validation on updates](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/58).

So the less than ideal workaround to this is that you cannot apply changes in the way you're used to with TF. Instead, the infra must be destroyed and then applied from scratch. EVEN IF you're not updating any records, this still applies. So a `make plan` will give you the changes to be implemented but a `make apply` will destroy the resources first. When that issue is resolved, I'll remove this notice and update this in the background.

## Running

Before each `make` operation, we parse the yaml config and generate JSON which terraform will understand. Notoriously troubled varaibles in TF mean the generated output is more verbose than anyone would like. The good news is, you don't need to see it - you only work with much cleaner yaml. When TF gets better, this project will follow and we hopefully won't need `flatten.go`.

So, first copy the config template and the variable file:

```
cp  ./config/domains.tpl.yaml ./config/domains.yaml && \
    ./config/variables.tpl.yaml ./config/variables.yaml 
```

Using your favourite editor, dive into the config files and add your domains (see below for domain details).

You should get an API token from DO and add it to the `DIGITALOCEAN_API_TOKEN` field in `./config/variables.yaml`.

Once you have completed your config, you can run:

```
make plan
```

Inspect the output and check you're happy to proceed. When you are:

```
make apply
```

Your records should now be created in DO. Go profit!

#### Domain Configuration

There is a top level `domains` attribute which is a list of the domains you want to manage. Each domain has several attributes:

Name   | Type       | Description
------ | ---------- | -----------
`name` | `string`   | The domain name (eg `github.com`)
`ip`   | `ipv4`     | The primary IP address (eg `192.30.253.112`)
`a`    | `[]record` | an array of A records
`mx`   | `[]record` | an array of MX records
`c`    | `[]record` | an array of CNAME records

#### Record Configuration

The record object referenced above depends on the type of record you're creating. It currently supports `A`, `MX` and `CNAME` records and the following properties:

Name       | Type     | Required On        | Description
---------- | -------- | ------------------ | -----------
`name`     | `string` | `A`, `CNAME`, `MX` | The zone name
`value`    | `string` | `A`, `CNAME`, `MX` | The record data
`priority` | `int`    | `MX`               | Priority for mail servers ([1](https://github.com/terraform-providers/terraform-provider-digitalocean/issues/6)-100)
