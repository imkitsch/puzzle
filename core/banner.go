package core

import (
	"Allin/gologger"
)

const Version = "0.1"
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
	gologger.Infof("Current Version: %s\n", Version)
}
