package util

import (
	"io"
	"net/http"
	"os"
	"puzzle/gologger"
)

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
