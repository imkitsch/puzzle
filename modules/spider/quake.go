package spider

import (
	"path/filepath"
	"puzzle/gologger"
	"puzzle/util"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
)

// QuakeEngine quake测绘引擎。
// quake api可用量较少,且证书检索比fofa全面,仅用于证书检索,不提供IP查询
type QuakeEngine struct {
	client *req.Client
	keys   []string
}

func (e *QuakeEngine) Name() string {
	return "quake"
}

// QueryCert 证书检索:cert:"域" 以5个域为单位OR打包,500/页自动翻页,每页随机换key
func (e *QuakeEngine) QueryCert(domains []string) []*searchResult {
	var quakeRaw []string
	var tmp []string
	for k, domain := range domains {
		tmp = append(tmp, "cert:\""+domain+"\"")
		if (k+1)%5 == 0 || k == len(domains)-1 {
			quakeRaw = append(quakeRaw, strings.Join(tmp, " OR "))
			tmp = []string{}
		}
	}

	var res []*searchResult
	for _, v := range quakeRaw {
		res = append(res, e.query(v)...)
		time.Sleep(queryInterval)
	}
	return res
}

func (e *QuakeEngine) query(raw string) []*searchResult {
	var pageSize = 500
	var start = 0
	var totalResults = -1

	type QuakeRequest struct {
		Query     string   `json:"query"`
		Include   []string `json:"include"`
		Start     int      `json:"start"`
		Size      int      `json:"size"`
		Shortcuts []string `json:"shortcuts"`
	}

	var res []*searchResult

	for {
		reqBody := QuakeRequest{
			Query:   raw,
			Include: []string{"service.http.host", "ip", "port"},
			Start:   start,
			Size:    pageSize,
		}
		var quakeResp QuakeResponse
		key := util.RandomPick(e.keys)

		resp, err := e.client.R().
			SetHeader("X-QuakeToken", key). // 填入你的 Token
			SetBody(&reqBody).
			SetSuccessResult(&quakeResp). // 如果状态码为 2xx，自动解析到 quakeResult
			Post("https://quake.360.net/api/v3/search/quake_service")
		if err != nil {
			gologger.Errorf("quake查询请求失败: %s", err.Error())
			return res
		}
		if !resp.IsSuccessState() {
			gologger.Errorf("quake查询响应异常: %d", resp.StatusCode)
			return res
		}
		if quakeResp.Code != 0 {
			gologger.Errorf("quake查询失败: %s", quakeResp.Message)
			return res
		}

		if totalResults == -1 {
			totalResults = quakeResp.Meta.Pagination.Total
		}
		for _, v := range quakeResp.Data {
			var rootDomain string
			var host string
			rootDomain, err = GetRootDomain(v.Service.HTTP.Host)
			if err != nil || v.Service.HTTP.Host == v.Ip {
				rootDomain = ""
			}

			if v.Service.HTTP.Host == "" {
				host = v.Ip
			} else {
				host = v.Service.HTTP.Host
			}
			tmp := &searchResult{
				host:     host,
				ip:       v.Ip,
				port:     strconv.Itoa(v.Port),
				domain:   rootDomain,
				protocol: "tls",
			}
			res = append(res, tmp)
		}

		if len(quakeResp.Data) == 0 || start+pageSize >= totalResults {
			break
		}
		start += pageSize

	}

	return res
}

func GetQuakeKey() []string {
	path := "/config/cert"
	v := viper.New()
	// 设置配置文件的名字
	v.SetConfigName("config")
	// 设置配置文件的类型
	v.SetConfigType("yaml")
	// 添加配置文件的路径，指定目录下寻找
	v.AddConfigPath(filepath.Join(util.GetRunDir() + path))
	// 寻找配置文件并读取
	err := v.ReadInConfig()
	if err != nil {
		gologger.Debugf("读取quake配置失败: %v", err)
		return []string{}
	}
	dataRaw := v.GetStringSlice("quake")
	if len(dataRaw) == 0 {
		return []string{}
	}
	return dataRaw
}
