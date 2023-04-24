package util

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"puzzle/gologger"
)

type GithubBody struct {
	Assets []Assets `json:"assets"`
}

type Assets struct {
	BrowserDownloadURL string `json:"browser_download_url"`
}

func GetGithubLatestUrl(repo string) string {
	url := "https://api.github.com/repos/" + repo + "/releases/latest"
	r, err := http.Get(url)
	if err != nil {
		gologger.Debugf("从GithubApi获取url失败: %s", err.Error())
		return ""
	}
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
	resp, err := http.Get(url)
	if err != nil {
		gologger.Fatalf("下载文件连接报错:%s", err.Error())
	}
	defer resp.Body.Close()

	if FileExists(filepath) {
		err := os.Remove(filepath)
		if err != nil {
			gologger.Fatalf("删除旧文件出错:%s", err.Error())
		}
	}
	out, err := os.Create(filepath)
	if err != nil {
		gologger.Fatalf("创建新文件出错:%s", err.Error())
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		gologger.Fatalf("连接句柄出错%s", err.Error())
	}
}
