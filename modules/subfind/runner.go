package subfind

import (
	"puzzle/gologger"
	"puzzle/util"
)

type Options struct {
	Domains       []string
	Level3        bool
	SubdomainDict []string
	SubNextDict   []string
	DeviceConfig  *EtherTable
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
		wildcard := IsWildCard(domain)
		if wildcard == false {
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
		domainResult := DomainBlast(subDomains, r.options.DeviceConfig)
		dr = append(dr, domainResult...)
		gologger.Printf("\n")

		//三级子域名爆破:仅对本轮域名的一二级结果展开,避免遍历累积结果导致重复爆破;
		//泛解析域名下所有子域都会命中,同样跳过,避免假阳性
		if r.options.Level3 && wildcard == false {
			gologger.Infof("域名 %s 三级子域名爆破", domain)
			for _, sub := range domainResult {
				//清空
				subDomains = []string{}
				for _, subNext := range r.options.SubNextDict {
					subDomains = append(subDomains, subNext+"."+sub.Domain)
				}
				dr = append(dr, DomainBlast(subDomains, r.options.DeviceConfig)...)
			}
		}

	}

	//按域名去重,避免重复结果进入报表
	seen := make(map[string]struct{}, len(dr))
	unique := dr[:0]
	for _, res := range dr {
		if _, ok := seen[res.Domain]; ok {
			continue
		}
		seen[res.Domain] = struct{}{}
		unique = append(unique, res)
	}

	return unique
}
