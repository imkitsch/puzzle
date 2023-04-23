package webscan

import (
	"github.com/imroc/req/v3"
	"net"
	"net/url"
	"path/filepath"
	"puzzle/gologger"
	"puzzle/util"
	"regexp"
	"strings"
)

var (
	/*
		<head>
			<meta http-equiv="refresh" content="1;URL='/admin'"/>
		</head>

		<!--
		<meta http-equiv="refresh" content="0.1;url=https://www.xxx.cn/">
		-->

		<head>
		<meta http-equiv=refresh content=0;url=index.jsp>
		</head>
	*/
	reg1 = regexp.MustCompile(`(?i)<meta.*?http-equiv=.*?refresh.*?url=(.*?)/?>`)
	/*
		<script type="text/javascript">
			location.href = "./ui/";
		</script>
	*/
	reg2 = regexp.MustCompile(`(?i)[window\.]?location[\.href]?.*?=.*?["'](.*?)["']`)
	/*
		<script language="javascript">
			window.location.replace("/mymeetings/");
		</script>
	*/
	reg3 = regexp.MustCompile(`(?i)window\.location\.replace\(['"](.*?)['"]\)`)
)

var (
	regHost = regexp.MustCompile(`(?i)https?://(.*?)/`)
)

func Jsjump(resp *req.Response) (jumpurl string) {
	res := regexJsjump(resp)
	if res != "" {
		gologger.Debugf("regexJsjump(), %v", res)
		res = strings.TrimSpace(res)
		res = strings.ReplaceAll(res, "\"", "")
		res = strings.ReplaceAll(res, "'", "")
		res = strings.ReplaceAll(res, "./", "/")
		if strings.HasPrefix(res, "http") {
			matches := regHost.FindAllStringSubmatch(res, -1)
			if len(matches) > 0 {
				var ip net.IP
				if strings.Contains(matches[0][1], ":") {
					ip = net.ParseIP(strings.Split(matches[0][1], ":")[0])
				} else {
					ip = net.ParseIP(matches[0][1])
				}
				if util.HasLocalIP(ip) {
					baseUrl := resp.Request.URL.Host
					res = strings.ReplaceAll(res, matches[0][1], baseUrl)
				}
			}
			jumpurl = res
		} else if strings.HasPrefix(res, "/") {
			// 前缀存在 / 时拼接绝对目录
			baseUrl := resp.Request.URL.Scheme + "://" + resp.Request.URL.Host
			jumpurl = baseUrl + res
		} else {
			// 前缀不存在 / 时拼接相对目录
			baseUrl := resp.Request.URL.Scheme + "://" + resp.Request.URL.Host + "/" + filepath.Dir(resp.Request.URL.Path) + "/"
			jumpurl = baseUrl + res
		}

		//判断是否跳转至同一网站
		u, err := url.Parse(jumpurl)
		if err != nil || u.Host != resp.Request.RawRequest.Host {
			jumpurl = ""
		}
	}
	return
}

func regexJsjump(resp *req.Response) (match string) {
	matches := reg1.FindAllStringSubmatch(resp.String(), -1)
	if len(matches) > 0 {
		// 去除注释的情况
		if !strings.Contains(resp.String(), "<!--\r\n"+matches[0][0]) && !strings.Contains(matches[0][1], "nojavascript.html") && !strings.Contains(resp.String(), "<!--[if lt IE 7]>\n"+matches[0][0]) {
			match = matches[0][1]
		}
	}
	body := resp.String()
	if len(body) > 700 {
		body = body[:700]
	}
	matches = reg2.FindAllStringSubmatch(body, -1)
	if len(matches) > 0 {
		match = matches[0][1]
	}
	matches = reg3.FindAllStringSubmatch(body, -1)
	if len(matches) > 0 {
		match = matches[0][1]
	}
	return ""
}
