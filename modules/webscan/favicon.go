package webscan

import (
	"bytes"
	"github.com/antchfx/htmlquery"
	"github.com/imroc/req/v3"
	"golang.org/x/net/html"
	"strings"
)

func (r *Runner) getFavicon(resp *req.Response) (favicon string, iconHash string) {
	htmlDoc, err := html.Parse(bytes.NewReader(resp.Bytes()))
	if err != nil {
		return favicon, iconHash
	}
	if nodes, err := htmlquery.QueryAll(htmlDoc, "//link"); err != nil {
		return favicon, iconHash
	} else {
		for _, node := range nodes {
			if htmlquery.SelectAttr(node, "href") != "" && strings.Contains(htmlquery.SelectAttr(node, "rel"), "icon") {
				favicon = htmlquery.SelectAttr(node, "href")
				break
			}
		}
	}
	if favicon == "" {
		// 页面未声明icon链接时,尝试站点默认favicon路径
		favicon = resp.Request.URL.Scheme + "://" + resp.Request.URL.Host + "/favicon.ico"
	}
	if !strings.Contains(favicon, "http") {
		favicon = strings.TrimSpace(favicon)
		favicon = strings.TrimLeft(favicon, "../") // nolint
		favicon = strings.TrimLeft(favicon, "./")
		if !strings.HasPrefix(favicon, "/") {
			favicon = "/" + favicon
		}
		favicon = resp.Request.URL.Scheme + "://" + resp.Request.URL.Host + favicon
	}
	iconHash = r.GetHash(favicon)
	if iconHash == "" {
		favicon = ""
	}
	return favicon, iconHash
}
