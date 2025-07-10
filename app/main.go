package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

var _ = net.Listen
var _ = os.Exit

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err2 := listener.Accept()
	if err2 != nil {
		fmt.Println("Error accepting connection: ", err2.Error())
		os.Exit(1)
	}
	defer conn.Close()

	for {
		handle(err2, conn)
	}
}

func handle(err2 error, conn net.Conn) {
	msgSizeArr := make([]byte, 4)
	_, err2 = conn.Read(msgSizeArr)
	if err2 != nil {
		fmt.Println("Error reading from connection: ", err2.Error())
		os.Exit(1)
	}

	msgSize := binary.BigEndian.Uint32(msgSizeArr)
	reqHeader := make([]byte, msgSize)
	_, err2 = conn.Read(reqHeader)
	if err2 != nil {
		fmt.Println("**Error reading from connection: ", err2.Error())
		os.Exit(1)
	}

	req := MakeRequest(reqHeader, msgSize)
	res := MakeResponse(req)

	_, err2 = conn.Write(res.toByteArray())
	if err2 != nil {
		fmt.Println("Error writing to connection: ", err2.Error())
		os.Exit(1)
	}
}
