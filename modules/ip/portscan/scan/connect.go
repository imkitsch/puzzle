package scan

import (
	"fmt"
	"net"
	"puzzle/gologger"
)

// ConnectVerify is used to verify if ports are accurate using a connect request
func (s *Scanner) ConnectVerify(host string, ports []int) []int {
	var verifiedPorts []int
	for _, port := range ports {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), s.timeout)
		if err != nil {
			continue
		}
		gologger.Debugf("Validated active port %d on %s\n", port, host)
		conn.Close()
		verifiedPorts = append(verifiedPorts, port)
	}
	return verifiedPorts
}
