package spider

import (
	"crypto/tls"
	"net"
)

// source:https://github.com/glebarez/cero
func grabCert(addr string, dialer *net.Dialer, onlyValidDomainNames bool) ([]string, error) {
	// dial
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// get first certificate in chain
	cert := conn.ConnectionState().PeerCertificates[0]

	// get CommonName and all SANs into a slice
	names := make([]string, 0, len(cert.DNSNames)+1)
	if onlyValidDomainNames && isDomainName(cert.Subject.CommonName) || !onlyValidDomainNames {
		names = append(names, cert.Subject.CommonName)
	}

	// append all SANs, excluding one that is equal to CN (if any)
	for _, name := range cert.DNSNames {
		if name != cert.Subject.CommonName {
			if onlyValidDomainNames && isDomainName(name) || !onlyValidDomainNames {
				names = append(names, name)
			}
		}
	}

	return names, nil
}

func isDomainName(s string) bool {
	// See RFC 1035, RFC 3696.
	// Presentation format has dots before every label except the first, and the
	// terminal empty label is optional here because we assume fully-qualified
	// (absolute) input. We must therefore reserve space for the first and last
	// labels' length octets in wire format, where they are necessary and the
	// maximum total length is 255.
	// So our _effective_ maximum is 253, but 254 is not rejected if the last
	// character is a dot.
	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}

	last := byte('.')
	nonNumeric := false // true once we've seen a letter or hyphen
	partlen := 0
	parts := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
			nonNumeric = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
			nonNumeric = true
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
			parts++
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return nonNumeric && (parts > 1 || parts > 0 && last != '.')
}
