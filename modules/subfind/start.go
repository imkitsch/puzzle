package subfind

import (
	"Allin/core"
	"Allin/gologger"
	_ "github.com/google/gopacket/pcap"
)

func PrintStatus() {
	gologger.Printf("\rSuccess:%d Sent:%d Recved:%d Faild:%d", SuccessIndex, SentIndex, RecvIndex, FaildIndex)
}
func Start(options *core.Options) {
	//version := pcap.Version()

}
