package core

import (
	"puzzle/config"
	"puzzle/gologger"
	"puzzle/util"
)

const banner = `
	██████╗ ██╗   ██╗███████╗███████╗██╗     ███████╗
	██╔══██╗██║   ██║╚══███╔╝╚══███╔╝██║     ██╔════╝
	██████╔╝██║   ██║  ███╔╝   ███╔╝ ██║     █████╗  
	██╔═══╝ ██║   ██║ ███╔╝   ███╔╝  ██║     ██╔══╝  
	██║     ╚██████╔╝███████╗███████╗███████╗███████╗
	╚═╝      ╚═════╝ ╚══════╝╚══════╝╚══════╝╚══════╝
`

func ShowBanner() {
	gologger.Printf(banner)
	gologger.Infof("Current Version: %s", config.Version)
	if !InitOutFile() {
		gologger.Fatalf("初始化输出文件夹失败,请检查是否有权限")
	}
}

func InitOutFile() bool {
	if !util.FileExists(config.OutDir) {
		err := util.CreateDir(config.OutDir)
		if err != nil {
			return false
		}
	}
	return true
}
