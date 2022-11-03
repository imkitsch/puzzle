package main

import "puzzle/core"

func main() {
	options := core.ParseOptions()
	switch options.Model {
	case "all":
		core.AllStart(options)
	case "domain":
		core.DomainStart(options)
	case "ip":
		core.IpStart(options)

	}

}
