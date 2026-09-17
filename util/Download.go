package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"puzzle/gologger"
	"strings"
	"sync"
	"time"

	"github.com/taythebot/cdn-ranges/provider"
	"golang.org/x/sync/errgroup"
)

// downloadClient 带超时的下载客户端,避免网络挂起时永久阻塞
var downloadClient = &http.Client{Timeout: 5 * time.Minute}

type GithubBody struct {
	Assets []Assets `json:"assets"`
}

type Assets struct {
	BrowserDownloadURL string `json:"browser_download_url"`
}

func GetGithubLatestUrl(repo string) string {
	url := "https://api.github.com/repos/" + repo + "/releases/latest"
	r, err := downloadClient.Get(url)
	if err != nil {
		gologger.Debugf("从GithubApi获取url失败: %s", err.Error())
		return ""
	}
	defer r.Body.Close()
	bodyDate, _ := io.ReadAll(r.Body)
	var jsonData GithubBody
	err = json.Unmarshal(bodyDate, &jsonData)
	if err != nil {
		gologger.Debugf("解析GithubRepos失败: %s", err.Error())
		return ""
	}
	if len(jsonData.Assets) > 0 {
		return jsonData.Assets[0].BrowserDownloadURL
	} else {
		return ""
	}
}

func Download(url string, filepath string) {
	if err := download(url, filepath); err != nil {
		gologger.Fatalf("下载失败:%s", err.Error())
	}
}

// download 下载url内容并原子替换到filepath:
// 校验HTTP状态码防止把错误页面写成数据文件,先写临时文件再重命名,中途失败不破坏旧文件
func download(url string, filepath string) error {
	resp, err := downloadClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &httpStatusError{code: resp.StatusCode, url: url}
	}

	tmpFile := filepath + ".tmp"
	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		out.Close()
		os.Remove(tmpFile)
		return err
	}
	if err := out.Close(); err != nil {
		os.Remove(tmpFile)
		return err
	}
	if err := os.Rename(tmpFile, filepath); err != nil {
		os.Remove(tmpFile)
		return err
	}
	return nil
}

type httpStatusError struct {
	code int
	url  string
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("HTTP状态码异常 %d %s: %s", e.code, http.StatusText(e.code), e.url)
}

func CdnRangeGenerate(filepath string) {
	providers := []provider.Provider{
		&provider.Akamai{},
		&provider.Bunny{},
		&provider.CacheFly{},
		&provider.CDN77{},
		&provider.CDNetworks{},
		&provider.Cloudflare{},
		&provider.Cloudfront{},
		&provider.DDoSGuard{},
		&provider.Edgecast{},
		&provider.EdgeNext{},
		&provider.Edgio{},
		&provider.F5{},
		&provider.Fastly{},
		&provider.GCore{},
		&provider.Imperva{},
		&provider.Leaseweb{},
		&provider.Limelight{},
		&provider.Medianova{},
		&provider.Qrator{},
		&provider.StormWall{},
		&provider.Sucuri{},
		&provider.X4B{},
	}

	// 先写临时文件,全部抓取结束后原子替换,失败不破坏旧文件
	tmpFile := filepath + ".tmp"
	outputFile, err := os.Create(tmpFile)
	if err != nil {
		gologger.Fatalf("创建临时文件出错:%s", err.Error())
	}

	var (
		count  int // 写入协程唯一,无需原子操作
		output = make(chan string, 1)
		queue  = make(chan provider.Provider, 4)
		g, ctx = errgroup.WithContext(context.Background())
		wg     sync.WaitGroup
		errMu  sync.Mutex
		failed []string
	)

	// 结果写入协程
	wg.Add(1)
	go func() {
		defer wg.Done()
		for o := range output {
			outputFile.WriteString(o + "\n")
			count++
		}
	}()

	// 任务队列
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, p := range providers {
			queue <- p
		}
		close(queue)
	}()

	// 抓取worker:单个源失败告警跳过,不影响其他源
	for i := 0; i < 4; i++ {
		g.Go(func() error {
			for p := range queue {
				gologger.Infof("Fetching %s ranges", p.Name())

				v4, v6, err := p.Fetch(ctx)
				if err != nil {
					gologger.Warningf("获取 %s ranges 失败,已跳过: %s", p.Name(), err.Error())
					errMu.Lock()
					failed = append(failed, p.Name())
					errMu.Unlock()
					continue
				}

				for _, r := range v4 {
					output <- r
				}
				for _, r := range v6 {
					output <- r
				}
			}
			return nil
		})
	}

	g.Wait()
	close(output)
	wg.Wait()

	if count == 0 {
		outputFile.Close()
		os.Remove(tmpFile)
		gologger.Fatalf("所有CDN数据源均获取失败,保留原有数据文件")
	}

	if err := outputFile.Close(); err != nil {
		os.Remove(tmpFile)
		gologger.Fatalf("关闭临时文件出错:%s", err.Error())
	}
	if err := os.Rename(tmpFile, filepath); err != nil {
		os.Remove(tmpFile)
		gologger.Fatalf("替换CDN数据文件出错:%s", err.Error())
	}

	gologger.Infof("CDN数据文件更新完成,共 %d 条", count)
	if len(failed) != 0 {
		gologger.Warningf("以下数据源获取失败已跳过: %s", strings.Join(failed, ", "))
	}
}
