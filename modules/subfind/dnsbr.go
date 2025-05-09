package subfind

import (
	"context"
	"fmt"
	"github.com/boy-hack/ksubdomain/core/device"
	"github.com/boy-hack/ksubdomain/core/options"
	"github.com/boy-hack/ksubdomain/runner"
	"github.com/boy-hack/ksubdomain/runner/outputter"
	"path/filepath"
	"puzzle/config"
	"puzzle/gologger"
	"puzzle/util"
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
		Resolvers:   []string{"8.8.8.8", "8.8.4.4"},
		Silent:      false,
		TimeOut:     6,
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
	fmt.Println("")
	gologger.Infof("获取域名数量: %d", len(buffPrinter.OutPut()))
	return buffPrinter.OutPut()
}

func GetDeviceConfig() *device.EtherTable {
	var ether *device.EtherTable
	deviceFile := filepath.Join(util.GetRunDir() + config.DeviceConfig)
	if util.FileExists(deviceFile) == true {
		var err error
		ether, err = device.ReadConfig(deviceFile)
		if err != nil {
			gologger.Fatalf("读取配置失败:%v", err)
		}
		gologger.Infof("读取配置%s成功!", deviceFile)
	} else {
		ether = device.AutoGetDevices()
		gologger.Infof("获取配置成功!")
	}
	return ether
}
