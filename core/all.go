package core

import (
	"net"
	"puzzle/gologger"
	"puzzle/modules/ip/portscan"
	"puzzle/modules/ip/qqwry"
	"puzzle/modules/spider"
	"puzzle/modules/subfind"
	"puzzle/modules/webscan"
	"puzzle/util"
	"regexp"
	"strconv"
	"strings"
)

func AllStart(options *Options) {
	var domains []string
	var ipsTmp []string
	var ips []string

	//拆分ip和域名
	ipReg := `^((0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])\.){3}(0|[1-9]\d?|1\d\d|2[0-4]\d|25[0-5])`
	reg, _ := regexp.Compile(ipReg)
	for _, value := range options.Domain {
		if reg.MatchString(value) {
			ipsTmp = append(ipsTmp, value)
		} else {
			domains = append(domains, value)
		}
	}

	//解析ip格式
	for _, ipTmp := range ipsTmp {
		tmp, err := portscan.ParseIP(ipTmp)
		if err != nil {
			gologger.Errorf("ParseIP失败:%s", err.Error())
			return
		}
		ips = append(ips, tmp...)
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

	//写入域名爆破结果
	ReportWrite(options.Output, "子域名", domainResult)

	//提取ip
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

	//获取可能ip段
	if options.SerialIp == true {
		cloudIp, normalIp := qqwry.ExcludeCloud(ips)
		ips = util.GetSerialIp(normalIp)
		ips = append(ips, cloudIp...)
	}

	//ping检测
	if options.Ping == true {
		ips = portscan.Ping(ips)
	}

	//ip扫描
	ipOptions := &portscan.Options{
		Hosts:     ips,
		PortRange: options.Port,
		Threads:   options.PortThread,
	}
	ipRunner := portscan.NewRunner(ipOptions)

	//位置信息获取
	QQwry := qqwry.GetQqwry()
	var ipInfoRes []*qqwry.ResultQQwry
	for _, ip := range ips {
		info := QQwry.Find(ip)
		ipInfoRes = append(ipInfoRes, &info)
	}

	ReportWrite(options.Output, "IP地址", ipInfoRes)

	//端口扫描
	portscanResults := ipRunner.Run()
	if len(portscanResults) == 0 {
		gologger.Infof("端口扫描结果为空")
	} else {
		ReportWrite(options.Output, "端口服务", portscanResults)
	}

	var urls []string

	//爬虫获取遗漏信息
	var subdomains []string

	for _, domainRes := range domainResult {
		subdomains = append(subdomains, domainRes.Domain)
	}
	fofaEmail, fofaKey := spider.GetFofaKey()
	if fofaEmail != "" && fofaKey != "" {
		spiderOptions := &spider.Options{
			Domains:    domains,
			Subdomains: subdomains,
			Ips:        ips,
			Email:      fofaEmail,
			Key:        fofaKey,
		}

		spiderRunner := spider.NewRunner(spiderOptions)
		spiderResults := spiderRunner.Run()
		urls = append(urls, spiderResults.Urls...)

		if len(spiderResults.AddSubdomains) != 0 {
			ReportWrite(options.Output, "子域名", spiderResults.AddSubdomains)
		}

		if len(spiderResults.AddDomains) != 0 {
			ReportWrite(options.Output, "Spider", spiderResults.AddDomains)
		}

		gologger.Infof("结束爬虫信息收集")
	}

	//web扫描
	for _, result := range portscanResults {
		if result.ServiceName == "http" || result.ServiceName == "https" || result.ServiceName == "ssl" {
			url := result.Addr + ":" + strconv.Itoa(result.Port)
			urls = append(urls, url)
		}
	}
	for _, result := range domainResult {
		urls = append(urls, result.Domain)
	}

	urls = util.RemoveRepeatedStringElement(urls)

	webOptions := &webscan.Options{
		Url:     urls,
		Threads: options.WebThread,
		Timeout: options.WebTimeout,
		Proxy:   options.Proxy,
	}
	webRunner := webscan.NewRunner(webOptions)
	webResult := webRunner.Run()

	ReportWrite(options.Output, "WEB指纹", webResult)

}
