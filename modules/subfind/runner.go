package subfind

import (
	"puzzle/gologger"
	"puzzle/util"
)

func Run(domains []string, level3 bool) []domainResult {
	var dr []domainResult = []domainResult{}
	var subDomains = []string{}

	var subdomainDict = GetSubdomainData()
	var subNextDict = GetSubNextData()

	for _, domain := range domains {
		subDomains = append(subDomains, domain)
		//api获取
		subDomains = append(subDomains, DoSubFinder(domain)...)
		//判断泛解析
		if IsWildCard(domain) == false {
			//加载字典
			gologger.Infof("域名 %s 装载爆破字典", domain)
			for _, sub := range subdomainDict {
				subDomains = append(subDomains, sub+"."+domain)
			}
		} else {
			gologger.Infof("域名 %s 存在泛解析,自动跳过爆破", domain)
		}

		//去重
		subDomains = util.RemoveRepeatedStringElement(subDomains)

		//dns爆破验证
		gologger.Infof("域名 %s 开始验证DNS", domain)
		dr = append(dr, DomainBlast(subDomains)...)

		//三级子域名爆破
		if level3 {
			gologger.Infof("域名 %s 三级子域名爆破", domain)
			for _, sub := range dr {
				//清空
				subDomains = nil
				for _, subNext := range subNextDict {
					subDomains = append(subDomains, subNext+"."+sub.domain)
				}
				dr = append(dr, DomainBlast(subDomains)...)
			}
		}
	}

	return dr
}
