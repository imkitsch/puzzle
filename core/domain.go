package core

import (
	"net"
	"puzzle/gologger"
	"puzzle/modules/ip/qqwry"
	"puzzle/modules/subfind"
	"puzzle/util"
	"strings"
)

func DomainStart(options *Options) {
	//域名爆破
	domainOptions := &subfind.Options{
		Domains:       options.Domain,
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

	//写入域名爆破结果
	ReportWrite(options.Output, "子域名", domainResult)

	//提取ip
	var ips []string
	for _, domainRes := range domainResult {
		if domainRes.Cdn == false && domainRes.Address != "" {
			addr := strings.Split(domainRes.Address, ",")
			for _, v := range addr {
				tmpIp := net.ParseIP(v)
				if tmpIp != nil && util.HasLocalIP(tmpIp) == false {
					ips = append(ips, v)
				}
			}
		}
	}

	//ip去重
	ips = util.RemoveRepeatedStringElement(ips)

	QQwry := qqwry.GetQqwry()
	var ipInfoRes []*qqwry.ResultQQwry
	for _, ip := range ips {
		info := QQwry.Find(ip)
		ipInfoRes = append(ipInfoRes, &info)
	}

	ReportWrite(options.Output, "IP地址", ipInfoRes)
}
