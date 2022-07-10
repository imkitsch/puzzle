package ip

import (
	"Allin/gologger"
	"github.com/kayon/iploc"
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
