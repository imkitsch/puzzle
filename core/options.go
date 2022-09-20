package core

import (
	"Allin/gologger"
	"Allin/util"
	"flag"
	"os"
)

type Options struct {
	Model      string   //使用模式
	Domain     []string //域名
	Level3     bool     //爆破三级域名
	Ip         []string //IP
	Port       []int    //端口号
	WebScan    bool     //web服务指纹扫描
	PortThread int      //端口爆破线程数
	WebThread  int      //指纹爆破线程数
	Ping       bool     //存活探测
	Output     string   // 输出名
}

func ParseOptions() *Options {
	options := &Options{}
	model := flag.String("m", "all", "选择模式,现有all,domain,ip,默认为all模式")
	domain := flag.String("d", "", "从单独域名中读取扫描")
	domainList := flag.String("dl", "", "从文件中读取扫描域名")
	output := flag.String("o", "", "输出文件名,如result")

	ip := flag.String("i", "", "ip输入,仅支持单次输入如192.168.1.1 or 192.168.1.1/24")
	ipList := flag.String("ipl", "", "从文件中获取ip")
	port := flag.String("p", "", "端口号,如1-65535,21,22,3306,默认为top1000")

	flag.IntVar(&options.PortThread, "pt", 800, "端口爆破线程,默认800")
	flag.IntVar(&options.WebThread, "wt", 25, "web指纹爆破线程,默认25")
	flag.BoolVar(&options.Ping, "ping", false, "是否开启ping探测,默认为false")
	flag.BoolVar(&options.Level3, "l3", false, "是否爆破三级域名，默认为false")
	flag.BoolVar(&options.WebScan, "ws", false, "是否开启web指纹扫描,all模式默认开启")
	flag.Parse()
	ShowBanner()

	//判断模式
	if *model != "ip" && *model != "domain" && *model != "all" {
		gologger.Fatalf("模式%s不存在", *model)
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
			gologger.Fatalf("读取地址文件失败:%s\n", err.Error())
		}
		options.Ip = append(ipl, options.Ip...)
	}

	//域名去重
	options.Domain = util.RemoveRepeatedStringElement(options.Domain)

	//port
	if *port != "" {
		options.Port = util.RemoveRepeatedIntElement(util.ParsePorts(*port))
	} else {
		options.Port = util.RemoveRepeatedIntElement(util.ParsePorts(PortTop1000))
	}

	// 输出文件
	if *output != "" {
		options.Output = OutDir + "/" + *output + ".xlsx"
		if util.FileExists(options.Output) {
			gologger.Fatalf("该文件已存在")
		}
		if XlsxInit(options.Output) != nil {
			gologger.Fatalf("初始化xlsx文件失败")
		}
	} else {
		gologger.Fatalf("请定义一个输出位置,参数为-o")
	}
	//显示参数信息
	if len(options.Domain) == 0 || options.Output == "" {
		flag.Usage()
		os.Exit(0)
	}

	return options
}
