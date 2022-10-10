package portscan

import (
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	"github.com/remeh/sizedwaitgroup"
	"go.uber.org/ratelimit"
	"puzzle/modules/ip/portscan/scan"
)

type Runner struct {
	Options   *Options
	Scanner   *scan.Scanner
	limiter   ratelimit.Limiter
	wgScan    sizedwaitgroup.SizedWaitGroup
	dnsClient *dnsx.DNSX
}
