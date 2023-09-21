package spider

import (
	"encoding/base64"
	"github.com/imroc/req/v3"
	"puzzle/gologger"
	"puzzle/util"
	"strings"
)

func NewRunner(options *Options) *Runner {
	return &Runner{
		options:   options,
		reqClient: req.C(),
	}
}

func (r *Runner) Run() *Result {

	gologger.Infof("开始爬虫信息收集")

	var result Result
	//以10为单位打包ip,减少查询次数
	var raw []string
	var tmp []string
	for k, ip := range r.options.Ips {
		tmp = append(tmp, "ip=\""+ip+"\"")
		if (k+1)%10 == 0 || k == len(r.options.Ips)-1 {
			raw = append(raw, strings.Join(tmp, "||"))
			tmp = []string{}
		}
	}

	var fofaRes []*fofaResult
	for _, v := range raw {
		res := r.getFofaResult(base64.StdEncoding.EncodeToString([]byte(v)))
		if res != nil {
			fofaRes = append(fofaRes, res...)
		}
	}

	gologger.Infof("获取fofa信息")

	var addDomains []string
	var addSubdomains [][]string

	for _, value := range fofaRes {
		//添加url
		if value.protocol == "http" || value.protocol == "https" || value.protocol == "tls" || value.protocol == "unknown" {
			result.Urls = append(result.Urls, value.host)
		}

		//如果存在domain
		if value.domain != "" {
			rootDomain := domain_parse(value.domain)
			if util.InSlice(r.options.Domains, rootDomain) == true {
				if util.InSlice(r.options.Subdomains, value.domain) == false {
					gologger.Infof("查询到遗漏子域名: %s", value.domain)
					addSubdomains = append(addSubdomains, []string{value.domain, value.ip})
				}
			} else {
				gologger.Infof("查询到同ip段内域名: %s", rootDomain)
				addDomains = append(addDomains, rootDomain)
			}
		}

	}

	addDomains = util.RemoveRepeatedStringElement(addDomains)
	addSubdomains = util.RemoveRepeatedStringArrayElement(addSubdomains)

	for _, v := range addDomains {
		result.AddDomains = append(result.AddDomains, &DResult{
			Domain: v,
		})
	}

	for _, v := range addSubdomains {
		result.AddSubdomains = append(result.AddSubdomains, &SdResult{
			Domain: v[0],
			CDN:    false,
			Ip:     v[1],
			CName:  "",
		})
	}

	result.Urls = util.RemoveRepeatedStringElement(result.Urls)
	return &result
}
