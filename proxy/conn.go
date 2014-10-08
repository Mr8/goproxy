package main

import (
	"bufio"
	"fmt"
	"net"
)

type Connection struct {
	buffer []byte
	con    net.Conn
	bufsz  int
	reader *bufio.Reader
}

func NewCon(c net.Conn) *Connection {
	return &Connection{
		buffer: make([]byte, 10240),
		con:    c,
		bufsz:  0,
		reader: bufio.NewReader(c)}
}

func (self *Connection) Read() (int, error) {
	sz, err := self.reader.Read(self.buffer[self.bufsz:])
	if err != nil {
		fmt.Printf("Failed to Read:%s\n", err.Error())
		return 0, err
	}
	self.bufsz += sz
	return sz, err
}

func (self *Connection) Write(msg []byte) (int, error) {
	sz, err := self.con.Write(msg)
	return sz, err
}

func (self *Connection) GetMsg() []byte {
	return self.buffer
}
