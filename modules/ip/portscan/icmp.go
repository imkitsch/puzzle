package portscan

import (
	"github.com/go-ping/ping"
	"puzzle/gologger"
	"sync"
	"time"
)

func icmpLive(host string) bool {
	timeout := 1000
	pinger, err := ping.NewPinger(host)
	pingTimeOut := time.Duration(timeout)
	pinger.Timeout = time.Duration(pingTimeOut * time.Millisecond)
	pinger.SetPrivileged(true)
	if err != nil {
		return false
	}
	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		return false
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv >= 1 {
		gologger.Printf("发现存活主机:%s", host)
		return true
	}
	return false
}

func Ping(host []string) []string {
	var LiveIp []string
	limit := 600
	out := make(chan string, len(host))
	taskChan := make(chan bool, limit)
	defer close(taskChan)
	var wg sync.WaitGroup
	for _, ip := range host {
		wg.Add(1)
		if icmpLive(ip) {
			out <- ip
		}
		taskChan <- true
		go func() {
			<-taskChan
			defer wg.Done()
		}()
	}
	wg.Wait()

	close(out)
	for i := range out {
		LiveIp = append(LiveIp, i)
	}

	return LiveIp
}
