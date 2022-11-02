package core

import (
	"net"
	"path/filepath"
	"puzzle/gologger"
	"puzzle/modules/ip"
	"puzzle/modules/ip/portscan"
	"puzzle/modules/subfind"
	"puzzle/util"
	"strings"
)

func allStart(options *Options) {
	var domains []string
	var ipsTmp []string
	var ips []string

	//拆分ip和域名
	for _, value := range options.Domain {
		if net.ParseIP(value) == nil {
			domains = append(domains, value)
		} else {
			ipsTmp = append(ipsTmp, value)
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
		ips := portscan.Ping(ips)
	}

	//ip扫描
	ipOptions := &ip.Options{
		Hosts:     ips,
		PortRange: options.Port,
		Threads:   options.PortThread,
		Rate:      3000,
		MaxPort:   200,
		ScanType:  getScanType(),
		QQwry:     getQqwry(),
	}

}
func getQqwry() *ip.QQwry {
	ip.IPData.FilePath = filepath.Join(util.GetRunDir() + "/config/qqwry.dat")
	res := ip.IPData.InitIPData()
	if v, ok := res.(error); ok {
		gologger.Fatalf(v.Error())
	}
	qqWry := ip.NewQQwry()
	return &qqWry
}
func getScanType() string {
	if util.IsOSX() || util.IsLinux() {
		return "s"
	} else {
		return "c"
	}
}
