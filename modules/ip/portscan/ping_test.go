package portscan

import "testing"

// 实测两级ICMP降级在当前(非root)环境下的实际行为,
// go test -v 可观察特权/非特权ICMP各自的真实结果
func TestPingFallbackChain(t *testing.T) {
	priv := rawPing("127.0.0.1", true)
	unpriv := rawPing("127.0.0.1", false)
	t.Logf("特权ICMP=%v 非特权ICMP=%v", priv, unpriv)
	if !icmpLive("127.0.0.1") {
		t.Error("回环地址应判存活")
	}
}
