package main

import (
	"encoding/binary"
	"fmt"
	"io"
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
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		err := handleRequest(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed the connection")
			} else {
				fmt.Println("Error handling request:", err)
			}
			break
		}
	}
}

func handleRequest(conn net.Conn) error {
	msgSizeArr := make([]byte, 4)
	if _, err := io.ReadFull(conn, msgSizeArr); err != nil {
		return err
	}

	msgSize := binary.BigEndian.Uint32(msgSizeArr)
	reqHeader := make([]byte, msgSize)
	if _, err := io.ReadFull(conn, reqHeader); err != nil {
		return err
	}

	req := MakeRequest(reqHeader, msgSize)
	res := MakeResponse(req)

	if _, err := conn.Write(res.toByteArray()); err != nil {
		return err
	}

	return nil
}
