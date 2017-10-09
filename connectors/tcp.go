package connectors

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/kartiksura/kvstore/store"
)

var kv store.KVStore

//TCPConnector exposes a tcp server which will accept the requests in tcp and return the response
type TCPConnector struct {
	Port string
}

//Start takes in some params such as port, etc
func (t *TCPConnector) Start(port string) {
	t.Port = port
	t.Listener()
}

//Listener continuously listens for new connections and spawns a go-routine for each connection
func (t *TCPConnector) Listener() error {
	l, err := net.Listen("tcp", t.Port)
	if err != nil {
		return err
	}

	log.Println("Listen on", l.Addr().String())
	for {
		log.Println("Accept a connection request.")
		conn, err := l.Accept()
		if err != nil {
			log.Println("Failed accepting a connection request:", err)
			continue
		}
		log.Println("Handle incoming messages.")
		go t.handleCommands(conn)
	}
}

func (t *TCPConnector) handleCommands(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	var err error
	for err == nil {
		log.Print("Receive STRING message:")
		s, err := rw.ReadString('\n')
		if err != nil {
			log.Println("Cannot read from connection.\n", err)
		}
		s = strings.Trim(s, "\r\n ")
		log.Println(s)
		cmd := strings.Split(s, " ")
		switch cmd[0] {
		case "set":
			fallthrough
		case "SET":
			keyName := ""
			valSize := 0
			keyName, valSize, err = checkSETParams(cmd)
			if err == nil {
				log.Print("Parsed :", cmd, keyName)
				val, err := rw.ReadBytes('\n')
				log.Print("Parsed data:", val)

				if err == nil {
					val = bytes.Trim(val, "\r\n ")
					log.Print("Parsed data:", val)
					if valSize != len(val) {
						log.Print("Incorrect val size:expected:", valSize, len(val))
						err = fmt.Errorf("Incorrect val size")
					} else {
						log.Print("Processing command")
						err = kv.Set(keyName, val)
						if err == nil {
							rw.WriteString("OK")
						}
					}
				}
			}
		case "get":
			fallthrough
		case "GET":
			keyName := ""
			keyName, err = checkGETParams(cmd)
			if err == nil {
				log.Print("Parsed :", cmd, keyName)
				var data []byte
				data, _ = kv.Get(keyName)
				rw.WriteString("VALUE " + strconv.Itoa(len(data)) + "\n")
				data = append(data, '\n')
				_, err = rw.WriteString(string(data))
			}
		case "delete":
			fallthrough
		case "DELETE":
			if len(cmd) != 2 {
				err = fmt.Errorf("incorrect number of arguments: DELETE KEY")
			} else {
				kv.Delete(cmd[1])
			}
		case "stats":
			fallthrough
		case "STATS":
			st := kv.GetStats()
			ans := fmt.Sprintf("%+v", st)
			_, err = rw.WriteString(ans)

		}
		if err != nil {
			log.Print(err)
			_, err = rw.WriteString(err.Error())
			if err != nil {
				log.Println("Cannot write to connection.\n", err)
				return
			}
		}

		err = rw.Flush()
		if err != nil {
			log.Println("Flush failed.", err)
			return
		}
	}
}

func checkSETParams(cmd []string) (key string, valSize int, err error) {
	if len(cmd) != 3 {
		err = fmt.Errorf("incorrect number of arguments: SET KEY VALUELEN")
		return
	}
	if len(cmd[1]) > 100 {
		err = fmt.Errorf("key size exceeded: GET KEY")
		return

	}
	key = cmd[1]
	valSize, err = strconv.Atoi(cmd[2])
	return
}

func checkGETParams(cmd []string) (key string, err error) {
	if len(cmd) != 2 {
		err = fmt.Errorf("incorrect number of arguments: GET KEY")
		return
	}
	if len(cmd[1]) > 100 {
		err = fmt.Errorf("key size exceeded: GET KEY")
		return

	}
	key = cmd[1]
	return
}

func init() {
	kv = store.InitLRU(100)
}
