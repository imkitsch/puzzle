package vulscan

import (
	"bytes"
	"io"
	"os/exec"
	"puzzle/config"
	"puzzle/gologger"
	"puzzle/util"
)

func xpocScan(urls io.Reader, output string) error {
	var stdout bytes.Buffer
	cmd := exec.Command(util.GetRunDir()+config.XpocPath, "-o", output)
	cmd.Stdin = urls
	cmd.Stdout = &stdout
	err := cmd.Run()
	if util.FileExists(output) {
		gologger.Infof("存在漏洞,请查看:%s", output)
	} else {
		gologger.Infof("未扫描到漏洞")
	}
	if err != nil {
		return err
	}
	return nil
}
