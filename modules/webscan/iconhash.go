package webscan

import (
	"crypto/md5"
	"fmt"
)

// GetHash 取favicon内容并计算指纹库使用的md5 hex。
// 指纹库(FingerprintHub)的favicon_hash为favicon文件的md5 hex,
// 与fofa平台的mmh3 icon hash是两套不同的约定
func (r *Runner) GetHash(url string) (iconhash string) {
	content, err := r.FromURLGetContent(url)
	if err != nil {
		return
	}
	if len(content) > 0 {
		iconhash = fmt.Sprintf("%x", md5.Sum(content))
	}
	return
}

func (r *Runner) FromURLGetContent(url string) (content []byte, err error) {
	resp, err := r.reqClient.R().Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Req %v status is not 200", url)
	}
	content = resp.Bytes()
	return
}
