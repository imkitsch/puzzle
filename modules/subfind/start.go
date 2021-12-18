package subfind

import (
	"Allin/core"
	"Allin/gologger"
	"Allin/util"
	"bufio"
	"github.com/google/gopacket/pcap"
	"io"
	"os"
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

	var f io.Reader
	// handle Stdin
	if options.Stdin {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			options.Domain = append(options.Domain, scanner.Text())
		}
	}

}
