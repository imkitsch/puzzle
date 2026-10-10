package portscan

import (
	"fmt"
	"net"
	"puzzle/gologger"
	"strconv"
	"strings"
)

// ParseIP 参考 fscan
// 仅接受IPv4输入(单IP/IP段/CIDR),IPv6地址直接抛弃
func ParseIP(ip string) ([]string, error) {
	var hosts []string
	if strings.Contains(ip, "/") {
		_, ipNet, err := net.ParseCIDR(ip)
		if err != nil {
			return hosts, err
		}
		if ipNet.IP.To4() == nil {
			gologger.Warningf("忽略IPv6地址输入: %s", ip)
			return hosts, nil
		}
		hosts = parseIP1(IPRange(ipNet))
	} else if strings.Contains(ip, "-") {
		hosts = parseIP1(ip)
	} else {
		testIP := net.ParseIP(ip)
		if testIP == nil {
			gologger.Warningf("无法解析的IP输入: %s", ip)
		} else if testIP.To4() == nil {
			gologger.Warningf("忽略IPv6地址输入: %s", ip)
		} else {
			hosts = []string{testIP.To4().String()}
		}
	}
	return hosts, nil
}

// IPRange 把 192.168.x.x/xx 转换成 192.168.x.x-192.168.x.x
func IPRange(c *net.IPNet) string {
	start := c.IP.String()
	mask := c.Mask
	bcst := make(net.IP, len(c.IP))
	copy(bcst, c.IP)
	for i := 0; i < len(mask); i++ {
		ipIdx := len(bcst) - i - 1
		bcst[ipIdx] = c.IP[ipIdx] | ^mask[len(mask)-i-1]
	}
	end := bcst.String()
	return fmt.Sprintf("%s-%s", start, end) //返回用-表示的ip段,192.168.1.0-192.168.255.255
}

// parseIP1 解析ip段: 192.168.111.1-255 | 192.168.111.1-192.168.112.255(支持跨八位组)
func parseIP1(ip string) []string {
	IPRange := strings.Split(ip, "-")
	if len(IPRange) != 2 {
		gologger.Warningf("无法解析的IP段: %s", ip)
		return nil
	}
	var AllIP []string
	if len(IPRange[1]) < 4 {
		// 短格式,只允许写最后一段: 192.168.111.1-255
		SplitIP := strings.Split(IPRange[0], ".")
		if len(SplitIP) != 4 {
			gologger.Warningf("无法解析的IP段: %s", ip)
			return nil
		}
		endNum, err := strconv.Atoi(IPRange[1])
		if err != nil || endNum > 255 {
			gologger.Warningf("无法解析的IP段: %s", ip)
			return nil
		}
		startNum, err := strconv.Atoi(SplitIP[3])
		if err != nil || startNum > endNum {
			gologger.Warningf("无法解析的IP段: %s", ip)
			return nil
		}
		PrefixIP := strings.Join(SplitIP[0:3], ".")
		for i := startNum; i <= endNum; i++ {
			AllIP = append(AllIP, PrefixIP+"."+strconv.Itoa(i))
		}
	} else {
		// 完整格式,按数值比较,支持跨八位组区间
		startIP := net.ParseIP(IPRange[0])
		endIP := net.ParseIP(IPRange[1])
		if startIP == nil || endIP == nil {
			gologger.Warningf("无法解析的IP段: %s", ip)
			return nil
		}
		start4 := startIP.To4()
		end4 := endIP.To4()
		if start4 == nil || end4 == nil {
			// IPv6地址直接抛弃,仅扫描IPv4
			gologger.Warningf("忽略IPv6地址段: %s", ip)
			return nil
		}
		startNum := ipToUint32(start4)
		endNum := ipToUint32(end4)
		if startNum > endNum {
			gologger.Warningf("无法解析的IP段: %s", ip)
			return nil
		}
		for num := startNum; ; num++ {
			AllIP = append(AllIP, uint32ToIP(num))
			if num == endNum {
				break
			}
		}
	}
	return AllIP
}

func ipToUint32(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func uint32ToIP(num uint32) string {
	return strconv.Itoa(int(num>>24)) + "." + strconv.Itoa(int(num>>16&0xff)) + "." + strconv.Itoa(int(num>>8&0xff)) + "." + strconv.Itoa(int(num&0xff))
}
