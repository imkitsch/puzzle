package subfind

import (
	"github.com/boy-hack/ksubdomain/core/device"
	"puzzle/gologger"
	"puzzle/util"
)

type Options struct {
	Domains       []string
	Level3        bool
	SubdomainDict []string
	SubNextDict   []string
	DeviceConfig  *device.EtherTable
}

type Runner struct {
	options *Options
}

func NewRunner(options *Options) (*Runner, error) {
	return &Runner{
		options: options,
	}, nil
}

func (r *Runner) Run() (dr []*domainResult) {
	for _, domain := range r.options.Domains {
		subDomains := []string{}
		subDomains = append(subDomains, domain)
		//api获取
		subDomains = append(subDomains, DoSubFinder(domain)...)
		//判断泛解析
		if IsWildCard(domain) == false {
			//加载字典
			gologger.Infof("域名 %s 装载爆破字典", domain)
			for _, sub := range r.options.SubdomainDict {
				subDomains = append(subDomains, sub+"."+domain)
			}
		} else {
			gologger.Infof("域名 %s 存在泛解析,自动跳过爆破", domain)
		}

		//去重
		subDomains = util.RemoveRepeatedStringElement(subDomains)

		//dns爆破验证
		gologger.Infof("域名 %s 开始验证DNS", domain)
		dr = append(dr, DomainBlast(subDomains, r.options.DeviceConfig)...)
		gologger.Printf("\n")

		//三级子域名爆破
		if r.options.Level3 {
			gologger.Infof("域名 %s 三级子域名爆破", domain)
			for _, sub := range dr {
				//清空
				subDomains = []string{}
				for _, subNext := range r.options.SubNextDict {
					subDomains = append(subDomains, subNext+"."+sub.Domain)
				}
				dr = append(dr, DomainBlast(subDomains, r.options.DeviceConfig)...)
			}
			gologger.Printf("\n")
		}

	}

	return
}
