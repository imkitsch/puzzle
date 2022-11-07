package portscan

import (
	"github.com/lcvvvv/gonmap"
	"time"
)

func (r *Runner) Scan(ip string, port int) (*Result, bool) {
	defer r.wgScan.Done()
	timeout := time.Millisecond * 1500
	r.nmapRunner.SetTimeout(timeout)
	status, response := r.nmapRunner.ScanTimeout(ip, port, 100*timeout)
	switch status {
	case gonmap.Closed:
		return nil, false
	case gonmap.Open:
		protocol := gonmap.GuessProtocol(port)
		return &Result{
			Addr:        ip,
			Port:        port,
			ServiceName: protocol,
		}, true
	case gonmap.Matched:
		return &Result{
			Addr:        ip,
			Port:        port,
			ServiceName: response.FingerPrint.Service,
			ProbeName:   response.FingerPrint.ProbeName,
			ProductName: response.FingerPrint.ProductName,
			Version:     response.FingerPrint.Version,
		}, true
	case gonmap.NotMatched:
		return &Result{
			Addr:        ip,
			Port:        port,
			ServiceName: "unknown",
			Raw:         response.Raw,
		}, true
	}
	return nil, false
}
