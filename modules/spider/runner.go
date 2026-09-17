package spider

import (
	"time"

	"github.com/imroc/req/v3"
)

func NewRunner(options *Options) *Runner {
	reqClient := req.C()
	reqClient.SetTimeout(20 * time.Second)

	r := &Runner{
		options:   options,
		reqClient: reqClient,
	}

	// fofa:证书与IP检索
	if accounts := GetFofaKeys(); len(accounts) != 0 {
		fofa := &FofaEngine{client: reqClient, accounts: accounts}
		r.certEngines = append(r.certEngines, fofa)
		r.ipEngines = append(r.ipEngines, fofa)
	}
	// quake:api可用量较少且证书检索比fofa全面,仅用于证书检索
	if keys := GetQuakeKey(); len(keys) != 0 {
		r.certEngines = append(r.certEngines, &QuakeEngine{client: reqClient, keys: keys})
	}
	return r
}
