package core

import (
	"puzzle/gologger"
	"puzzle/modules/ip"
	"puzzle/modules/ip/portscan"
	"puzzle/util"
)

func IpStart(options *Options) {
	var ips []string

	//解析ip格式
	for _, ipTmp := range options.Ip {
		tmp, err := portscan.ParseIP(ipTmp)
		if err != nil {
			gologger.Errorf("ParseIP失败:%s", err.Error())
			return
		}
		ips = append(ips, tmp...)
	}

	//ip去重
	ips = util.RemoveRepeatedStringElement(ips)

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
	gologger.Infof("获取ip位置信息")
	for _, ip := range ips {
		info := ipOptions.QQwry.Find(ip)
		gologger.Infof("%v", info)
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
