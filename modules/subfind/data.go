package subfind

import (
	"Allin/gologger"
	_ "embed"
	"github.com/rakyll/statik/fs"
	"os"
	"path/filepath"
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

func GetDefaultSubNextData() []string {
	return strings.Split(subnext, "\n")
}

func GetAsnData() []byte {
	return asnData
}

func getDefaultScripts() []string {
	var scripts []string
	StatikFS, err := fs.New()
	if err != nil {
		gologger.Fatalf(err.Error())
	}
	fs.Walk(StatikFS, "/scripts", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Is this file not a script?
		if info.IsDir() || filepath.Ext(info.Name()) != ".lua" {
			return nil
		}
		// Get the script content
		data, err := fs.ReadFile(StatikFS, path)
		if err != nil {
			return err
		}
		scripts = append(scripts, string(data))
		return nil
	})

	return scripts
}
