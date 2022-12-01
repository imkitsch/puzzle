package webscan

import "github.com/imroc/req/v3"

func (r *Runner) GetFinger(resp *req.Response) string {
	var fingerResult []string
	bodyString := resp.String()
	header := getHeaderString(resp)
	content := resp.Bytes()
	iconhash := Mmh3Hash32(StandBase64(content))

	return nil
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
