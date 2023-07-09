package subfind

import (
	"context"
	"github.com/boy-hack/ksubdomain/core/device"
	"github.com/boy-hack/ksubdomain/core/options"
	"github.com/boy-hack/ksubdomain/runner"
	"github.com/boy-hack/ksubdomain/runner/outputter"
	"puzzle/gologger"
	"time"
)

func DomainBlast(domains []string, DeviceConfig *device.EtherTable) []*domainResult {
	buffPrinter, _ := NewDomainResult()
	domainChanel := make(chan string)
	go func() {
		for _, d := range domains {
			domainChanel <- d
		}
		close(domainChanel)
	}()
	opt := &options.Options{
		Rate:        options.Band2Rate("1m"),
		Domain:      domainChanel,
		DomainTotal: len(domains),
		Resolvers:   options.GetResolvers(""),
		Silent:      false,
		TimeOut:     12,
		Retry:       3,
		Method:      runner.VerifyType,
		DnsType:     "a",
		Writer: []outputter.Output{
			buffPrinter,
		},
		EtherInfo: DeviceConfig,
	}
	r, err := runner.New(opt)
	if err != nil {
		gologger.Fatalf(err.Error())
	}
	ctx := context.Background()
	r.RunEnumeration(ctx)
	time.Sleep(10 * time.Second)
	r.Close()
	return buffPrinter.OutPut()
}

func GetDeviceConfig() *device.EtherTable {
	ether := device.AutoGetDevices()
	gologger.Infof("获取配置成功!")
	return ether
}
