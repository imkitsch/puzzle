package portscan

import (
	"sync"
	"time"

	"github.com/go-ping/ping"
	"puzzle/gologger"
)

// icmpLive 主机存活探测,仅使用ICMP ping,不做TCP探测兜底:
// 优先特权模式(raw socket,需root),失败时降级为非特权模式
// (内核ping socket,受 net.ipv4.ping_group_range 限制)
func icmpLive(host string) bool {
	if rawPing(host, true) || rawPing(host, false) {
		gologger.Printf("发现存活主机:%s", host)
		return true
	}
	return false
}

func rawPing(host string, privileged bool) bool {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return false
	}
	pinger.Timeout = time.Second
	pinger.SetPrivileged(privileged)
	pinger.Count = 3
	if err := pinger.Run(); err != nil {
		gologger.Debugf("ping %s (privileged=%v) 失败: %v", host, privileged, err)
		return false
	}
	return pinger.Statistics().PacketsRecv >= 1
}

func Ping(host []string) []string {
	var LiveIp []string
	limit := 600
	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, limit)
	for _, ip := range host {
		wg.Add(1)
		sem <- struct{}{}
		go func(ip string) {
			defer wg.Done()
			defer func() { <-sem }()
			if icmpLive(ip) {
				mu.Lock()
				LiveIp = append(LiveIp, ip)
				mu.Unlock()
			}
		}(ip)
	}
	wg.Wait()

	return LiveIp
}
