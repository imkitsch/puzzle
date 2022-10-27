package ip

import (
	"puzzle/gologger"
	"puzzle/modules/ip/portfinger"
	"puzzle/modules/ip/portscan"
)

type Options struct {
	Hosts     []string
	PortRange string
	Threads   int
	Rate      int
	MaxPort   int
	QQwry     *QQwry
	ScanType  string
	NmapProbe *portfinger.NmapProbe
}

type Runner struct {
	options          *Options
	scanRunner       *portscan.Runner
	portFingerEngine *portfinger.Engine
}

func NewRunner(options *Options) (*Runner, error) {
	scanOptions := &portscan.Options{
		Host:      options.Hosts,
		Ports:     options.PortRange,
		Threads:   options.Threads,
		Timeout:   portscan.DefaultPortTimeoutSynScan,
		ScanType:  options.ScanType,
		Rate:      options.Rate,
		Retries:   portscan.DefaultRetriesSynScan,
		Interface: "",
		Debug:     false,
	}
	scanRunner, err := scanOptions.NewRunner(scanOptions)
	if err != nil {
		return nil, err
	}
	portFingerEngine, err := portfinger.NewEngine(200, options.NmapProbe)
	if err != nil {
		return nil, err
	}
	return &Runner{
		options:          options,
		scanRunner:       scanRunner,
		portFingerEngine: portFingerEngine,
	}, nil
}

func (r *Runner) Run() (results []*portfinger.Result) {
	gologger.Infof("开始端口扫描")
	err := r.scanRunner.Run()
	if err != nil {
		gologger.Errorf("scanRunner.Run() err, %v", err)
		return
	}
	portscanResult := r.scanRunner.Scanner.ScanResults.IPPorts
	if len(portscanResult) == 0 {
		return
	}
	// 去除开放端口数大于maxPort
	for k := range portscanResult {
		ports := portscanResult[k]
		if len(ports) > r.options.MaxPort {
			gologger.Infof("%v 开放端口大于 %v", k, r.options.MaxPort)
			portscanResult[k] = map[int]struct{}{}
		}
	}
	gologger.Infof("端口协议识别")
	results = r.portFingerEngine.Run(portscanResult)

	return
}
