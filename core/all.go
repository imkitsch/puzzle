package core

import (
	"path/filepath"
	"puzzle/gologger"
	"puzzle/modules/ip/portscan"
	"puzzle/modules/ip/qqwry"
	"puzzle/modules/subfind"
	"puzzle/util"
	"regexp"
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
		if domainRes.Cdn == false {
			addr := strings.Split(domainRes.Address, ",")
			ips = append(ips, addr...)
		}
	}

	//ip去重
	ips = util.RemoveRepeatedStringElement(ips)

	//获取可能ip段
	ips = util.GetSerialIp(ips)

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
	QQwry := getQqwry()
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
}

func getQqwry() *qqwry.QQwry {
	qqwry.IPData.FilePath = filepath.Join(util.GetRunDir() + QqwryPath)
	res := qqwry.IPData.InitIPData()
	if v, ok := res.(error); ok {
		gologger.Fatalf(v.Error())
	}
	qqWry := qqwry.NewQQwry()
	return &qqWry
}
