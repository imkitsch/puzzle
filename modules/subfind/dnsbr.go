package subfind

import (
	"context"
	"github.com/boy-hack/ksubdomain/core/device"
	"github.com/boy-hack/ksubdomain/core/options"
	"github.com/boy-hack/ksubdomain/runner"
	"github.com/boy-hack/ksubdomain/runner/outputter"
	"puzzle/gologger"
	"puzzle/util"
)

func DomainBlast(domains []string) []domainResult {
	buffPrinter, _ := NewDomainResult()

	domainChanel := make(chan string)
	go func() {
		for _, d := range domains {
			domainChanel <- d
		}
		close(domainChanel)
	}()
	opt := &options.Options{
		Rate:        options.Band2Rate("2m"),
		Domain:      domainChanel,
		DomainTotal: len(domains),
		Resolvers:   options.GetResolvers(""),
		Silent:      false,
		TimeOut:     10,
		Retry:       3,
		Method:      runner.VerifyType,
		DnsType:     "a",
		Writer: []outputter.Output{
			buffPrinter,
		},
		EtherInfo: getDeviceConfig(),
	}
	r, err := runner.New(opt)
	if err != nil {
		gologger.Fatalf(err.Error())
	}
	ctx := context.Background()
	r.RunEnumeration(ctx)
	r.Close()

	return buffPrinter.OutPut()
}

func getDeviceConfig() *device.EtherTable {
	filename := "config/device.yaml"
	var ether *device.EtherTable
	var err error
	if util.FileExists(filename) {
		ether, err = device.ReadConfig(filename)
		if err != nil {
			gologger.Fatalf("读取配置失败:%v", err)
		}
		gologger.Infof("读取配置%s成功!", filename)
	} else {
		ether = device.AutoGetDevices()
		err = ether.SaveConfig(filename)
		if err != nil {
			gologger.Fatalf("保存配置失败:%v", err)
		}
	}
	return ether
}
