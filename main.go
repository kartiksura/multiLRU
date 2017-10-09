package main

import "github.com/kartiksura/kvstore/connectors"

func main() {
	var tcp connectors.TCPConnector
	tcp.Start(":61000")

}
