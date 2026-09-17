package portscan

import (
	"puzzle/gologger"
	"sync"
	"time"

	pb "github.com/cheggaaa/pb/v3"
	"github.com/projectdiscovery/blackrock"
	"github.com/remeh/sizedwaitgroup"
	"go.uber.org/ratelimit"
)

type Options struct {
	Hosts     []string
	PortRange string
	Threads   int
}

type Runner struct {
	ips     []string
	ports   []int
	Rate    int
	LiveIPs *LiveIPs
	limiter ratelimit.Limiter
	wgScan  sizedwaitgroup.SizedWaitGroup
}

func NewRunner(options *Options) *Runner {
	live := make(map[string]struct{})
	ports, err := ParsePorts(options.PortRange)
	if err != nil {
		panic(err)
	}
	return &Runner{
		ips:     options.Hosts,
		ports:   ports,
		Rate:    options.Threads,
		LiveIPs: &LiveIPs{IPS: live},
	}
}

func (r *Runner) Run() (result []*Result) {
	gologger.Infof("开始端口扫描")

	// sizedwaitgroup.New 最大允许启动的goroutine数量
	// rate-limit.New 限制单位时间访问的频率
	r.wgScan = sizedwaitgroup.New(r.Rate)
	r.limiter = ratelimit.New(r.Rate)
	timeout := getTimeout(len(r.ports))

	portsCount := uint64(len(r.ports))
	ipsCount := uint64(len(r.ips))
	Range := portsCount * ipsCount

	bar := pb.StartNew(int(Range))
	currentSeed := time.Now().UnixNano()

	var mu sync.Mutex
	b := blackrock.New(int64(Range), currentSeed)
	for index := int64(0); index < int64(Range); index++ {
		xxx := b.Shuffle(index)
		ipIndex := int(xxx / int64(portsCount))
		portIndex := int(xxx % int64(portsCount))
		ip := r.PickIP(ipIndex)
		port := r.PickPort(portIndex)
		r.limiter.Take()
		r.wgScan.Add()
		go func() {
			defer bar.Increment()
			res, flag := r.Scan(ip, port, timeout)
			if flag == true {
				r.LiveIPs.AddIP(ip)
				gologger.Infof("%s:%d,finger:%v", ip, port, *res)
				mu.Lock()
				result = append(result, res)
				mu.Unlock()
			}
		}()
	}
	r.wgScan.Wait()
	bar.Finish()
	gologger.Infof("当前共 %v IP 存活", len(r.LiveIPs.IPS))
	gologger.Infof("端口扫描完毕")
	return
}

// PickPort 通过算法随机获取端口
func (r *Runner) PickPort(index int) int {
	return r.ports[index]
}

func (r *Runner) PickIP(index int) string {
	return r.ips[index]
}

func getTimeout(i int) time.Duration {
	switch {
	case i > 10000:
		return time.Millisecond * 200
	case i > 5000:
		return time.Millisecond * 300
	case i > 1000:
		return time.Millisecond * 400
	default:
		return time.Millisecond * 500
	}
}
