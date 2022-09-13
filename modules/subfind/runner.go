package subfind

import (
	"Allin/gologger"
	"Allin/util"
)

func Run(domains []string) []domainResult {
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
			for _, sub := range subdomainDict {
				subDomains = append(subDomains, sub+"."+domain)
				for _, subNext := range subNextDict {
					subDomains = append(subDomains, subNext+"."+sub+"."+domain)
				}
			}
		} else {
			gologger.Infof("域名%s存在泛解析,自动给跳过爆破", domain)
		}

		//去重
		subDomains = util.RemoveRepeatedStringElement(subDomains)

		//dns爆破验证
		dr = append(dr, DomainBlast(subDomains)...)
	}

	return dr
}
