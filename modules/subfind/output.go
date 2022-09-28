package subfind

import (
	"github.com/boy-hack/ksubdomain/runner/result"
)

type ResultOutput struct {
	dr []domainResult
}

func NewDomainResult() (*ResultOutput, error) {
	r := &ResultOutput{}
	r.dr = []domainResult{}

	return r, nil
}

func (r *ResultOutput) WriteDomainResult(domain result.Result) error {
	var domainRes domainResult
	domainRes.domain = domain.Subdomain
	for _, item := range domain.Answers {
		if item[0:5] == "CNAME" {
			domainRes.cname = item[6:]
		} else {
			domainRes.a = append(domainRes.a, item)
		}
	}
	domainRes.cdn = IsCdn(domainRes.cname, domainRes.a)
	r.dr = append(r.dr, domainRes)
	return nil
}
func (r *ResultOutput) Close() {
	r.dr = nil
}

func (r *ResultOutput) OutPut() []domainResult {
	return r.dr
}
