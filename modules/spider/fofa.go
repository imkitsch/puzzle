package spider

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"net"
	"path/filepath"
	"puzzle/util"
	"strings"
)

func (r *Runner) getFofaResult(qbase string) []*fofaResult {
	resp, err := r.reqClient.R().Get(fmt.Sprintf("https://fofa.info/api/v1/search/all?fields=host,ip,port,domain,protocol&email=%s&size=5000&key=%s&qbase64=%s", r.options.Email, r.options.Key, qbase))
	if err != nil {
		return nil
	}
	var res []*fofaResult
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	//fmt.Println(string(respBody))
	var apiResults ApiResults
	json.Unmarshal(respBody, &apiResults)
	if apiResults.Error == true {
		return nil
	}
	for _, result := range apiResults.Results {
		if result[0][:5] == "https" {
			result[0] = result[0][8:]
		} else if result[0][:4] == "http" {
			result[0] = result[0][7:]
		}

		if net.ParseIP(result[0]) != nil && strings.Contains(result[0], ":") == false {
			result[0] = result[0] + ":" + result[2]
		}

		res = append(res, &fofaResult{
			host:     result[0],
			ip:       result[1],
			port:     result[2],
			domain:   result[3],
			protocol: result[4],
		})
	}
	return res
}

func GetFofaKey() (email string, key string) {
	path := "/config/subfinder"
	// 设置配置文件的名字
	viper.SetConfigName("provider-config")
	// 设置配置文件的类型
	viper.SetConfigType("yaml")
	// 添加配置文件的路径，指定目录下寻找
	viper.AddConfigPath(filepath.Join(util.GetRunDir() + path))
	// 寻找配置文件并读取
	err := viper.ReadInConfig()
	dataRaw := viper.GetStringSlice("fofa")
	if err != nil || len(dataRaw) == 0 {
		return "", ""
	}
	data := strings.Split(dataRaw[0], ":")
	return data[0], data[1]
}
