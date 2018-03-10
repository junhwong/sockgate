package main

import(
	"fmt"
	"log"
	"github.com/sockgate/listen"
	"github.com/sockgate/common"
	"github.com/sockgate/handler"
)

func main() {
	common.Set("sock5", handler.Socket5Handler)
	listener, err := listen.NewTCPListener(":1081")
	if err!= nil {
		log.Fatalln(err)
		return
	}
	listener.Listen()
    fmt.Printf("hello, world\n")
}