package main

import (
	"xpnsec.com/reverse-tunnel/v2/pkg/cli"
	"github.com/fatih/color"
)

func main() {
	theme1 := color.RGB(247, 71, 130)
	theme2 := color.RGB(90, 142, 255)

	theme1.Println(`
 __   ___       ___  __   __   ___    ___                 ___
|__) |__  \  / |__  |__) /__  |__  __  |  |  | |\ | |\ | |__  |
|  \ |___  \/  |___ |  \ .__/ |___     |  \__/ | \| | \| |___ |___`)

	theme2.Println(`
		@_xpn_
	`)
	cli.Execute()
}
