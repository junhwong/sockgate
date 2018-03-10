package common

import(
	"net"
)

type Handler func(*net.TCPConn) bool