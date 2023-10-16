package spider

import (
	"bufio"
	_ "embed"
	"puzzle/util"
	"strings"
)

//go:embed data/tld-data.txt
var rootDomainList string

func getSubdomainData() []string {
	reader := bufio.NewScanner(strings.NewReader(rootDomainList))
	reader.Split(bufio.ScanLines)
	var ret []string
	for reader.Scan() {
		ret = append(ret, reader.Text())
	}
	return ret
}

func domain_parse(domain string) string {
	list := util.ConvertStrSlice2Map(getSubdomainData())
	step := strings.Split(domain, ".")
	tmp := step[len(step)-2] + "." + step[len(step)-1]
	if util.InMap(list, tmp) == true && len(step) >= 3 {
		return step[len(step)-3] + "." + tmp
	} else {
		return tmp
	}
}
