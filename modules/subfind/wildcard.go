package subfind

import (
	"errors"
	"net"
	"puzzle/gologger"
	"puzzle/util"
)

func IsWildCard(domain string) bool {
	for i := 0; i < 2; i++ {
		subdomain := util.RandomStr(6) + "." + domain
		_, err := net.LookupIP(subdomain)
		if err == nil {
			return true
		}
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) && dnsErr.IsNotFound {
			// 明确的NXDOMAIN,该随机子域名不存在,继续验证
			continue
		}
		// 超时或DNS服务器故障等无法判断的情况,保守按存在泛解析处理,
		// 跳过爆破,避免对泛解析域名产生大量假阳性子域名
		gologger.Warningf("域名 %s 泛解析判断受DNS故障影响,默认跳过爆破", domain)
		return true
	}
	return false
}
