package webscan

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"puzzle/util"
	"strings"
)

func FmtResult(result *Result) (res string) {
	finger := strings.Replace(result.Finger, "\n", ",", -1)
	wappalyzer := strings.Replace(result.Wappalyzer, "\n", ",", -1)
	if util.IsWindows() {
		res = fmt.Sprintf("%v [%v] [%v] [%v] [%v] [%v]\n", result.Url, result.StatusCode, result.Length, result.Title, finger, wappalyzer)
	} else {
		res = fmt.Sprintf("%v [%v] [%v] [%v] [%v] [%v]\n", result.Url, aurora.Red(result.StatusCode), aurora.Yellow(result.Length), aurora.Green(result.Title), aurora.Blue(finger), aurora.Cyan(wappalyzer))
	}
	return
}
