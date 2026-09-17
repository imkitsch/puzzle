package subfind

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"puzzle/gologger"
	"puzzle/util"

	"github.com/miekg/dns"
)

// 参考 gobuster 的 dns 模式实现:标准网络栈 + goroutine 工作池并发查询,
// 替换原 ksubdomain 的 pcap 收发,移除对 libpcap 的依赖。
// 查询参数(超时/重试/DNS服务器列表)与原实现保持一致。

const (
	dnsTimeout = 6 * time.Second
	dnsRetry   = 3
	dnsWorkers = 1000
)

// defaultResolvers 与原ksubdomain实现使用相同的DNS服务器列表
var defaultResolvers = []string{"8.8.8.8", "8.8.4.4", "223.5.5.5", "119.29.29.29", "114.114.114.114"}

// EtherTable 原ksubdomain网卡设备配置的占位类型,
// 现仅承载DNS服务器列表,用于保持对外函数签名兼容
type EtherTable struct {
	Resolvers []string
}

// GetDeviceConfig 返回DNS解析配置,保持原接口签名
func GetDeviceConfig() *EtherTable {
	return &EtherTable{
		Resolvers: defaultResolvers,
	}
}

// DomainBlast 对域名列表并发做A记录解析,返回可解析的域名及其CNAME/A记录,
// 输出结构(Domain/Cname/Address/Cdn)与原ksubdomain实现一致
func DomainBlast(domains []string, DeviceConfig *EtherTable) []*domainResult {
	resolvers := defaultResolvers
	if DeviceConfig != nil && len(DeviceConfig.Resolvers) != 0 {
		resolvers = DeviceConfig.Resolvers
	}

	gologger.Infof("开始DNS爆破,共 %d 个域名", len(domains))

	domainChanel := make(chan string)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var results []*domainResult

	// 工作池:每个worker从公共队列消费域名并发查询,并发数即限速手段;
	// 小任务时按任务量收缩,避免为几个域名也起满额goroutine
	workers := dnsWorkers
	if len(domains) < workers {
		workers = len(domains)
	}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			udpClient := &dns.Client{Timeout: dnsTimeout, UDPSize: dns.MaxMsgSize}
			tcpClient := &dns.Client{Net: "tcp", Timeout: dnsTimeout}
			for domain := range domainChanel {
				msg := queryDomain(udpClient, tcpClient, domain, resolvers)
				if msg == nil {
					continue
				}
				res := parseAnswers(domain, msg)
				if res == nil {
					continue
				}
				mu.Lock()
				results = append(results, res)
				mu.Unlock()
			}
		}()
	}

	for _, d := range domains {
		domainChanel <- d
	}
	close(domainChanel)
	wg.Wait()

	fmt.Println("")
	gologger.Infof("获取域名数量: %d", len(results))
	return results
}

// queryDomain 按重试次数随机选择DNS服务器发起A记录查询,返回应答消息;
// 网络错误时重试,明确的否定应答(NXDOMAIN等)不重试
func queryDomain(udpClient, tcpClient *dns.Client, domain string, resolvers []string) *dns.Msg {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	m.RecursionDesired = true

	for i := 0; i < dnsRetry; i++ {
		server := util.RandomPick(resolvers)
		if !strings.Contains(server, ":") {
			server = server + ":53"
		}
		resp, _, err := udpClient.Exchange(m, server)
		if err != nil {
			continue
		}
		// UDP应答被截断时改用TCP重查
		if resp.Truncated {
			resp, _, err = tcpClient.Exchange(m, server)
			if err != nil {
				continue
			}
		}
		if resp.Rcode == dns.RcodeSuccess {
			return resp
		}
		if resp.Rcode == dns.RcodeNameError {
			// NXDOMAIN为明确否定应答,无需重试
			return nil
		}
		// SERVFAIL/REFUSED等临时错误,换DNS服务器重试
	}
	return nil
}

// parseAnswers 从应答中提取CNAME与A记录,组装与原实现一致的domainResult;
// 一次A查询的应答 section 同时包含CNAME链与最终A记录,与ksubdomain输出等价
func parseAnswers(domain string, msg *dns.Msg) *domainResult {
	var cnames []string
	var addresses []string
	for _, rr := range msg.Answer {
		switch record := rr.(type) {
		case *dns.CNAME:
			cnames = append(cnames, strings.TrimSuffix(record.Target, "."))
		case *dns.A:
			addresses = append(addresses, record.A.String())
		}
	}
	if len(addresses) == 0 && len(cnames) == 0 {
		return nil
	}
	cname := ""
	if len(cnames) != 0 {
		// 记录字段保留最后一级,与原实现一致
		cname = cnames[len(cnames)-1]
	}
	// 多级CNAME链的中间节点常含CDN特征域,整链参与CDN判断
	return &domainResult{
		Domain:  domain,
		Cname:   cname,
		Address: strings.Join(addresses, ","),
		Cdn:     IsCdn(strings.Join(cnames, ","), addresses),
	}
}
