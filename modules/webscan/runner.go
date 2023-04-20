package webscan

import (
	"encoding/json"
	"github.com/imroc/req/v3"
	"puzzle/gologger"
	"puzzle/util"
	"time"
)

var FingerPrintPath = "/config/web/web_fingerprint_v3.json"

type Options struct {
	Url     []string
	Threads int
	Timeout int
	Proxy   string
}

type Runner struct {
	Fingerprint []FingerPrint
	Options     *Options
	reqClient   *req.Client
}

func NewRunner(options *Options) *Runner {
	var fpSlice []FingerPrint
	content, err := util.ReadFile(util.GetRunDir() + FingerPrintPath)
	if err != nil {
		gologger.Fatalf("读取指纹文件失败:%s", err.Error())
	}
	if err := json.Unmarshal(content, &fpSlice); err != nil {
		gologger.Fatalf("json格式化失败:%s", err.Error())
	}

	return &Runner{
		Fingerprint: fpSlice,
		Options:     options,
		reqClient:   NewReqClient(options.Proxy, options.Timeout),
	}
}

func (r *Runner) Run() {

}

func NewReqClient(proxy string, timeout int) *req.Client {
	reqClient := req.C()
	reqClient.GetTLSClientConfig().InsecureSkipVerify = true
	reqClient.SetCommonHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
		"Cookie":     "rememberMe=admin;rememberMe-K=admin", // check shiro
	})
	reqClient.SetRedirectPolicy(req.AlwaysCopyHeaderRedirectPolicy("Cookie"))
	if proxy != "" {
		reqClient.SetProxyURL(proxy)
	}
	reqClient.SetTimeout(time.Duration(timeout) * time.Second)
	return reqClient
}
