package main

import "github.com/kartiksura/kvstore/connectors"

func main() {
	var tcp connectors.TCPConnector
	tcp.Start(":61000")

	// l := store.InitConcurrentLRU(9)
	// err := l.Set("jai", []byte("guru"))
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }
	// l.PrintEntries()

	// err = l.Set("yo", []byte("mama"))
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }
	// l.PrintEntries()

	// ans, err := l.Get("jai")
	// if err != nil {
	// 	log.Print(err)
	// } else {
	// 	log.Println("Found:", string(ans))
	// }
	// l.PrintEntries()

	// err = l.Set("hello", []byte("wo"))
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }
	// l.PrintEntries()

	// err = l.Set("yo", []byte("mama"))
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }
	// l.PrintEntries()
}
