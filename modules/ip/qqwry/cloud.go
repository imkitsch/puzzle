package qqwry

func ExcludeCloud(ips []string) (cloudIp []string, normalIp []string) {
	var cloudName map[string]string
	cloudName = make(map[string]string)
	cloudName["阿里云"] = ""
	cloudName["腾讯云"] = ""
	cloudName["华为云"] = ""
	cloudName["京东云"] = ""
	cloudName["UCloud"] = ""

	QQwry := GetQqwry()
	for _, ip := range ips {
		info := QQwry.Find(ip)
		_, ok := cloudName[info.Area]
		if ok {
			cloudIp = append(cloudIp, ip)
		} else {
			normalIp = append(normalIp, ip)
		}
	}
	return
}
