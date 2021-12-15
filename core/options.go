package core

import "flag"

type Options struct {
	Rate            int64
	Test            bool
	Domain          []string // 单域名
	FileName        string   // 域名字典
	Output          string   // 输出文件
	DomainLevel     int      // 爆破层级
	SkipWildCard    bool     // 跳过泛解析
	API             bool     // 显示api
	SubNameFileName string   // 三级域名字典文件
	FilterWildCard  bool     // 过滤泛解析结果
	CheckOrigin     bool     // 会从返回包检查DNS是否为设定的，防止其他包的干扰
	Model           string   // 调用模式
	Poc             string   // 指定poc

}

func ParseOptions() *Options {
	options := &Options{}
	//bandwith := flag.String("b", "1M", "宽带的下行速度，可以5M,5K,5G")
	//domain := flag.String("d", "", "扫描地址")
	//domain_list := flag.String("dl", "", "从文件中读取扫描")
	//help := flag.Bool("h",false,"查看使用帮助")
	flag.StringVar(&options.Model, "m", "full", "使用模式,目前支持:full模式\\subfind模式\\pocscan模式")

	flag.StringVar(&options.FileName, "f", "", "域名字典路径")
	flag.StringVar(&options.SubNameFileName, "sf", "", "三级域名爆破字典文件(默认内置)")
	flag.StringVar(&options.Output, "o", "", "输出文件路径")
	flag.BoolVar(&options.Test, "test", false, "测试本地最大发包数")
	flag.BoolVar(&options.API, "api", false, "查看调用api")
	flag.IntVar(&options.DomainLevel, "l", 1, "爆破域名层级,默认爆破一级域名")
	flag.BoolVar(&options.SkipWildCard, "skip-wild", false, "跳过泛解析的域名")
	flag.BoolVar(&options.CheckOrigin, "check-origin", false, "会从返回包检查DNS是否为设定的，防止其他包的干扰")
	flag.StringVar(&options.Poc, "poc", "", "指定poc扫描,默认为扫描全部")

	return options
}
