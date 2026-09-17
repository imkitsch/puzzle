package qqwry

import "strings"

// cloudKeywords 云厂商关键字,命中则不参与C段扩展
var cloudKeywords = []string{
	"阿里云", "腾讯云", "华为云", "京东云", "UCloud", "火山引擎",
	"百度云", "青云", "金山云", "亚马逊", "AWS", "Azure", "微软云", "谷歌云", "Oracle",
}

func ExcludeCloud(ips []string) (cloudIp []string, normalIp []string) {
	QQwry := GetQqwry()
	for _, ip := range ips {
		info := QQwry.Find(ip)
		cloud := false
		// 包含匹配,兼容"阿里云计算有限公司"等完整归属地描述
		for _, name := range cloudKeywords {
			if strings.Contains(info.Area, name) || strings.Contains(info.Country, name) {
				cloud = true
				break
			}
		}
		if cloud {
			cloudIp = append(cloudIp, ip)
		} else {
			normalIp = append(normalIp, ip)
		}
	}
	return
}
