package main

import (
	"xpnsec.com/node-hijack/v2/pkg/cli"

	"github.com/fatih/color"
)

func main() {
	// server := sshserver.NewSSHServer("/tmp/hijack-certs/host_ssh_signed.crt", "/tmp/hijack-certs/host.key")
	// err := server.Start()
	// if err != nil {
	// 	panic(err)
	// }
	//
	//

	theme1 := color.RGB(247, 71, 130)
	theme2 := color.RGB(90, 142, 255)

	theme1.Println(`
      __   __   ___                      __
|\ | /  \ |  \ |__  __ |__| |    |  /\  /    |__/
| \| \__/ |__/ |___    |  | | \__/ /~~\ \__, |  \
                                                  `)

	theme2.Println(`      @_xpn_
                                                 `)
	cli.Execute()

	// client, err := api.NewClient("/tmp/node.crt", "/tmp/node.key", "10.1.10.1:8443")
	// if err != nil {
	// 	panic(err)
	// }

	// err = client.GetAllNodes(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	// node, err := client.GetNode(context.Background(), "58102c12-cf6a-4fd9-b74f-8a6c0e93765f")
	// if err != nil {
	// 	panic(err)
	// }

	// node.Spec.Hostname = "older-one"

	// err = client.UpdateNode(context.Background(), node)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("%+v\n", node)

	// err = client.UpsertNode(context.Background(), "teleport-node-2")
	// if err != nil {
	// 	panic(err)
	// }

	// time.Sleep(100000)

	// err = client.DeleteNode(context.Background(), "teleport-node-2")
	// if err != nil {
	// 	panic(err)
	// }
}
