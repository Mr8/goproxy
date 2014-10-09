package main

import (
	"fmt"
	"github.com/Mr8/goproxy/utils"
	"net"
)

var CLRF = []byte{0x0d, 0x0a, 0x0d, 0x0d}

func RemoteCon(hp *utils.HttpParser) (net.Conn, error) {
	return net.Dial("tcp", fmt.Sprintf("%s:%s", hp.Host, hp.Port))
}

func HandlerCon(c net.Conn) {
	defer c.Close()
	var (
		header *utils.HttpParser
		msg    []byte
	)

	C := NewCon(c)

	sz, err := C.Read()
	if sz == 0 || err != nil {
		return
	}
	msg = C.GetMsg()
	header, err = utils.Parser(msg)
	fmt.Printf("Parsed HTTP Header--> \n%s\n", header)
	fmt.Printf("Received HTTP Header--> \n%s\n", string(msg))
	tBuf, err := TransferHTTP(C.GetMsg(), header)
	fmt.Printf("TransferHTTP--> \n%s\n", string(tBuf))

	if header == nil || msg == nil {
		return
	}

	fmt.Printf(fmt.Sprintf("Make connection to remote %s:%s", header.Host, header.Port))
	conn, err := RemoteCon(header)
	if err != nil {
		return
	}
	conn.Write(msg)

	for {
		buf := make([]byte, 10240)
		sz, err := conn.Read(buf)
		if sz == 0 || err != nil {
			return
		}
		C.Write(buf)
	}
}

func main() {
	fmt.Printf("Start proxy server\n")
	l, err := net.Listen("tcp4", "0.0.0.0:1121")
	if err != nil {
		fmt.Printf("Listen to port error :%s\n", err.Error())
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Printf("Accept connection error")
		} else {
			fmt.Printf("Accept connection %s", c.RemoteAddr())
			go HandlerCon(c)
		}
	}
}
