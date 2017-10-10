package main

import "github.com/kartiksura/multiLRU/connectors"

func main() {
	var tcp connectors.TCPConnector
	tcp.Start(":61000")

}
