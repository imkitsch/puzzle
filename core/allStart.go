package core

import (
	"net"
	"path/filepath"
	"puzzle/gologger"
	"puzzle/modules/ip"
	"puzzle/modules/ip/portfinger"
	"puzzle/modules/ip/portscan"
	"puzzle/modules/subfind"
	"puzzle/util"
	"strings"
)

func AllStart(options *Options) {
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
		ips = portscan.Ping(ips)
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
		NmapProbe: getNmapProbe(),
	}
	ipRunner, err := ip.NewRunner(ipOptions)
	if err != nil {
		gologger.Fatalf(err.Error())
	}

	//位置信息获取
	var ipInfoRes []*ip.ResultQQwry
	for _, ip := range ips {
		info := ipOptions.QQwry.Find(ip)
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

func getQqwry() *ip.QQwry {
	ip.IPData.FilePath = filepath.Join(util.GetRunDir() + QqwryPath)
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

func getNmapProbe() *portfinger.NmapProbe {
	NmapProbe := portfinger.NmapProbe{}
	nmapData, err := util.ReadFile(util.GetRunDir() + NmapPath)
	if err != nil {
		gologger.Fatalf("读取nmap指纹文件失败:%s", err.Error())
	}
	if err = NmapProbe.Init(nmapData); err != nil {
		gologger.Fatalf("nmap指纹初始化失败:%s", err.Error())
	}
	gologger.Infof("nmap指纹数量: %v个探针,%v条正则", len(NmapProbe.Probes), NmapProbe.Count())
	return &NmapProbe
}
