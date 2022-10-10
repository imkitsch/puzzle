package portscan

import (
	"fmt"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	"puzzle/gologger"
	"puzzle/modules/ip/portscan/scan"
	"time"
)

// NewRunner 通过解析配置选项、配置来源、阅读列表等创建一个新的 runner 结构体实例
func (o *Options) NewRunner(options *Options) (*Runner, error) {
	runner := &Runner{
		Options: options,
	}
	Scanner, err := scan.NewScanner(&scan.Options{
		Timeout: time.Duration(options.Timeout) * time.Millisecond,
		Retries: options.Retries,
		Rate:    options.Rate,
		Debug:   options.Debug,
		Proxy:   options.Proxy,
	})
	if err != nil {
		return nil, err
	}
	runner.Scanner = Scanner

	dnsOptions := dnsx.DefaultOptions
	dnsOptions.MaxRetries = runner.Options.Retries
	dnsOptions.Hostsfile = true

	dnsClient, err := dnsx.New(dnsOptions)
	if err != nil {
		return nil, err
	}
	runner.dnsClient = dnsClient

	// 解析扫描目标
	err = runner.ParseTarget()
	if err != nil {
		gologger.Fatalf("parse target failed: %v", err)
	}

	// 解析扫描端口
	runner.Scanner.Ports, err = ParsePorts(options)
	if err != nil {
		return nil, fmt.Errorf("could not parse ports: %s", err)
	}

	return runner, nil
}
