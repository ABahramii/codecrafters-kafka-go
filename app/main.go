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
	fmt.Println("Logs from your program will appear here!")

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

	req := make([]byte, 1024)
	_, err = conn.Read(req)
	if err != nil {
		fmt.Println("Error reading from connection: ", err.Error())
		os.Exit(1)
	}

	corrId := binary.BigEndian.Uint32(req[8:12])

	apiVersion := binary.BigEndian.Uint16(req[6:8])
	var errCode uint16 = 0
	if apiVersion < 0 || apiVersion > 4 {
		errCode = 35
	}

	res := make([]byte, 4+4+2)
	binary.BigEndian.PutUint32(res[0:4], 0)
	binary.BigEndian.PutUint32(res[4:8], corrId)
	binary.BigEndian.PutUint16(res[8:10], errCode)
	_, err = conn.Write(res)
	if err != nil {
		fmt.Println("Error writing to connection: ", err.Error())
		os.Exit(1)
	}

}
