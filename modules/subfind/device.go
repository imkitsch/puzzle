package subfind

import (
	"Allin/core"
	"Allin/gologger"
	"Allin/util"
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
	"time"
)

func AutoGetDevices() EthTable {
	domain := util.RandomStr(4) + "paper.seebug.org"
	signal := make(chan EthTable)
	devices, err := pcap.FindAllDevs()
	if err != nil {
		gologger.Fatalf("获取网络设备失败:%s\n", err.Error())
	}
	data := make(map[string]net.IP)
	keys := []string{}
	for _, d := range devices {
		for _, address := range d.Addresses {
			ip := address.IP
			if ip.To4() != nil && !ip.IsLoopback() {
				data[d.Name] = ip
				keys = append(keys, d.Name)
			}
		}
	}
	ctx := context.Background()
	// 在初始上下文的基础上创建一个有取消功能的上下文
	ctx, cancel := context.WithCancel(ctx)
	for _, drviceName := range keys {
		go func(drviceName string, domain string, signal chan EthTable, ctx context.Context) {
			var (
				snapshot_len int32         = 1024
				promiscuous  bool          = false
				timeout      time.Duration = -1 * time.Second
				handle       *pcap.Handle
			)
			var err error
			handle, err = pcap.OpenLive(
				drviceName,
				snapshot_len,
				promiscuous,
				timeout,
			)
			if err != nil {
				gologger.Fatalf("pcap打开失败:%s\n", err.Error())
			}
			defer handle.Close()
			// Use the handle as a packet source to process all packets
			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			for {
				select {
				case <-ctx.Done():
					return
				default:
					packet, err := packetSource.NextPacket()
					gologger.Printf(".")
					if err != nil {
						continue
					}
					if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
						dns, _ := dnsLayer.(*layers.DNS)
						if !dns.QR {
							continue
						}
						for _, v := range dns.Questions {
							if string(v.Name) == domain {
								ethLayer := packet.Layer(layers.LayerTypeEthernet)
								if ethLayer != nil {
									eth := ethLayer.(*layers.Ethernet)
									signal <- EthTable{SrcIp: data[drviceName], Device: drviceName, SrcMac: eth.DstMAC, DstMac: eth.SrcMAC}
									// 网关mac 和 本地mac
									return
								}
							}
						}
					}
				}
			}
		}(drviceName, domain, signal, ctx)
	}
	var c EthTable
	for {
		select {
		case c = <-signal:
			cancel()
			goto END
		default:
			_, _ = net.LookupHost(domain)
			time.Sleep(time.Millisecond * 20)
		}
	}
END:
	gologger.Printf("\n")
	gologger.Infof("Use Device: %s\n", c.Device)
	gologger.Infof("Use IP:%s\n", c.SrcIp.String())
	gologger.Infof("Local Mac:%s\n", c.SrcMac.String())
	gologger.Infof("GateWay Mac:%s\n", c.DstMac.String())
	return c
}
func GetIpv4Devices() (keys []string, data map[string]net.IP) {
	devices, err := pcap.FindAllDevs()
	data = make(map[string]net.IP)
	if err != nil {
		gologger.Fatalf("获取网络设备失败:%s\n", err.Error())
	}
	for _, d := range devices {
		for _, address := range d.Addresses {
			ip := address.IP
			if ip.To4() != nil && !ip.IsLoopback() {
				gologger.Printf("  [%d] Name: %s\n", len(keys), d.Name)
				gologger.Printf("  Description: %s\n", d.Description)
				gologger.Printf("  Devices addresses: %s\n", d.Description)
				gologger.Printf("  IP address: %s\n", ip)
				gologger.Printf("  Subnet mask: %s\n\n", address.Netmask.String())
				data[d.Name] = ip
				keys = append(keys, d.Name)
			}
		}
	}
	return
}
func GetDevices(options *core.Options) EthTable {
	// Find all devices
	defaultSelect := options.NetworkId
	keys, data := GetIpv4Devices()

	if len(keys) == 0 {
		gologger.Fatalf("获取不到可用的设备名称\n")
	} else if len(keys) == 1 {
		defaultSelect = 0
	}
	if defaultSelect == -1 {
		if options.Stdin {
			gologger.Fatalf("Stdin模式下需要指定-e参数\n")
		}
		gologger.Infof("选择一个可用网卡ID:")
		var i int
		_, err2 := fmt.Scanln(&i)
		if err2 != nil {
			gologger.Fatalf("\n读入ID失败，确认输入的是数字?\n")
		}

		if i < 0 || i >= len(keys) {
			gologger.Fatalf("ID超出了范围\n")
		}
		defaultSelect = i
	}
	deviceName := keys[defaultSelect]
	ip := data[deviceName]
	gologger.Infof("Use Device: %s\n", deviceName)
	gologger.Infof("Use IP:%s\n", ip.String())
	c := GetGateMacAddress(deviceName)
	gologger.Infof("Local Mac:%s\n", c[1])
	gologger.Infof("GateWay Mac:%s\n", c[0])
	return EthTable{ip, deviceName, c[1], c[0]}
}

func GetGateMacAddress(dvice string) [2]net.HardwareAddr {
	// 获取网关mac地址
	domain := util.RandomStr(4) + "paper.seebug.org"
	_signal := make(chan [2]net.HardwareAddr)
	go func(device string, domain string, signal chan [2]net.HardwareAddr) {
		var (
			snapshot_len int32         = 1024
			promiscuous  bool          = false
			timeout      time.Duration = -1 * time.Second
			handle       *pcap.Handle
		)
		var err error
		handle, err = pcap.OpenLive(
			device,
			snapshot_len,
			promiscuous,
			timeout,
		)
		if err != nil {
			gologger.Fatalf("pcap打开失败:%s\n", err.Error())
		}
		defer handle.Close()
		// Use the handle as a packet source to process all packets
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for {
			packet, err := packetSource.NextPacket()
			gologger.Printf(".")
			if err != nil {
				continue
			}
			if dnsLayer := packet.Layer(layers.LayerTypeDNS); dnsLayer != nil {
				dns, _ := dnsLayer.(*layers.DNS)
				if !dns.QR {
					continue
				}
				for _, v := range dns.Questions {
					if string(v.Name) == domain {
						ethLayer := packet.Layer(layers.LayerTypeEthernet)
						if ethLayer != nil {
							eth := ethLayer.(*layers.Ethernet)
							srcMac := eth.SrcMAC
							dstMac := eth.DstMAC

							signal <- [2]net.HardwareAddr{srcMac, dstMac}
							// 网关mac 和 本地mac
							return
						}
					}
				}
			}

		}
	}(dvice, domain, _signal)
	var c [2]net.HardwareAddr
	for {
		select {
		case c = <-_signal:
			gologger.Printf("\n")
			goto END
		default:
			_, _ = net.LookupHost(domain)
			time.Sleep(time.Millisecond * 10)
		}
	}
END:
	return c
}
