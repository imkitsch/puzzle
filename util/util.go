package util

import (
	"bufio"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"puzzle/gologger"
	"sort"
	"strconv"
	"strings"
	"time"
)

func RandomStr(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ReadFile(filename string) (bytes []byte, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	return data, nil
}

func LinesInFile(fileName string) ([]string, error) {
	result := []string{}
	f, err := os.Open(fileName)
	if err != nil {
		return result, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			result = append(result, line)
		}
	}
	return result, nil
}

// FileExists 判断文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// RemoveRepeatedStringElement string切片内元素去重
func RemoveRepeatedStringElement(arr []string) []string {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}
	return arr[:j]
}

// RemoveRepeatedIntElement int切片元素去重
func RemoveRepeatedIntElement(s []int) []int {
	result := make([]int, 0)
	m := make(map[int]bool)
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

// CreateDir 调用递归创建文件夹
func CreateDir(filePath string) error {
	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

// DeleteDir 删除文件夹
func DeleteDir(filePath string) error {
	err := os.RemoveAll(filePath)
	if err != nil {
		return err
	}
	return nil
}

func GetRunDir() string {
	dir, err := os.Getwd()
	if err != nil {
		gologger.Fatalf("Could not get Run directory: %s\n", err)
	}
	return dir
}

// GetSerialIp 获取连续ip段
func GetSerialIp(ips []string) []string {
	var IpMap map[string][]string
	IpMap = make(map[string][]string)

	var cRangeElements []string
	var ipLastPartList []string

	var newIpList []string

	for _, ip := range ips {
		ipParts := strings.Split(ip, ".")
		cRange := ipParts[0] + "." + ipParts[1] + "." + ipParts[2]
		_, keyIs := IpMap[cRange]
		if keyIs == false {
			cRangeElements = []string{}
			cRangeElements = append(cRangeElements, ip)
			IpMap[cRange] = append(IpMap[cRange], ip)
		} else {
			IpMap[cRange] = append(IpMap[cRange], ip)
		}
	}

	for _, valueList := range IpMap {
		if len(valueList) == 1 {
			newIpList = append(newIpList, valueList[0])
		} else {
			ipParts := strings.Split(valueList[0], ".")
			cRange := ipParts[0] + "." + ipParts[1] + "." + ipParts[2]
			ipLastPartList = []string{}
			for _, ip := range valueList {
				ipLastPartList = append(ipLastPartList, strings.Split(ip, ".")[3])
			}
			sort.Strings(ipLastPartList)
			start, _ := strconv.Atoi(ipLastPartList[0])
			end, _ := strconv.Atoi(ipLastPartList[len(ipLastPartList)-1])
			for i := start; i <= end; i++ {
				newIp := cRange + "." + strconv.Itoa(i)
				newIpList = append(newIpList, newIp)
			}
		}
	}
	return newIpList
}

// HasLocalIP 检测 IP 地址是否是内网地址
func HasLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}

func StringSearch(s string, sub string) bool {
	s = strings.ToLower(s)
	sub = strings.ToLower(sub)
	return strings.Contains(s, sub)
}
