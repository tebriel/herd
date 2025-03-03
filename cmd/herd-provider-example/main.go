package main

import (
	// Import the provider you wish to serve over grpc
	_ "github.com/seveas/herd/provider/example"

	// And the helper library to serve it
	"github.com/seveas/herd/provider/plugin/server"
)

func main() {
	if err := server.ProviderPluginServer("example"); err != nil {
		panic(err)
	}
}
