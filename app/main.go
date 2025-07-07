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
	l, err := net.Listen("tcp", "0.0.0.0:9092")
	if err != nil {
		fmt.Println("Failed to bind to port 9092")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	msgSizeArr := make([]byte, 4)
	_, err = conn.Read(msgSizeArr)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	msgSize := binary.BigEndian.Uint32(msgSizeArr)
	reqHeader := make([]byte, msgSize)
	_, err = conn.Read(reqHeader)
	if err != nil {
		fmt.Println("**Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	req := MakeRequest(reqHeader, msgSize)
	res := MakeResponse(req)

	_, err = conn.Write(res.toByteArray())
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}
}
