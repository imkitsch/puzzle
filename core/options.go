package core

import (
	"Allin/gologger"
	"Allin/util"
	"flag"
	"os"
)

type Options struct {
	Domain     []string //域名
	DomainFile string   //域名输入文件
	PortThread int      //端口爆破线程数
	WebThread  int      //指纹爆破线程数
	NoPing     bool     //存活探测
	Output     string   // 输出名

}

func ParseOptions() *Options {
	options := &Options{}
	domain := flag.String("d", "", "从文件中读取扫描")
	domainList := flag.String("dl", "", "从文件中读取扫描")
	output := flag.String("o", "", "输出文件名,如result")

	flag.IntVar(&options.PortThread, "pt", 800, "端口爆破线程,默认800")
	flag.IntVar(&options.WebThread, "wt", 25, "web指纹爆破线程,默认25")
	flag.BoolVar(&options.NoPing, "np", false, "是否开启存活探测,默认为false")
	flag.Parse()
	ShowBanner()

	// 读取资源文件
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

	// 输出文件
	if *output != "" {
		options.Output = OutDir + "/" + *output + ".xlsx"
		if util.FileExists(options.Output) {
			gologger.Fatalf("该文件已存在")
		}
	} else {
		gologger.Fatalf("请定义一个输出位置,参数为-o")
	}
	//显示参数信息
	if len(options.Domain) == 0 && options.Output == "" && options.DomainFile == "" {
		flag.Usage()
		os.Exit(0)
	}

	return options
}
