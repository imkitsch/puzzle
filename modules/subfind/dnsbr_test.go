package subfind

import (
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/miekg/dns"
)

// servailQueryCount 统计SERVFAIL域名的查询次数,用于验证重试行为
var servailQueryCount atomic.Int32

func TestGetDeviceConfig(t *testing.T) {
	cfg := GetDeviceConfig()
	if cfg == nil || len(cfg.Resolvers) == 0 {
		t.Fatal("GetDeviceConfig应返回默认DNS服务器列表")
	}
}

func TestParseAnswers(t *testing.T) {
	msg := new(dns.Msg)
	msg.Answer = []dns.RR{
		&dns.CNAME{
			Hdr:    dns.RR_Header{Name: "www.example.com.", Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 60},
			Target: "cname.example.com.",
		},
		&dns.A{
			Hdr: dns.RR_Header{Name: "cname.example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			A:   net.ParseIP("192.0.2.1"),
		},
	}
	res := parseAnswers("www.example.com", msg)
	if res == nil {
		t.Fatal("CNAME+A应答应解析出结果")
	}
	if res.Domain != "www.example.com" {
		t.Errorf("Domain = %s, want www.example.com", res.Domain)
	}
	if res.Cname != "cname.example.com" {
		t.Errorf("Cname = %s, want cname.example.com", res.Cname)
	}
	if res.Address != "192.0.2.1" {
		t.Errorf("Address = %s, want 192.0.2.1", res.Address)
	}

	// 空应答不应产生结果
	if parseAnswers("empty.example.com", new(dns.Msg)) != nil {
		t.Error("空应答不应产生结果")
	}
}

func TestParseAnswersCnameChain(t *testing.T) {
	// 两级CNAME链:中间节点命中CDN特征域,末级为普通域
	msg := new(dns.Msg)
	msg.Answer = []dns.RR{
		&dns.CNAME{
			Hdr:    dns.RR_Header{Name: "www.example.com.", Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 60},
			Target: "e123.x.akamaiedge.net.",
		},
		&dns.CNAME{
			Hdr:    dns.RR_Header{Name: "e123.x.akamaiedge.net.", Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 60},
			Target: "final.example.net.",
		},
		&dns.A{
			Hdr: dns.RR_Header{Name: "final.example.net.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			A:   net.ParseIP("192.0.2.1"),
		},
	}
	res := parseAnswers("www.example.com", msg)
	if res == nil {
		t.Fatal("CNAME链应答应解析出结果")
	}
	// 记录字段保留最后一级
	if res.Cname != "final.example.net" {
		t.Errorf("Cname = %s, want final.example.net", res.Cname)
	}
	// 中间节点的akamai特征域应使整链命中CDN
	if res.Cdn != true {
		t.Error("链中含akamaiedge.net应判定为CDN")
	}
}

func TestParseAnswersMultiA(t *testing.T) {
	msg := new(dns.Msg)
	for _, ip := range []string{"192.0.2.1", "192.0.2.2", "198.51.100.1", "198.51.100.2"} {
		msg.Answer = append(msg.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: "multi.example.com.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			A:   net.ParseIP(ip),
		})
	}
	res := parseAnswers("multi.example.com", msg)
	if res == nil {
		t.Fatal("多A记录应答应解析出结果")
	}
	want := "192.0.2.1,192.0.2.2,198.51.100.1,198.51.100.2"
	if res.Address != want {
		t.Errorf("Address = %s, want %s", res.Address, want)
	}
	// 跨/16网段的多A记录应命中CDN启发式
	if res.Cdn != true {
		t.Errorf("跨/16网段的4条A记录应判定为CDN")
	}
}

// startTestDNSServer 启动本地UDP DNS服务器模拟解析,避免测试依赖外部网络
func startTestDNSServer(tb testing.TB) string {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("127.0.0.1")})
	if err != nil {
		tb.Skipf("无法监听本地UDP端口: %v", err)
	}
	// 扩大收包缓冲,避免高并发下UDP丢包干扰基准测量
	_ = conn.SetReadBuffer(8 << 20)
	tb.Cleanup(func() { conn.Close() })

	go dns.ActivateAndServe(nil, conn, dns.HandlerFunc(dnsServeveHandler))
	return conn.LocalAddr().String()
}

func dnsServeveHandler(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	q := r.Question[0]
	name := strings.TrimSuffix(q.Name, ".")
	switch {
	case name == "a.test.example":
		m.Answer = append(m.Answer,
			&dns.CNAME{
				Hdr:    dns.RR_Header{Name: q.Name, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: 60},
				Target: "c.test.example.",
			},
			&dns.A{
				Hdr: dns.RR_Header{Name: "c.test.example.", Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP("192.0.2.10"),
			})
	case name == "b.test.example":
		m.Answer = append(m.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			A:   net.ParseIP("192.0.2.11"),
		})
	case strings.HasPrefix(name, "bench-"):
		m.Answer = append(m.Answer, &dns.A{
			Hdr: dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
			A:   net.ParseIP("192.0.2.1"),
		})
	case name == "servfail.test.example":
		m.Rcode = dns.RcodeServerFailure
		servailQueryCount.Add(1)
	default:
		m.Rcode = dns.RcodeNameError
	}
	_ = w.WriteMsg(m)
}

// TestDomainBlast 通过本地DNS服务器验证完整爆破流程,
// 需在仓库根目录运行(包内cdn.go/qqwry依赖工作目录下的config):
//	go test -c ./modules/subfind -o /tmp/subfind.test && /tmp/subfind.test
func TestDomainBlast(t *testing.T) {
	server := startTestDNSServer(t)
	cfg := &EtherTable{Resolvers: []string{server}}

	results := DomainBlast([]string{"a.test.example", "nx.test.example", "b.test.example"}, cfg)

	byName := map[string]*domainResult{}
	for _, r := range results {
		byName[r.Domain] = r
	}
	if len(results) != 2 {
		t.Fatalf("应解析出2个域名,实际%d个: %+v", len(results), results)
	}
	a := byName["a.test.example"]
	if a == nil {
		t.Fatal("a.test.example应被解析")
	}
	if a.Cname != "c.test.example" {
		t.Errorf("Cname = %s, want c.test.example", a.Cname)
	}
	if a.Address != "192.0.2.10" {
		t.Errorf("Address = %s, want 192.0.2.10", a.Address)
	}
	b := byName["b.test.example"]
	if b == nil || b.Address != "192.0.2.11" || b.Cname != "" {
		t.Errorf("b.test.example字段不符: %+v", b)
	}
	if _, ok := byName["nx.test.example"]; ok {
		t.Error("NXDOMAIN域名不应出现在结果中")
	}
}

// TestDomainBlastServfailRetry SERVFAIL应换DNS服务器重试dnsRetry次,且不产出结果
func TestDomainBlastServfailRetry(t *testing.T) {
	server := startTestDNSServer(t)
	cfg := &EtherTable{Resolvers: []string{server}}
	servailQueryCount.Store(0)

	results := DomainBlast([]string{"servfail.test.example"}, cfg)
	if len(results) != 0 {
		t.Errorf("SERVFAIL域名不应有结果: %+v", results)
	}
	if got := servailQueryCount.Load(); got != dnsRetry {
		t.Errorf("SERVFAIL应重试%d次,实际查询%d次", dnsRetry, got)
	}
}

// BenchmarkDomainBlast 度量5000个域名的端到端爆破耗时,评估与pcap实现的性能差距
func BenchmarkDomainBlast(b *testing.B) {
	server := startTestDNSServer(b)
	cfg := &EtherTable{Resolvers: []string{server}}
	domains := make([]string, 5000)
	for i := range domains {
		domains[i] = fmt.Sprintf("bench-%d.test.example", i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		results := DomainBlast(domains, cfg)
		if len(results) != len(domains) {
			b.Fatalf("应解析出%d个域名,实际%d个", len(domains), len(results))
		}
	}
}
