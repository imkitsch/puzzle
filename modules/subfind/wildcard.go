package subfind

import "net"
import "puzzle/util"

func IsWildCard(domain string) bool {
	ranges := [2]int{}
	for _, _ = range ranges {
		subdomain := util.RandomStr(6) + "." + domain
		_, err := net.LookupIP(subdomain)
		if err != nil {
			continue
		}
		return true
	}
	return false
}
