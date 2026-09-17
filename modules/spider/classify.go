package spider

import (
	"net/url"
	"puzzle/gologger"
	"puzzle/modules/subfind"
	"puzzle/util"
	"strings"
)

type classifyMode int

const (
	// certMode 证书检索结果分类:IP记录做CDN过滤后作为扫描目标产出
	certMode classifyMode = iota
	// ipMode IP检索结果分类:输入已是扫描目标IP,IP记录直接进Urls
	ipMode
)

// classifySets 预构建的目标集合,避免分类时对切片反复线性查找
type classifySets struct {
	domains    map[string]struct{}
	subdomains map[string]struct{}
}

func newClassifySets(domains []string, subdomains []string) *classifySets {
	sets := &classifySets{
		domains:    make(map[string]struct{}, len(domains)),
		subdomains: make(map[string]struct{}, len(subdomains)),
	}
	for _, d := range domains {
		sets.domains[d] = struct{}{}
	}
	for _, s := range subdomains {
		sets.subdomains[s] = struct{}{}
	}
	return sets
}

// classify 对测绘检索结果统一分类,组装Result。
// 分类为纯CPU工作(CDN判定带per-IP缓存),串行执行即可,无需并发
func classify(results []*searchResult, sets *classifySets, mode classifyMode) *Result {
	var result Result

	addDomains := make(map[string]struct{})
	// host与ip成对去重,等价原[][]string整行去重
	type subdomainPair struct{ host, ip string }
	subdomainSeen := make(map[string]struct{})
	var subdomainOrder []subdomainPair

	for _, value := range results {
		if value.protocol != "http" && value.protocol != "https" && value.protocol != "tls" && value.protocol != "unknown" {
			continue
		}
		if !strings.HasPrefix(value.host, "http") {
			value.host = "http://" + value.host
		}

		tmpUrl, err := url.Parse(value.host)
		if err != nil {
			continue
		}
		host := tmpUrl.Hostname()

		if value.domain != "" {
			if _, ok := sets.domains[value.domain]; ok {
				if _, ok := sets.subdomains[host]; !ok {
					key := host + "\x00" + value.ip
					if _, dup := subdomainSeen[key]; !dup {
						gologger.Infof("查询到遗漏子域名: %s", host)
						subdomainSeen[key] = struct{}{}
						subdomainOrder = append(subdomainOrder, subdomainPair{host: host, ip: value.ip})
					}
					result.Urls = append(result.Urls, host)
				}
				// 测绘中暴露的端口服务纳入web扫描,覆盖已知子域名的非常规端口
				if value.port != "" && value.port != "0" {
					result.Urls = append(result.Urls, host+":"+value.port)
				}
			} else {
				if _, dup := addDomains[value.domain]; !dup {
					gologger.Infof("查询到同ip段内域名: %s", value.domain)
					addDomains[value.domain] = struct{}{}
					result.AddDomains = append(result.AddDomains, &DResult{
						Domain: value.domain,
						Ip:     value.ip,
					})
				}
			}
		} else {
			// IP记录
			if mode == certMode {
				// 证书发现的IP将作为扫描目标,过滤CDN节点
				if subfind.IsCdn("", []string{value.ip}) == false {
					gologger.Infof("查询到IP绑定目标证书: %s", value.ip)
					result.Urls = append(result.Urls, host+":"+value.port)
					result.Ips = append(result.Ips, value.ip)
				}
			} else {
				// IP检索的输入已是扫描目标IP,直接收集端口服务
				result.Urls = append(result.Urls, host+":"+value.port)
			}
		}
	}

	for _, pair := range subdomainOrder {
		result.AddSubdomains = append(result.AddSubdomains, &SdResult{
			Domain: pair.host,
			Ip:     pair.ip,
			CDN:    subfind.IsCdn("", []string{pair.ip}),
			CName:  "",
		})
	}

	result.Ips = util.RemoveRepeatedStringElement(result.Ips)
	result.Urls = util.RemoveRepeatedStringElement(result.Urls)
	return &result
}
