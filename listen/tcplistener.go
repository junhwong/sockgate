package listen

import(
	"net"
	"fmt"
	"log"
	"github.com/sockgate/common"
)

type TCPListener struct {

	Local *net.TCPAddr
}

func NewTCPListener(addr string) (*TCPListener, error) {
	listenAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &TCPListener{Local: listenAddr}, nil
}

func (listener *TCPListener) Listen() error{
	inner, err := net.ListenTCP("tcp", listener.Local)
	if err != nil {
		return err
	}
	defer inner.Close()
	log.Println(fmt.Sprintf("start linstening on: %s", listener.Local))
	for {
		conn, err:= inner.AcceptTCP()
		if err!= nil {
			log.Println(err)
			continue
		}
		go common.Handle(conn)
	}

	return nil
}