package spider

import (
	"fmt"
	"testing"
	"time"

	"github.com/imroc/req/v3"
)

func TestQuakeEngineQueryCert(t *testing.T) {
	t.Skip("需要有效quake key并访问外部API,手动执行")
	keys := GetQuakeKey()
	if len(keys) == 0 {
		t.Fatal("未配置quake key")
	}
	client := req.C()
	client.SetTimeout(20 * time.Second)
	engine := &QuakeEngine{client: client, keys: keys}
	res := engine.QueryCert([]string{"example.com"})
	fmt.Println("结果数:", len(res))
	for _, v := range res[:min(10, len(res))] {
		fmt.Println("host:" + v.host, "ip:"+v.ip, "port:"+v.port)
	}
}
