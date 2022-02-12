package subfind

import (
	"Allin/core"
	"Allin/gologger"
	"Allin/util"
	"bufio"
	"github.com/google/gopacket/pcap"
	"io"
	"os"
	"strings"
)

func PrintStatus() {
	gologger.Printf("\rSuccess:%d Sent:%d Recved:%d Faild:%d", SuccessIndex, SentIndex, RecvIndex, FaildIndex)
}
func Start(options *core.Options) {
	version := pcap.Version()
	gologger.Infof(version + "\n")
	var ether EthTable
	if options.ListNetwork {
		GetIpv4Devices()
		os.Exit(0)
	}
	if options.NetworkId == -1 {
		ether = AutoGetDevices()
	} else {
		ether = GetDevices(options)
	}
	LocalStack = NewStack()

	// 设定接收的ID
	flagID := uint16(util.RandInt64(400, 654))
	retryChan := make(chan RetryStruct, options.Rate)
	go Recv(ether.Device, options, flagID, retryChan)
	sendog := SendDog{}
	sendog.Init(ether, options.Resolvers, flagID, true)

	var _ io.Reader
	// handle Stdin
	if options.Stdin {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			options.Domain = append(options.Domain, scanner.Text())
		}
	}

	// handle dict
	if len(options.Domain) > 0 {
		if options.FileName == "" {
			gologger.Infof("加载内置字典\n")
			_ = strings.NewReader(GetSubdomainData())
		} else {
			f2, err := os.Open(options.FileName)
			defer f2.Close()
			if err != nil {
				gologger.Fatalf("打开文件:%s 出现错误:%s\n", options.FileName, err.Error())
			}
			_ = f2
		}
	}

}
