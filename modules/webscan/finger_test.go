package webscan

import (
	"crypto/md5"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	wappalyzer "github.com/projectdiscovery/wappalyzergo"
)

func newFingerTestRunner(t *testing.T, fp []FingerPrint) *Runner {
	w, err := wappalyzer.New()
	if err != nil {
		t.Fatalf("wappalyzer初始化失败: %v", err)
	}
	return &Runner{
		Fingerprint:      fp,
		Options:          &Options{Threads: 5, Timeout: 5},
		reqClient:        NewReqClient("", 5),
		wappalyzerClient: w,
		dialer:           &net.Dialer{Timeout: 4 * time.Second},
	}
}

// 验证js跳转跟随:起始页只包含location.href跳转,指纹与标题应来自落地页
func TestGetFingerJsjump(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/start":
			fmt.Fprint(w, `<html><body><script type="text/javascript">location.href = "./landed";</script></body></html>`)
		case "/landed":
			fmt.Fprint(w, `<html><head><title>LandedPage</title></head><body>unique-jump-marker</body></html>`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	r := newFingerTestRunner(t, []FingerPrint{{Name: "jump-app", Keyword: []string{"unique-jump-marker"}}})
	res := r.GetFinger(srv.URL + "/start")
	if res.Title != "LandedPage" {
		t.Errorf("jsjump未生效, Title=%q", res.Title)
	}
	if res.Finger != "jump-app" {
		t.Errorf("跳转后页面指纹不符: %q", res.Finger)
	}
}

// 验证favicon指纹:页面声明icon链接,hash应为favicon内容的md5 hex并命中指纹库规则
func TestGetFingerFaviconHash(t *testing.T) {
	fav := []byte("fake-favicon-content-bytes")
	wantHash := fmt.Sprintf("%x", md5.Sum(fav))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			fmt.Fprint(w, `<html><head><link rel="icon" href="/favicon.ico"><title>FavPage</title></head></html>`)
		case "/favicon.ico":
			_, _ = w.Write(fav)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	r := newFingerTestRunner(t, []FingerPrint{{Name: "fav-app", FaviconHash: []string{wantHash}}})
	res := r.GetFinger(srv.URL + "/")
	if res.Finger != "fav-app" {
		t.Errorf("favicon指纹未命中: Finger=%q", res.Finger)
	}
	if res.Title != "FavPage" {
		t.Errorf("Title不符: %q", res.Title)
	}
}

// 验证页面未声明icon链接时回退到站点默认/favicon.ico
func TestGetFingerFaviconDefaultPath(t *testing.T) {
	fav := []byte("default-path-favicon")
	wantHash := fmt.Sprintf("%x", md5.Sum(fav))

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/":
			fmt.Fprint(w, `<html><head><title>DefaultFavPage</title></head></html>`)
		case "/favicon.ico":
			_, _ = w.Write(fav)
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	r := newFingerTestRunner(t, []FingerPrint{{Name: "default-fav-app", FaviconHash: []string{wantHash}}})
	res := r.GetFinger(srv.URL + "/")
	if res.Finger != "default-fav-app" {
		t.Errorf("默认路径favicon指纹未命中: Finger=%q", res.Finger)
	}
}
