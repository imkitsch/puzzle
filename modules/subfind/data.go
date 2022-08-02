package subfind

import (
	_ "embed"
	"strings"
)

//go:embed data/subnext.txt
var subnext string

//go:embed data/subdomain.txt
var subdomain string

//go:embed data/GeoLite2-ASN.mmdb
var asnData []byte

func GetSubdomainData() []string {
	return strings.Split(subdomain, "\n")
}

func GetSubNextData() []string {
	return strings.Split(subnext, "\n")
}

func GetAsnData() []byte {
	return asnData
}
