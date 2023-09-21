package spider

import "github.com/imroc/req/v3"

type Options struct {
	Domains    []string
	Subdomains []string
	Ips        []string
	Key        string
	Email      string
}

type Runner struct {
	options   *Options
	reqClient *req.Client
}

type fofaResult struct {
	host     string
	ip       string
	port     string
	domain   string
	protocol string
}

type Result struct {
	Urls          []string
	AddDomains    []*DResult
	AddSubdomains []*SdResult
}

type SdResult struct {
	Domain string
	CDN    bool
	Ip     string
	CName  string
}

type DResult struct {
	Domain string
}

type ApiResults struct {
	Mode    string     `json:"mode"`
	Error   bool       `json:"error"`
	ErrMsg  string     `json:"errmsg"`
	Query   string     `json:"query"`
	Page    int        `json:"page"`
	Size    int        `json:"size"`
	Results [][]string `json:"results"`
}
