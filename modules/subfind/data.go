package subfind

import (
	"Allin/gologger"
	_ "Allin/statik"
	"github.com/rakyll/statik/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetAsnData() []byte {
	statikFS, err := fs.New()
	if err != nil {
		gologger.Fatalf(err.Error())
	}
	r, err := statikFS.Open("/GeoLite2-ASN.mmdb")
	if err != nil {
		gologger.Fatalf("打开资源文件失败:%s", err.Error())
	}
	defer r.Close()
	asnData, err := ioutil.ReadAll(r)
	if err != nil {
		gologger.Fatalf("读取资源文件失败:%s", err.Error())
	}
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
