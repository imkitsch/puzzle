package ip

import (
	"github.com/kayon/iploc"
	"puzzle/gologger"
	"sort"
	"strconv"
	"strings"
)

func GetIpInfo(ip string) (string, string) {
	loc, err := iploc.Open("config/qqwry_utf8.dat")
	if err != nil {
		gologger.Fatalf(err.Error())
	}

	detail := loc.Find(ip)
	result := strings.Fields(detail.String())
	Country := result[0]
	City := result[0]
	return Country, City
}

func GetSerialIp(ips []string) []string {
	var IpMap map[string][]string
	IpMap = make(map[string][]string)

	var cRangeElements []string
	var ipLastPartList []string

	var newIpList []string

	for _, ip := range ips {
		ipParts := strings.Split(ip, ".")
		cRange := ipParts[0] + "." + ipParts[1] + "." + ipParts[2]
		_, keyIs := IpMap[cRange]
		if keyIs == false {
			cRangeElements = []string{}
			cRangeElements = append(cRangeElements, ip)
			IpMap[cRange] = append(IpMap[cRange], ip)
		} else {
			IpMap[cRange] = append(IpMap[cRange], ip)
		}
	}

	//fmt.Println(IpMap)

	for _, valueList := range IpMap {
		if len(valueList) == 1 {
			newIpList = append(newIpList, valueList[0])
		} else {
			ipParts := strings.Split(valueList[0], ".")
			cRange := ipParts[0] + "." + ipParts[1] + "." + ipParts[2]
			ipLastPartList = []string{}
			for _, ip := range valueList {
				ipLastPartList = append(ipLastPartList, strings.Split(ip, ".")[3])
			}
			sort.Strings(ipLastPartList)
			start, _ := strconv.Atoi(ipLastPartList[0])
			end, _ := strconv.Atoi(ipLastPartList[len(ipLastPartList)-1])
			for i := start; i <= end; i++ {
				newIp := cRange + "." + strconv.Itoa(i)
				newIpList = append(newIpList, newIp)
			}
		}
	}
	return newIpList
}
