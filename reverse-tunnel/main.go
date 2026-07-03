package main

import (
	"fmt"

	"xpnsec.com/reverse-tunnel/v2/pkg/cli"
)

func main() {
	fmt.Printf(`
 __   ___       ___  __   __   ___    ___                 ___
|__) |__  \  / |__  |__) /__  |__  __  |  |  | |\ | |\ | |__  |
|  \ |___  \/  |___ |  \ .__/ |___     |  \__/ | \| | \| |___ |___
         @_xpn_
`)
	cli.Execute()
}
