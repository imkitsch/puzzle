package vulscan

import (
	"bytes"
	"puzzle/gologger"
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
	gologger.Infof("开始漏洞扫描,模块为xpoc")

	urls := ""
	for _, url := range r.options.UrlList {
		urls = urls + url + "\n"
	}
	urls = urls[:len(urls)-1]

	err := xpocScan(bytes.NewBufferString(urls), r.options.Output+".html")
	if err != nil {
		gologger.Errorf("漏洞扫描报错:%s", err.Error())
	} else {
		gologger.Infof("漏洞扫描完成")
	}
}
