package portscan

import "sync"

type Result struct {
	Addr        string
	Port        int
	ServiceName string
	ProbeName   string
	ProductName string
	Version     string
	Raw         string
}

type LiveIPs struct {
	sync.RWMutex
	IPS map[string]struct{}
}

func (l *LiveIPs) AddIP(k string) {
	l.Lock()
	defer l.Unlock()
	if _, ok := l.IPS[k]; !ok {
		l.IPS[k] = struct{}{}
	}
}
