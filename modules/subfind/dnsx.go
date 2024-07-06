package subfind

import (
	"github.com/miekg/dns"
	"github.com/projectdiscovery/dnsx/libs/dnsx"
	"math"
)

func DomainBlastx(domains []string) {
	options := dnsx.Options{
		BaseResolvers:     []string{},
		MaxRetries:        5,
		QuestionTypes:     []uint16{dns.TypeA, dns.TypeCNAME, dns.TypeNS},
		TraceMaxRecursion: math.MaxUint16,
		Hostsfile:         true,
		OutputCDN:         true,
	}
	dnsClient, err := dnsx.New(options)
	if err != nil {

	}
	_ = dnsClient
	//dnsClient.QueryOne()
}
