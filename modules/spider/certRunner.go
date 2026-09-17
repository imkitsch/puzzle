package spider

import (
	"puzzle/gologger"
)

// CertRun 从测绘引擎收集证书:
// 发现遗漏子域名、同证书域名,以及绑定目标证书的非CDN IP(作为后续扫描目标)
func (r *Runner) CertRun() *Result {
	gologger.Infof("从测绘引擎收集证书")
	if len(r.certEngines) == 0 {
		return &Result{}
	}

	var searchRes []*searchResult
	for _, engine := range r.certEngines {
		gologger.Infof("使用%s检索证书", engine.Name())
		searchRes = append(searchRes, engine.QueryCert(r.options.Domains)...)
	}

	gologger.Infof("从网络测绘获取信息")
	sets := newClassifySets(r.options.Domains, r.options.Subdomains)
	return classify(searchRes, sets, certMode)
}
