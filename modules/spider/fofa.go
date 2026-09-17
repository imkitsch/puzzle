package spider

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"path/filepath"
	"puzzle/gologger"
	"puzzle/util"
	"strings"
	"time"

	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
)

// queryInterval 相邻查询的间隔,控制对测绘引擎的请求频率
const queryInterval = 2 * time.Second

type fofaAccount struct {
	Email string
	Key   string
}

// FofaEngine fofa测绘引擎,支持证书与IP检索
type FofaEngine struct {
	client   *req.Client
	accounts []fofaAccount
}

func (e *FofaEngine) Name() string {
	return "fofa"
}

// QueryCert 证书检索:cert="域" 以10个域为单位||打包,
// 单页查询(限定size保证结果在时效窗口内,且避免大数据量拖慢处理)
func (e *FofaEngine) QueryCert(domains []string) []*searchResult {
	var raw []string
	var tmp []string
	for k, domain := range domains {
		tmp = append(tmp, "cert=\""+domain+"\"")
		if (k+1)%10 == 0 || k == len(domains)-1 {
			raw = append(raw, strings.Join(tmp, "||"))
			tmp = []string{}
		}
	}

	var res []*searchResult
	for _, v := range raw {
		res = append(res, e.query(base64.StdEncoding.EncodeToString([]byte(v)))...)
		time.Sleep(queryInterval)
	}
	return res
}

// QueryIP IP检索:ip="x" 以10个为单位||打包,单页查询
func (e *FofaEngine) QueryIP(ips []string) []*searchResult {
	var raw []string
	var tmp []string
	for k, ip := range ips {
		tmp = append(tmp, "ip=\""+ip+"\"")
		if (k+1)%10 == 0 || k == len(ips)-1 {
			raw = append(raw, strings.Join(tmp, "||"))
			tmp = []string{}
		}
	}

	var res []*searchResult
	for _, v := range raw {
		res = append(res, e.query(base64.StdEncoding.EncodeToString([]byte(v)))...)
		time.Sleep(queryInterval)
	}
	return res
}

func (e *FofaEngine) query(qbase string) []*searchResult {
	if len(e.accounts) == 0 {
		return nil
	}
	// 多账号随机轮换
	account := e.accounts[rand.IntN(len(e.accounts))]
	resp, err := e.client.R().Get(fmt.Sprintf("https://fofa.info/api/v1/search/all?fields=host,ip,port,domain,protocol&email=%s&size=1000&key=%s&qbase64=%s", account.Email, account.Key, qbase))
	if err != nil {
		gologger.Errorf("fofa查询请求失败: %s", err.Error())
		return nil
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		gologger.Errorf("读取fofa响应失败: %s", err.Error())
		return nil
	}
	var apiResults ApiResults
	if err := json.Unmarshal(respBody, &apiResults); err != nil {
		gologger.Errorf("解析fofa响应失败: %s", err.Error())
		return nil
	}
	if apiResults.Error == true {
		gologger.Errorf("fofa查询失败: %s", apiResults.ErrMsg)
		return nil
	}
	var res []*searchResult
	for _, result := range apiResults.Results {
		if len(result) < 5 {
			continue
		}
		if len(result[0]) >= 5 {
			if result[0][:5] == "https" {
				result[0] = result[0][8:]
			} else if result[0][:4] == "http" {
				result[0] = result[0][7:]
			}
		}

		res = append(res, &searchResult{
			host:     result[0],
			ip:       result[1],
			port:     result[2],
			domain:   result[3],
			protocol: result[4],
		})
	}
	return res
}

// GetFofaKeys 解析provider-config中全部fofa账号,多账号查询时随机轮换
func GetFofaKeys() []fofaAccount {
	path := "/config/subfinder"
	v := viper.New()
	// 设置配置文件的名字
	v.SetConfigName("provider-config")
	// 设置配置文件的类型
	v.SetConfigType("yaml")
	// 添加配置文件的路径，指定目录下寻找
	v.AddConfigPath(filepath.Join(util.GetRunDir() + path))
	// 寻找配置文件并读取
	err := v.ReadInConfig()
	if err != nil {
		return nil
	}
	dataRaw := v.GetStringSlice("fofa")
	var accounts []fofaAccount
	for _, item := range dataRaw {
		data := strings.Split(item, ":")
		if len(data) != 2 {
			gologger.Errorf("fofa key格式错误,应为email:key: %s", item)
			continue
		}
		accounts = append(accounts, fofaAccount{Email: data[0], Key: data[1]})
	}
	return accounts
}

// GetFofaKey 返回第一个可用fofa账号,保持原接口签名供core调用
func GetFofaKey() (email string, key string) {
	accounts := GetFofaKeys()
	if len(accounts) == 0 {
		return "", ""
	}
	return accounts[0].Email, accounts[0].Key
}
