package subfind

import (
	"net"
	"testing"
)

// 以下用例不依赖qqwry与CDN段数据文件,可在任意工作目录运行
func TestAreaHit(t *testing.T) {
	positive := []string{
		"某CDN节点", "CloudFlare节点", "网宿科技", "白山云", "百度云加速", "蓝汛", "帝联网络",
	}
	for _, area := range positive {
		if !areaHit(area) {
			t.Errorf("%q 应命中CDN关键字", area)
		}
	}
	negative := []string{"杭州市", "北京市 海淀区", "电信", "联通机房"}
	for _, area := range negative {
		if areaHit(area) {
			t.Errorf("%q 不应命中CDN关键字", area)
		}
	}
}

func TestIsCdnByCnameChain(t *testing.T) {
	// 多级CNAME链拼接后,中间节点的CDN特征域应命中
	if !IsCdn("www.example.com,static.akamai.net,final.example.net", nil) {
		t.Error("链中包含akamai.net应判定为CDN")
	}
	// 单级CNAME命中
	if !IsCdn("x.cdn.cloudflare.net", nil) {
		t.Error("cloudflare特征域应判定为CDN")
	}
	// 普通CNAME不误判
	if IsCdn("plain.example.net", nil) {
		t.Error("普通CNAME不应判定为CDN")
	}
}

func TestIsCdnMultiASubnet(t *testing.T) {
	// 该用例经过ipDecision->qqwry,需在仓库根目录运行(见文件尾说明)
	// 同/16网段的4条A记录(轮询解析)不应误判
	same := []string{"192.0.2.1", "192.0.2.2", "192.0.2.3", "192.0.2.4"}
	for _, ip := range same {
		if cidrCheck(ip) {
			t.Skip("192.0.2.0/24命中CDN段,环境数据影响用例")
		}
	}
	if IsCdn("", same) {
		t.Error("同/16网段的多A记录不应判定为CDN")
	}
	// 跨/16网段的4条A记录应判定为CDN
	cross := []string{"192.0.2.1", "192.0.2.2", "198.51.100.1", "198.51.100.2"}
	if !IsCdn("", cross) {
		t.Error("跨/16网段的多A记录应判定为CDN")
	}
}

func TestIpDecisionCache(t *testing.T) {
	// 需在仓库根目录运行:ipDecision依赖qqwry数据
	ip := "203.0.113.1"
	cidr1, area1 := ipDecision(ip)
	_, ok := ipCdnCache.Load(ip)
	if !ok {
		t.Fatal("首次判定后应写入缓存")
	}
	cidr2, area2 := ipDecision(ip)
	if cidr1 != cidr2 || area1 != area2 {
		t.Errorf("缓存结果应一致: (%v,%v) vs (%v,%v)", cidr1, area1, cidr2, area2)
	}
	// 确认判定对公网IP可用(不panic即通过)
	_ = net.ParseIP(ip)
}
