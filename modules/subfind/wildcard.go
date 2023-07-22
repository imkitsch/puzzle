package subfind

import (
	"net"
)
import "puzzle/util"

func IsWildCard(domain string) bool {
	for i := 0; i < 2; i++ {
		subdomain := util.RandomStr(6) + "." + domain
		_, err := net.LookupIP(subdomain)
		if err != nil {
			continue
		}
		return true
	}
	return false
}
