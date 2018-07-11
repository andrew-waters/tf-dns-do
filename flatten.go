package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	yaml "gopkg.in/yaml.v1"
)

// Input is the main holder for the input data to parse
type Input struct {
	Domains []DomainInput
}

// DomainInput contains zone information for a single domain
type DomainInput struct {
	Name string
	IP   string
	A    []ARecordInput
	C    []CRecordInput
	MX   []MXRecordInput
}

// ARecordInput is information about an A record
type ARecordInput struct {
	Name  string
	Value string
}

// CRecordInput is information about a CNAME record
type CRecordInput struct {
	Name  string
	Value string
}

// MXRecordInput is information about an MX record
type MXRecordInput struct {
	Name     string
	Value    string
	Priority int
}

// BaseDomainOutput is the output for domains
type BaseDomainOutput struct {
	Output []DomainOutput `json:"domains_flattened"`
}

// BaseARecordOutput is the output for A records
type BaseARecordOutput struct {
	Output []ARecordOutput `json:"a_records_flattened"`
}

// BaseMXRecordOutput is the output for MX records
type BaseMXRecordOutput struct {
	Output []MXRecordOutput `json:"mx_records_flattened"`
}

// BaseCNAMERecordOutput is the output for CNAME records
type BaseCNAMERecordOutput struct {
	Output []CNAMERecordOutput `json:"cname_records_flattened"`
}

// DomainOutput is the target format for a domain
type DomainOutput struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
}

// ARecordOutput is the target format for an A Record
type ARecordOutput struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

// MXRecordOutput is the target format for an MX Record
type MXRecordOutput struct {
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Value    string `json:"value"`
	Priority int    `json:"priority"`
}

// CNAMERecordOutput is the target format for a CNAME Record
type CNAMERecordOutput struct {
	Domain string `json:"domain"`
	Name   string `json:"name"`
	Value  string `json:"value"`
}

func main() {

	input := Input{}

	data, err := ioutil.ReadFile("./config/domains.yaml")
	check(err)

	err = yaml.Unmarshal([]byte(data), &input)
	check(err)

	domainOutput, aRecordOutput, mxRecordOutput, cnameRecordOutput := generateDomainOutput(input.Domains)

	writeData(bastardise(domainOutput, "domains_flattened"), "./config/output/domains.json")
	writeData(bastardise(aRecordOutput, "a_records_flattened"), "./config/output/a_records.json")
	writeData(bastardise(mxRecordOutput, "mx_records_flattened"), "./config/output/mx_records.json")
	writeData(bastardise(cnameRecordOutput, "cname_records_flattened"), "./config/output/cname_records.json")
}

func generateDomainOutput(domains []DomainInput) ([]DomainOutput, []ARecordOutput, []MXRecordOutput, []CNAMERecordOutput) {
	domainOutput := make([]DomainOutput, 0)
	aRecordOutput := make([]ARecordOutput, 0)
	mxRecordOutput := make([]MXRecordOutput, 0)
	cnameRecordOutput := make([]CNAMERecordOutput, 0)
	for _, domain := range domains {
		domainOutput = append(domainOutput, DomainOutput{
			domain.Name,
			domain.IP,
		})
		for _, a := range domain.A {
			aRecordOutput = append(aRecordOutput, ARecordOutput{
				domain.Name,
				a.Name,
				a.Value,
			})
		}
		for _, mx := range domain.MX {
			mxRecordOutput = append(mxRecordOutput, MXRecordOutput{
				domain.Name,
				mx.Name,
				mx.Value,
				mx.Priority,
			})
		}
		for _, cname := range domain.C {
			cnameRecordOutput = append(cnameRecordOutput, CNAMERecordOutput{
				domain.Name,
				cname.Name,
				cname.Value,
			})
		}
	}
	return domainOutput, aRecordOutput, mxRecordOutput, cnameRecordOutput
}

func bastardise(input interface{}, key string) []byte {

	marshalledInput, _ := json.Marshal(input)

	var tmpOutput interface{}
	err := json.Unmarshal(marshalledInput, &tmpOutput)
	check(err)
	tmpOutputUnmarshalled := tmpOutput.([]interface{})

	newOutput := make(map[string]interface{})
	count := 0
	for _, tmpOutput := range tmpOutputUnmarshalled {
		newOutput[strconv.Itoa(count)] = []interface{}{tmpOutput}
		count++
	}

	newOutputInterface := make(map[string]interface{})
	newOutputInterface[key] = newOutput
	marshalled, err := json.MarshalIndent(newOutputInterface, "", "    ")
	check(err)

	return marshalled
}

func writeData(data []byte, location string) {
	file, err := os.Create(location)
	check(err)
	defer file.Close()

	w := bufio.NewWriter(file)
	_, err = w.WriteString(string(data))
	check(err)

	w.Flush()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
