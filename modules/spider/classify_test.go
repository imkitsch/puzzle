package spider

import "testing"

// classify用例中的CDN判定依赖qqwry数据,
// 需在仓库根目录运行: go test -c ./modules/spider -o /tmp/spider.test && /tmp/spider.test
func TestClassifyCertMode(t *testing.T) {
	sets := newClassifySets([]string{"example.com"}, []string{"known.example.com"})
	results := []*searchResult{
		// 遗漏子域名:根域在目标内,host不在已知子域
		{host: "miss.example.com", ip: "192.0.2.1", port: "443", domain: "example.com", protocol: "tls"},
		// 已知子域名:仅端口服务进Urls
		{host: "known.example.com", ip: "192.0.2.2", port: "8080", domain: "example.com", protocol: "http"},
		// 同证书域名:根域不在目标内,仅进报告
		{host: "other.example.net", ip: "192.0.2.3", port: "443", domain: "example.net", protocol: "tls"},
		// 非web协议:跳过
		{host: "skip.example.com", ip: "192.0.2.4", port: "22", domain: "example.com", protocol: "ssh"},
	}
	res := classify(results, sets, certMode)

	if len(res.AddSubdomains) != 1 || res.AddSubdomains[0].Domain != "miss.example.com" {
		t.Errorf("遗漏子域名分类不符: %+v", res.AddSubdomains)
	}
	if len(res.AddDomains) != 1 || res.AddDomains[0].Domain != "example.net" {
		t.Errorf("同证书域名分类不符: %+v", res.AddDomains)
	}
	urlSet := make(map[string]struct{})
	for _, u := range res.Urls {
		urlSet[u] = struct{}{}
	}
	for _, want := range []string{"miss.example.com", "miss.example.com:443", "known.example.com:8080"} {
		if _, ok := urlSet[want]; !ok {
			t.Errorf("Urls应包含 %s: %v", want, res.Urls)
		}
	}
	if _, ok := urlSet["skip.example.com:22"]; ok {
		t.Error("非web协议记录不应进入Urls")
	}
}

func TestClassifyCertModeIPRecord(t *testing.T) {
	// TEST-NET地址非CDN,应作为扫描目标产出(需在仓库根目录运行)
	sets := newClassifySets([]string{"example.com"}, nil)
	results := []*searchResult{
		{host: "192.0.2.10", ip: "192.0.2.10", port: "8443", domain: "", protocol: "tls"},
	}
	res := classify(results, sets, certMode)
	if len(res.Ips) != 1 || res.Ips[0] != "192.0.2.10" {
		t.Errorf("非CDN证书IP应进入Ips: %+v", res.Ips)
	}
	if len(res.Urls) != 1 || res.Urls[0] != "192.0.2.10:8443" {
		t.Errorf("IP记录Url不符: %v", res.Urls)
	}
}

func TestClassifyIpModeIPRecord(t *testing.T) {
	// ipMode输入已是扫描目标IP,不做CDN过滤,不产出Ips
	sets := newClassifySets([]string{"example.com"}, nil)
	results := []*searchResult{
		{host: "192.0.2.10", ip: "192.0.2.10", port: "8443", domain: "", protocol: "tls"},
	}
	res := classify(results, sets, ipMode)
	if len(res.Ips) != 0 {
		t.Errorf("ipMode不应产出Ips: %+v", res.Ips)
	}
	if len(res.Urls) != 1 || res.Urls[0] != "192.0.2.10:8443" {
		t.Errorf("IP记录Url不符: %v", res.Urls)
	}
}

func TestClassifyDedup(t *testing.T) {
	sets := newClassifySets([]string{"example.com"}, nil)
	results := []*searchResult{
		{host: "dup.example.com", ip: "192.0.2.1", port: "443", domain: "example.com", protocol: "tls"},
		{host: "dup.example.com", ip: "192.0.2.1", port: "443", domain: "example.com", protocol: "tls"},
	}
	res := classify(results, sets, certMode)
	if len(res.AddSubdomains) != 1 {
		t.Errorf("重复的遗漏子域对应去重: %+v", res.AddSubdomains)
	}
}
