package util

import (
	"bufio"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"math/rand"
	"os"
	"puzzle/gologger"
	"runtime"
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
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
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

func GetWindowWith() int {
	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 0
	}
	return w
}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
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

func GetRunDir() string {
	dir, err := os.Getwd()
	if err != nil {
		gologger.Fatalf("Could not get Run directory: %s\n", err)
	}
	return dir
}

func ParsePorts(portInput string) []int {
	var portOutput []int
	portTemp := strings.Split(portInput, ",")
	for _, value := range portTemp {
		temp := strings.Split(value, "-")
		if len(temp) == 2 {
			min, _ := strconv.Atoi(temp[0])
			max, _ := strconv.Atoi(temp[1])
			for i := min; i <= max; i++ {
				portOutput = append(portOutput, i)
			}
		} else {
			p, _ := strconv.Atoi(value)
			portOutput = append(portOutput, p)
		}
	}
	return portOutput
}

func isOSSupported() bool {
	return IsLinux() || IsOSX()
}

func IsOSX() bool {
	return runtime.GOOS == "darwin"
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}
