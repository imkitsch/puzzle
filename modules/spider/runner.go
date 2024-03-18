package spider

import (
	"encoding/base64"
	"net/url"
	"puzzle/gologger"
	"puzzle/util"
	"strings"
	"time"
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
	//以10为单位打包,减少查询次数
	//ip查询打包
	var raw []string
	var tmp []string
	for k, ip := range r.options.Ips {
		tmp = append(tmp, "ip=\""+ip+"\"")
		if (k+1)%10 == 0 || k == len(r.options.Ips)-1 {
			raw = append(raw, strings.Join(tmp, "||"))
			tmp = []string{}
		}
	}

	//证书查询打包
	tmp = []string{}
	for k, domain := range r.options.Domains {
		tmp = append(tmp, "cert=\""+domain+"\"")
		if (k+1)%10 == 0 || k == len(r.options.Domains)-1 {
			raw = append(raw, strings.Join(tmp, "||"))
			tmp = []string{}
		}
	}

	var fofaRes []*fofaResult

	// fofa查询
	for _, v := range raw {
		res := r.getFofaResult(base64.StdEncoding.EncodeToString([]byte(v)))
		if res != nil {
			fofaRes = append(fofaRes, res...)
		}
		time.Sleep(time.Second * 2)
	}

	gologger.Infof("从fofa获取信息")

	var addDomains []string
	var addSubdomains [][]string

	for _, value := range fofaRes {
		//添加url
		if value.protocol == "http" || value.protocol == "https" || value.protocol == "tls" || value.protocol == "unknown" {
			if value.host[0:4] != "http" {
				value.host = "http://" + value.host
			}

			tmpUrl, _ := url.Parse(value.host)
			host := tmpUrl.Hostname()

			// url为域名
			if value.domain != "" {
				if util.InSlice(r.options.Domains, value.domain) == true {
					if util.InSlice(r.options.Subdomains, host) == false {
						gologger.Infof("查询到遗漏子域名: %s", host)
						addSubdomains = append(addSubdomains, []string{host, value.ip})
						result.Urls = append(result.Urls, host)
					}
				} else {
					gologger.Infof("查询到同ip段内域名: %s", value.domain)
					addDomains = append(addDomains, value.domain)
				}
			} else {
				//为ip
				result.Urls = append(result.Urls, host+value.port)
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
