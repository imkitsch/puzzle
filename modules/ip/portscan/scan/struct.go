package scan

import (
	"github.com/google/gopacket"
	"github.com/projectdiscovery/ipranger"
	"golang.org/x/net/proxy"
	"net"
	"time"
)

// PkgSend is a TCP package
type PkgSend struct {
	ip       string
	port     int
	flag     PkgFlag
	SourceIP string
}

// PkgResult contains the results of sending TCP packages
type PkgResult struct {
	ip   string
	port int
}

type Scanner struct {
	SourceIP           net.IP
	NetworkInterface   *net.Interface
	tcpPacketListener  net.PacketConn
	icmpPacketListener net.PacketConn
	retries            int
	rate               int
	listenPort         int
	timeout            time.Duration
	proxyDialer        proxy.Dialer
	Ports              []int
	IPRanger           *ipranger.IPRanger
	tcpPacketSend      chan *PkgSend
	icmpPacketSend     chan *PkgSend
	tcpChan            chan *PkgResult
	icmpChan           chan *PkgResult
	State              State
	ScanResults        *Result
	tcpSequencer       *TCPSequencer
	serializeOptions   gopacket.SerializeOptions
	debug              bool
	handlers           interface{}
}

type Options struct {
	Timeout time.Duration
	Retries int
	Rate    int
	Debug   bool
	Proxy   string
	Stream  bool
}
