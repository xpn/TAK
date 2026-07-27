package main

import (
	"github.com/fatih/color"
	"xpnsec.com/log-viewer/v2/pkg/cli"
)

func main() {

	theme1 := color.RGB(247, 71, 130)
	theme2 := color.RGB(90, 142, 255)
	theme1.Println(`
_    ____ ____    _  _ _ ____ _ _ _ ____ ____
|    |  | | __ __ |  | | |___ | | | |___ |__/
|___ |__| |__]     \/  | |___ |_|_| |___ |  \`)
	theme2.Println(`           @_xpn_
`)

	cli.Execute()
}
