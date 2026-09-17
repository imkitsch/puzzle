package spider

import (
	"golang.org/x/net/publicsuffix"
)

func GetRootDomain(domain string) (string, error) {
	root, err := publicsuffix.EffectiveTLDPlusOne(domain)
	if err != nil {
		return "", err
	}
	return root, nil
}
