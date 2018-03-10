package handler

import(
	"fmt"
	"log"
	"time"
	"io"
	"net"
	"encoding/binary"
)
// Socks 5 协议 https://www.ietf.org/rfc/rfc1928.txt
func Socket5Handler(conn *net.TCPConn) bool {
	log.Println(fmt.Sprintf("[Socket5Handler] connected: %s", conn))
	// todo stream
	defer conn.Close()
	buf := make([]byte, 256)
	_, err := conn.Read(buf)
	if err != nil {
		log.Println(fmt.Sprintf("[Socket5Handler] reads error: %s", err))
		return false
	}

	/** 协商验证方式
	   The localConn connects to the dstServer, and sends a ver
	   identifier/method selection message:
		          +----+----------+----------+
		          |VER | NMETHODS | METHODS  |
		          +----+----------+----------+
		          | 1  |    1     | 1 to 255 |
		          +----+----------+----------+
	   The VER field is set to X'05' for this ver of the protocol.  The
	   NMETHODS field contains the number of method identifier octets that
	   appear in the METHODS field.
	*/
	if buf[0] != 0x05 {
		log.Println(fmt.Sprintf("[Socket5Handler] not supported socket version"))
		return false
	}

	/**
	   The dstServer selects from one of the methods given in METHODS, and
	   sends a METHOD selection message:
		          +----+--------+
		          |VER | METHOD |
		          +----+--------+
		          | 1  |   1    |
		          +----+--------+
	*/
	// 不需要验证，直接验证通过
	conn.Write([]byte{0x05, 0x00})
	//c.SetReadDeadline(time.Now().Add(10c.SetReadDeadline(time.Now().Add(10 * time.Millisecond)) * time.Millisecond))


	/** 获取获取目标访问地址
		          +----+-----+-------+------+----------+----------+
		          |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
		          +----+-----+-------+------+----------+----------+
		          | 1  |  1  | X'00' |  1   | Variable |    2     |
		          +----+-----+-------+------+----------+----------+
	*/
	
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(fmt.Sprintf("[Socket5Handler] IO error: %s", err))
		return false
	}
	if n < 7 {
		log.Println(fmt.Sprintf("[Socket5Handler] IO error: %s", err))
		return false
	}

	if buf[0] != 0x05 && (buf[1] != 0x01) {
		log.Println(fmt.Sprintf("[Socket5Handler] not supported connect method"))
		return false
	}

	switch buf[3] {
	case 0x01:
		handleIPv4(conn, buf[4 : 4+net.IPv4len], buf[n-2:])
	case 0x03:
		handleDomain(conn, buf[5:n-2], buf[n-2:])
	case 0x04:
		handleIPv6(conn, buf[4 : 4+net.IPv6len], buf[n-2:])
	default: return false
	}



	return true
}

func handleIPv4(conn *net.TCPConn, ip []byte, port []byte){
	log.Println(fmt.Sprintf("[Socket5Handler] ipv4: %s:%d", string(ip), int(binary.BigEndian.Uint16((port)))))
}

func handleIPv6(conn *net.TCPConn, ip []byte, port []byte){
	log.Println(fmt.Sprintf("[Socket5Handler] ipv6: %s:%d", string(ip), int(binary.BigEndian.Uint16((port)))))
}

func handleDomain(conn *net.TCPConn, domain []byte, port []byte) {
	defer conn.Close()
	log.Println(fmt.Sprintf("[Socket5Handler] host: %s:%d", string(domain), int(binary.BigEndian.Uint16((port)))))
	ip, _ := net.ResolveIPAddr("ip", string(domain))
	dstAddr := &net.TCPAddr{
		IP:   ip.IP,
		Port: int(binary.BigEndian.Uint16(port)),
	}
	dstServer, err := net.DialTCP("tcp", nil, dstAddr)
	if err != nil {
		log.Println(fmt.Sprintf("[Socket5Handler] connect to host error: %s", err))
		return
	}
	
	defer dstServer.Close()

	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		log.Println(fmt.Sprintf("[Socket5Handler] IO error: %s", err))
		return
	}
	go func(){
		err = transform(dstServer, conn)
		if err != nil {
			conn.Close()
			dstServer.Close()

		} 
	}()
	transform(conn, dstServer)
	
}

func transform(src *net.TCPConn, dst *net.TCPConn) error{
	buf := make([]byte, 1024)
	for{
		src.SetReadDeadline(time.Now().Add(1000 * time.Millisecond))
		reads, err := src.Read(buf)
		if err != nil {
			if(err == io.EOF){
				return nil
			}
			log.Println(fmt.Sprintf("[Socket5Handler] [conn->dist] IO error: %s", err))
			return err
		}
		writes, err := dst.Write(buf[0 : reads])
		if err != nil {
			return err
		}
		if reads != writes {
			return io.ErrShortWrite
		}
	}
}