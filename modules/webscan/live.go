package webscan

import (
	"fmt"
	"puzzle/gologger"
	"strings"
	"sync"
)

func (r *Runner) CheckAlive() (results []string) {
	var (
		ToHttps = []string{
			"sent to HTTPS port",
			"This combination of host and port requires TLS",
			"Instead use the HTTPS scheme to",
			"This web server is running in SSL mode",
		}
	)

	if len(r.Options.Url) == 0 {
		return
	}
	gologger.Infof("开始HTTP探活: %v", len(r.Options.Url))

	// RunTask
	wg := &sync.WaitGroup{}
	mutex := sync.Mutex{}
	taskChan := make(chan string, r.Options.Threads)
	for i := 0; i < r.Options.Threads; i++ {
		go func() {
			for task := range taskChan {
				resp, err := r.reqClient.R().Get("http://" + task)
				if err == nil {
					mutex.Lock()
					flag := false
					for _, str := range ToHttps {
						if strings.Contains(resp.String(), str) {
							flag = true
							break
						}
					}
					if flag == true {
						resp, err = r.reqClient.R().Get("https://" + task)
						if err == nil {
							results = append(results, resp.Response.Request.URL.String())
						} else {
							gologger.Warningf("%v", err)
						}
					} else {
						results = append(results, resp.Response.Request.URL.String())
					}
					mutex.Unlock()
				} else {
					resp, err = r.reqClient.R().Get("https://" + task)
					if err == nil {
						mutex.Lock()
						results = append(results, resp.Response.Request.URL.String())
						mutex.Unlock()
					} else {
						gologger.Warningf("%v", err)
					}
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

	gologger.Infof("HTTP探活结束")
	fmt.Println(results)
	return
}
