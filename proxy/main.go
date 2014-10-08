package main

import (
	"fmt"
	"github.com/Mr8/goproxy/utils"
	"net"
	"strings"
)

func HandlerCon(c net.Conn) {
	defer c.Close()
	C := NewCon(c)
	for {
		sz, err := C.Read()
		if sz == 0 || err != nil {
			return
		}
		msg := C.GetMsg()
		header, err := utils.Parser(msg)
		if err != nil {
			C.Write(C.GetMsg())
			return
		}
		fmt.Printf("Parsed HTTP Header %s", header)
		C.Write(C.GetMsg())
	}
}

func main() {
	fmt.Printf("Start proxy server\n")
	l, err := net.Listen("tcp", "127.0.0.1:1121")
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
