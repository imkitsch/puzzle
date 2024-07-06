package webscan

import (
	"bytes"
	"errors"
	"github.com/imroc/req/v3"
	"io"
	urltool "net/url"
	"puzzle/util"
	"strconv"
	"strings"
)

func (r *Runner) GetFinger(url string) *Result {
	respBody, headers, rawHeaders, statusCode, length, title, iconHash, err := getResponseData(r.reqClient, url)
	if err != nil {
		return &Result{
			Url: url,
		}
	}

	fingerprints := r.wappalyzerClient.Fingerprint(rawHeaders, respBody)
	var wappzerTmp []string
	var wappzerRes string
	for name, _ := range fingerprints {
		wappzerTmp = append(wappzerTmp, name)
	}
	if len(wappzerTmp) > 0 {
		wappzerRes = strings.Join(wappzerTmp, "\n")
	} else {
		wappzerRes = ""
	}

	var fpName []string
	for _, v := range r.Fingerprint {
		hflag := true
		if len(v.Headers) > 0 {
			hflag = false
			for k, h := range v.Headers {
				if headers[k] != "" {
					// fmt.Println("header key ", headers[k][0], " => h :", h)
					if !util.StringSearch(headers[k], h) {
						hflag = false
						break
					}
					hflag = true
				} else {
					hflag = false
					break
				}
			}
		}

		kflag := true
		if len(v.Keyword) > 0 {
			kflag = false
			for _, k := range v.Keyword {
				if !util.StringSearch(string(respBody), k) {
					kflag = false
					break
				}
				kflag = true
			}
		}

		iflag := true
		if len(v.FaviconHash) > 0 {
			iflag = false
			for _, k := range v.FaviconHash {
				if k == iconHash {
					iflag = true
					break
				}
				iflag = false
			}
		}

		if iflag && kflag && hflag {
			fpName = append(fpName, v.Name)
		}
	}
	fpName = util.RemoveRepeatedStringElement(fpName)
	fingerRes := strings.Join(fpName, "\n")

	// 获取证书
	var cert string
	if url[0:5] == "https" {
		host, err := urltool.Parse(url)
		if err != nil {

		}
		tmp := strings.Split(host.Host, ":")
		var ip string
		var port string

		ip = tmp[0]
		if len(tmp) == 1 {
			port = "443"
		} else {
			port = tmp[1]
		}

		certList, err := GrabCert(ip+":"+port, r.dialer, false)
		if len(certList) > 0 && err != nil {
			cert = certList[0]
		} else {
			cert = ""
		}

	} else {
		cert = ""
	}

	return &Result{
		Url:        url,
		StatusCode: strconv.Itoa(statusCode),
		Length:     length,
		Title:      title,
		Finger:     fingerRes,
		Wappalyzer: wappzerRes,
		Cert:       cert,
	}
}

func getHeaderString(resp *req.Response) map[string]string {
	headerMap := map[string]string{}
	for k := range resp.Header {
		if k != "Set-Cookie" {
			headerMap[k] = resp.Header.Get(k)
		}
	}
	for _, ck := range resp.Cookies() {
		headerMap["Set-Cookie"] += ck.String() + ";"
	}
	return headerMap
}

func getResponseData(client *req.Client, url string) ([]byte, map[string]string, map[string][]string, int, int, string, string, error) {
	if url == "" {
		return []byte(""), nil, nil, 0, 0, "", "", errors.New("no target specified")
	}
	request := client.R()
	resp, err := request.Get(url)

	// 处理js跳转, 上限3次
	for i := 0; i < 3; i++ {
		jumpUrl := Jsjump(resp)
		if jumpUrl == "" {
			break
		}
		resp, err = request.Get(jumpUrl)
	}

	if err != nil {
		return []byte(""), nil, nil, 0, 0, "", "", err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return []byte(""), nil, nil, 0, 0, "", "", err
	}

	reader := io.NopCloser(bytes.NewBuffer(respBody))
	title := getHTTPTitle(&reader)

	header := getHeaderString(resp)

	content := resp.Bytes()
	iconHash := Mmh3Hash32(StandBase64(content))

	return respBody, header, resp.Header, resp.StatusCode, len(resp.String()), title, iconHash, nil
}
