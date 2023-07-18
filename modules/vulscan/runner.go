package vulscan

import (
	"bytes"
	"fmt"
	"puzzle/gologger"
	"strings"
)

type Options struct {
	UrlList []string
	Output  string
}

type Runner struct {
	options *Options
}

func NewRunner(options *Options) *Runner {
	return &Runner{
		options: options,
	}
}

func (r *Runner) Run() {
	if len(r.options.UrlList) == 0 {
		gologger.Infof("无存活web,跳过漏洞扫描")
		return
	}
	gologger.Infof("开始漏洞扫描,模块为xpoc")

	urls := fmt.Sprintf(strings.Join(r.options.UrlList, "\n"))

	err := xpocScan(bytes.NewBufferString(urls), r.options.Output+".html")
	if err != nil {
		gologger.Errorf("漏洞扫描报错:%s", err.Error())
	} else {
		gologger.Infof("漏洞扫描完成")
	}
}
