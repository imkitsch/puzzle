package util

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("hello"))
	}))
	defer srv.Close()

	target := filepath.Join(t.TempDir(), "data.bin")
	if err := download(srv.URL, target); err != nil {
		t.Fatalf("下载应成功: %v", err)
	}
	got, err := os.ReadFile(target)
	if err != nil || string(got) != "hello" {
		t.Errorf("文件内容不符: %q, err=%v", got, err)
	}
	// 临时文件不应残留
	if _, err := os.Stat(target + ".tmp"); !os.IsNotExist(err) {
		t.Error("临时文件应已被重命名清除")
	}
}

func TestDownloadBadStatusKeepsOldFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	dir := t.TempDir()
	target := filepath.Join(dir, "qqwry.dat")
	if err := os.WriteFile(target, []byte("old"), 0600); err != nil {
		t.Fatal(err)
	}

	err := download(srv.URL+"/missing", target)
	if err == nil {
		t.Fatal("404应返回错误")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Errorf("错误信息应包含状态码: %v", err)
	}
	// 旧文件必须保留,不能被错误页面覆盖
	got, rerr := os.ReadFile(target)
	if rerr != nil || string(got) != "old" {
		t.Errorf("旧文件应保持不变: %q, err=%v", got, rerr)
	}
	if _, err := os.Stat(target + ".tmp"); !os.IsNotExist(err) {
		t.Error("失败后不应残留临时文件")
	}
}
