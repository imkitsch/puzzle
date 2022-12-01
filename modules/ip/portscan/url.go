package portscan

import "fmt"

func GetUrl(result []*Result) (url []string) {
	for _, res := range result {
		if res.ServiceName == "http" || res.ServiceName == "https" {
			url = append(url, fmt.Sprintf("%s://%s:%d", res.ServiceName, res.Addr, res.Port))
		} else if res.ServiceName == "unknown" {
			url = append(url, fmt.Sprintf("%s:%d", res.Addr, res.Port))
		}
	}
	return url
}
