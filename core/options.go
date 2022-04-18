package core

import (
	"Allin/gologger"
	"Allin/util"
	"flag"
	"os"
	"strconv"
)

type Options struct {
	Rate            int64
	Test            bool
	Resolvers       []string
	TimeOut         int
	Retry           int
	Level           int
	Stdin           bool
	Domain          []string // 地址
	FileName        string   // 域名字典
	Output          string   // 输出文件
	DomainLevel     int      // 爆破层级
	SkipWildCard    bool     // 跳过泛解析
	SubNameFileName string   // 三级域名字典文件
	FilterWildCard  bool     // 过滤泛解析结果
	CheckOrigin     bool     // 会从返回包检查DNS是否为设定的，防止其他包的干扰

}

func ParseOptions() *Options {
	options := &Options{}
	bandwith := flag.String("b", "1M", "宽带的下行速度，可以5M,5K,5G")
	domain := flag.String("d", "", "扫描地址")
	domain_list := flag.String("dl", "", "从文件中读取扫描")
	resolvers := flag.String("s", "", "resolvers文件路径,默认使用内置DNS")

	flag.StringVar(&options.FileName, "f", "", "域名字典路径")
	flag.StringVar(&options.SubNameFileName, "sf", "", "三级域名爆破字典文件(默认内置)")
	flag.StringVar(&options.Output, "o", "", "输出文件路径")
	flag.BoolVar(&options.Test, "test", false, "测试本地最大发包数")
	flag.IntVar(&options.DomainLevel, "l", 1, "爆破域名层级,默认爆破一级域名")
	flag.BoolVar(&options.SkipWildCard, "skip-wild", false, "跳过泛解析的域名")
	flag.BoolVar(&options.CheckOrigin, "check-origin", false, "会从返回包检查DNS是否为设定的，防止其他包的干扰")

	flag.Parse()
	options.Stdin = hasStdin()

	ShowBanner()

	// 读取资源文件
	if *domain != "" {
		options.Domain = append(options.Domain, *domain)
	}
	if *domain_list != "" {
		dl, err := util.LinesInFile(*domain_list)
		if err != nil {
			gologger.Fatalf("读取地址文件失败:%s\n", err.Error())
		}
		options.Domain = append(dl, options.Domain...)
	}

	//subdomain
	// handle resolver
	if *resolvers != "" {
		rs, err := util.LinesInFile(*resolvers)
		if err != nil {
			gologger.Fatalf("读取resolvers文件失败:%s\n", err.Error())
		}
		if len(rs) == 0 {
			gologger.Fatalf("resolvers文件内容为空\n")
		}
		options.Resolvers = rs
	} else {
		defaultDns := []string{"223.5.5.5", "223.6.6.6", "180.76.76.76", "119.29.29.29", "182.254.116.116", "114.114.114.115"}
		options.Resolvers = defaultDns
	}

	var rate int64
	suffix := string([]rune(*bandwith)[len(*bandwith)-1])
	rate, _ = strconv.ParseInt(string([]rune(*bandwith)[0:len(*bandwith)-1]), 10, 64)
	switch suffix {
	case "G":
		fallthrough
	case "g":
		rate *= 1000000000
	case "M":
		fallthrough
	case "m":
		rate *= 1000000
	case "K":
		fallthrough
	case "k":
		rate *= 1000
	default:
		gologger.Fatalf("unknown bandwith suffix '%s' (supported suffixes are G,M and K)\n", suffix)
	}
	packSize := int64(100) // 一个DNS包大概有74byte
	rate = rate / packSize
	options.Rate = rate

	//显示参数信息
	if len(options.Domain) == 0 && !hasStdin() && options.FileName == "" && !options.Test {
		flag.Usage()
		os.Exit(0)
	}

	if options.FileName != "" && !util.FileExists(options.FileName) {
		gologger.Fatalf("文件:%s 不存在!\n", options.FileName)
	}

	if options.FilterWildCard && options.Output == "" {
		gologger.Fatalf("启用了 -filter-wild后，需要搭配一个输出文件 '-o'")
	}

	return options
}

func hasStdin() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	if fi.Mode()&os.ModeNamedPipe == 0 {
		return false
	}
	return true
}
