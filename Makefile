VARFILES = -var-file=./config/variables.json -var-file=./config/output/domains.json -var-file=./config/output/a_records.json -var-file=./config/output/cname_records.json -var-file=./config/output/mx_records.json

init: flatten
	terraform init $(VARFILES)

apply: destroy flatten
	terraform apply $(VARFILES) -target=module.domains -auto-approve && \
	sleep 5 && \
	terraform apply $(VARFILES) -target=module.records -auto-approve

plan: flatten
	terraform plan $(VARFILES)

refresh: flatten
	terraform refresh $(VARFILES)

destroy: flatten
	terraform destroy $(VARFILES)

flatten:
	go run ./flatten.go
