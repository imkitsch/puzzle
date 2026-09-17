package spider

import (
	"puzzle/gologger"
)

// IpRun 从测绘引擎按已知IP收集:遗漏子域名、同证书域名与端口服务
func (r *Runner) IpRun() *Result {
	gologger.Infof("从测绘引擎收集IP信息")
	if len(r.ipEngines) == 0 {
		return &Result{}
	}

	var searchRes []*searchResult
	for _, engine := range r.ipEngines {
		gologger.Infof("使用%s检索IP", engine.Name())
		searchRes = append(searchRes, engine.QueryIP(r.options.Ips)...)
	}

	gologger.Infof("从网络测绘获取信息")
	sets := newClassifySets(r.options.Domains, r.options.Subdomains)
	return classify(searchRes, sets, ipMode)
}
