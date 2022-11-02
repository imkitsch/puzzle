package subfind

import (
	"github.com/boy-hack/ksubdomain/runner/result"
	"strings"
)

type ResultOutput struct {
	dr []*domainResult
}

func NewDomainResult() (*ResultOutput, error) {
	r := &ResultOutput{}
	r.dr = []*domainResult{}

	return r, nil
}

func (r *ResultOutput) WriteDomainResult(domain result.Result) error {
	var domainRes domainResult
	var Address []string
	domainRes.Domain = domain.Subdomain
	for _, item := range domain.Answers {
		if item[0:5] == "CNAME" {
			domainRes.Cname = item[6:]
		} else {
			Address = append(Address, item)
		}
	}
	domainRes.Cdn = IsCdn(domainRes.Cname, Address)
	domainRes.Address = strings.Join(Address, `,`)
	r.dr = append(r.dr, &domainRes)
	return nil
}
func (r *ResultOutput) Close() {
	r.dr = nil
}

func (r *ResultOutput) OutPut() []*domainResult {
	return r.dr
}
