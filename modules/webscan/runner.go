package webscan

import (
	"encoding/json"
	"errors"
	"github.com/imroc/req/v3"
	wappalyzer "github.com/projectdiscovery/wappalyzergo"
	"net/http"
	"puzzle/config"
	"puzzle/gologger"
	"puzzle/util"
	"strings"
	"sync"
	"time"
)

type Options struct {
	Url     []string
	Threads int
	Timeout int
	Proxy   string
}

type Runner struct {
	Fingerprint      []FingerPrint
	Options          *Options
	reqClient        *req.Client
	wappalyzerClient *wappalyzer.Wappalyze
}

func NewRunner(options *Options) *Runner {
	var fpSlice []FingerPrint
	content, err := util.ReadFile(util.GetRunDir() + config.FingerPrintPath)
	if err != nil {
		gologger.Fatalf("读取指纹文件失败:%s", err.Error())
	}
	if err := json.Unmarshal(content, &fpSlice); err != nil {
		gologger.Fatalf("json格式化失败:%s", err.Error())
	}

	wappalyzer, err := wappalyzer.New()
	if err != nil {
		gologger.Fatalf("wappalyzer初始化失败:%s", err.Error())
	}

	return &Runner{
		Fingerprint:      fpSlice,
		Options:          options,
		reqClient:        NewReqClient(options.Proxy, options.Timeout),
		wappalyzerClient: wappalyzer,
	}
}

func (r *Runner) Run() (results []*Result) {
	//确定请求类型
	r.Options.Url = r.CheckAlive()

	gologger.Infof("开始WEB扫描")
	wg := &sync.WaitGroup{}
	mutex := sync.Mutex{}
	taskChan := make(chan string, r.Options.Threads)
	for i := 0; i < r.Options.Threads; i++ {
		go func() {
			for task := range taskChan {
				res := r.GetFinger(task)
				// 判断蜜罐匹配大量指纹的情况
				if len(strings.Split(res.Finger, "\n")) > 5 {
					gologger.Infof("%v 可能为蜜罐", res.Url)
				} else {
					gologger.Silentf(FmtResult(res))
					mutex.Lock()
					results = append(results, res)
					mutex.Unlock()
				}
				wg.Done()
			}
		}()
	}

	for _, task := range r.Options.Url {
		wg.Add(1)
		taskChan <- task
	}
	close(taskChan)
	wg.Wait()

	gologger.Infof("WEB扫描结束")
	return
}

func NewReqClient(proxy string, timeout int) *req.Client {
	reqClient := req.C()
	reqClient.SetCommonRetryCount(3)
	reqClient.SetRedirectPolicy(checkRedirect("Cookie"))
	reqClient.GetTLSClientConfig().InsecureSkipVerify = true
	reqClient.SetCommonHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
		"Cookie":     "rememberMe=admin;rememberMe-K=admin", // check shiro
	})
	if proxy != "" {
		reqClient.SetProxyURL(proxy)
	}
	reqClient.SetTimeout(time.Duration(timeout) * time.Second)
	return reqClient
}

func checkRedirect(headers ...string) req.RedirectPolicy {
	return func(req *http.Request, via []*http.Request) error {
		//fmt.Println(via[0].Host)
		//自用，将url根据需求进行组合
		if len(via) >= 1 {
			if req.URL.Host != via[0].Host {
				return errors.New("stopped after 1 redirects")
			}
		}
		for _, header := range headers {
			if len(req.Header.Values(header)) > 0 {
				continue
			}
			vals := via[0].Header.Values(header)
			for _, val := range vals {
				req.Header.Add(header, val)
			}
		}
		return nil
	}
}
