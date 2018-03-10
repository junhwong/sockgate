package common

import "net"

var handlers = make(map[string]Handler)

func Set(name string, handler Handler){
	handlers[name] = handler
}

func Handle(conn *net.TCPConn){
	for _, handle := range handlers {
		handle(conn)
	}
}