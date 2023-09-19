package core

import (
	"flag"
	"os"
	"path/filepath"
	"puzzle/config"
	"puzzle/gologger"
	"puzzle/util"
	"strings"
)

type Options struct {
	Model      string   //使用模式
	Domain     []string //域名
	Level3     bool     //爆破三级域名
	Ip         []string //IP
	Port       string   //端口
	WebScan    bool     //web服务指纹扫描
	PortThread int      //端口爆破线程数
	WebThread  int      //指纹爆破线程数
	WebTimeout int      //指纹扫描超时时间
	Proxy      string
	Ping       bool   //存活探测
	SerialIp   bool   //ip连续段
	Vul        bool   //漏洞扫描
	Output     string // 输出名
}

func ParseOptions() *Options {
	options := &Options{}
	arg := flag.NewFlagSet(os.Args[0], 1)
	update := arg.Bool("update", false, "更新qqwrt库和指纹库")
	model := arg.String("m", "all", "选择模式,现有all,domain,ip,默认为all模式")
	domain := arg.String("d", "", "从单独域名中读取扫描")
	domainList := arg.String("dl", "", "从文件中读取扫描域名")
	output := arg.String("o", "", "输出文件名,如result")

	ip := arg.String("i", "", "ip输入,仅支持单次输入如192.168.1.1 or 192.168.1.1/24")
	ipList := arg.String("ipl", "", "从文件中获取ip")
	port := arg.String("p", "top1000", "端口号,如1-65535,22,3306,默认为top1000")

	arg.StringVar(&options.Proxy, "proxy", "", "web扫描代理,如socks5://127.0.0.1:8080")
	arg.IntVar(&options.PortThread, "pt", 500, "端口爆破线程,默认500")
	arg.IntVar(&options.WebThread, "wt", 25, "web指纹爆破线程,默认25")
	arg.IntVar(&options.WebTimeout, "timeout", 10, "web指纹扫描超时数,默认10")

	arg.BoolVar(&options.Ping, "ping", false, "是否开启ping探测,默认为false")
	arg.BoolVar(&options.Level3, "l3", false, "是否爆破三级域名，默认为false")
	arg.BoolVar(&options.SerialIp, "serial", false, "是否开启连续ip段检测")
	arg.BoolVar(&options.Vul, "vul", false, "是否开启漏洞扫描")

	arg.Parse(os.Args[1:])
	ShowBanner()

	//显示参数信息
	if (*model == "" || *output == "") && *update == false {
		arg.Usage()
		os.Exit(0)
	}

	if *update == true {
		util.Download(config.GepLiteUrl, filepath.Join(util.GetRunDir()+config.GeoLitePath))
		util.Download(config.FingerUrl, filepath.Join(util.GetRunDir()+config.FingerPrintPath))
		url := util.GetGithubLatestUrl(config.QQwryGithubRepo)
		if url != "" {
			util.Download(url, filepath.Join(util.GetRunDir()+config.QqwryPath))
		} else {
			gologger.Warningf("获取纯真数据库失败")
		}
		gologger.Infof("完成更新")
		os.Exit(0)
	}

	//判断模式
	if *model != "ip" && *model != "domain" && *model != "all" {
		gologger.Fatalf("模式%s不存在", *model)
	} else {
		options.Model = *model
	}

	if *model == "all" {
		options.WebScan = true
	}

	// 读取域名资源文件
	if *domain != "" {
		options.Domain = append(options.Domain, *domain)
	}
	if *domainList != "" {
		dl, err := util.LinesInFile(*domainList)
		if err != nil {
			gologger.Fatalf("读取地址文件失败:%s\n", err.Error())
		}
		options.Domain = append(dl, options.Domain...)
	}

	// 读取ip资源文件
	if *ip != "" {
		options.Ip = append(options.Ip, *ip)
	}
	if *ipList != "" {
		ipl, err := util.LinesInFile(*ipList)
		if err != nil {
			gologger.Fatalf("读取地址文件失败:%s", err.Error())
		}
		options.Ip = append(ipl, options.Ip...)
	}

	//域名去重
	if len(options.Domain) != 0 {
		options.Domain = util.RemoveRepeatedStringElement(options.Domain)
	}

	//port
	if *port != "" {
		options.Port = strings.TrimSpace(*port)
	}

	// 输出文件
	if *output != "" {
		options.Output = config.OutDir + "/" + *output + ".xlsx"
		if util.FileExists(options.Output) {
			gologger.Fatalf("该文件已存在")
		}
		err := XlsxInit(options.Output)
		if err != nil {
			gologger.Fatalf("初始化xlsx文件失败:%s", err.Error())
		}
	} else {
		gologger.Fatalf("请定义一个输出位置,参数为-o")
	}

	return options
}
