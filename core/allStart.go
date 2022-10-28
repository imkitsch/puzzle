package core

import (
	"net"
	"puzzle/gologger"
	"puzzle/modules/subfind"
)

func allStart(options *Options) {
	var domains = []string{}
	var ips = []string{}
	for _, value := range options.Domain {
		if net.ParseIP(value) == nil {
			domains = append(domains, value)
		} else {
			ips = append(ips, value)
		}
	}

	//域名爆破
	domainOptions := &subfind.Options{
		Domains:       domains,
		SubdomainDict: subfind.GetSubdomainData(),
		SubNextDict:   subfind.GetSubNextData(),
		Level3:        options.Level3,
		DeviceConfig:  subfind.GetDeviceConfig(),
	}
	domainRunner, err := subfind.NewRunner(domainOptions)
	if err != nil {
		gologger.Fatalf("创建DNS爆破失败:%s", err.Error())
	}
	domainResult := domainRunner.Run()
	_ = domainResult

	//写入域名爆破结果

}
