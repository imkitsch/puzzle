package qqwry

import "testing"

// 畸形IP不应触发panic,直接返回空结果
func TestFindMalformedIP(t *testing.T) {
	q := &QQwry{Data: &fileData{Data: make([]byte, 16)}}
	for _, ip := range []string{"1.2.3.", "1.2.3.4.5", "", "abc.def.ghi.jkl", "999.1.1.1"} {
		res := q.Find(ip)
		if res.IP != ip {
			t.Errorf("回显IP不符: %s -> %s", ip, res.IP)
		}
		if res.Country != "" || res.Area != "" {
			t.Errorf("畸形IP %s 不应有归属地结果: %+v", ip, res)
		}
	}
}
