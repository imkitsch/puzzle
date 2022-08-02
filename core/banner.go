package core

import (
	"Allin/gologger"
	"Allin/util"
)

const banner = `                    	
	 █████╗ ██╗     ██╗     ██╗███╗   ██╗
	██╔══██╗██║     ██║     ██║████╗  ██║
	███████║██║     ██║     ██║██╔██╗ ██║
	██╔══██║██║     ██║     ██║██║╚██╗██║
	██║  ██║███████╗███████╗██║██║ ╚████║
	╚═╝  ╚═╝╚══════╝╚══════╝╚═╝╚═╝  ╚═══╝
`

func ShowBanner() {
	gologger.Printf(banner)
	gologger.Infof("Current Version: %s", Version)
	if InitOutDile() {
		gologger.Fatalf("初始化输出文件夹失败,请检查是否有权限")
	}
}

func InitOutDile() bool {
	if !util.FileExists(OutDir) {
		err := util.CreateDir(OutDir)
		if err != nil {
			return false
		}
	}
	return true
}
