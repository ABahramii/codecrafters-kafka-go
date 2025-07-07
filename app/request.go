package main

import (
	"encoding/binary"
	"fmt"
)

type Request struct {
	messageSize              uint32
	apiKey                   uint16
	apiVersion               uint16
	correlationId            uint32
	clientIdLen              uint16
	clientId                 []byte
	taggedBuffer1            uint8
	compactClientIdLen       uint8
	actualClientId           []byte
	clientSoftwareVersionLen uint8
	clientSoftwareVersion    []byte
	//taggedBuffer2            int8
}

func (r Request) String() string {
	return fmt.Sprintf(
		"Request{messageSize:%d, apiKey:%d, apiVersion:%d, correlationId:%d, clientIdLen:%d, clientId:%s, taggedBuffer1:%d, compactClientIdLen:%d, actualClientId:%s, clientSoftwareVersionLen:%d, clientSoftwareVersion:%s}",
		r.messageSize,
		r.apiKey,
		r.apiVersion,
		r.correlationId,
		r.clientIdLen,
		string(r.clientId),
		r.taggedBuffer1,
		r.compactClientIdLen,
		string(r.actualClientId),
		r.clientSoftwareVersionLen,
		string(r.clientSoftwareVersion),
	)
}

func MakeRequest(req []byte, size uint32) Request {
	//messageSize := binary.BigEndian.Uint32(req[0:4])
	/*if len(req) != int(size) {
		return Request{}, fmt.Errorf("input too short for header")
	}*/

	apiKey := binary.BigEndian.Uint16(req[0:2])
	apiVersion := binary.BigEndian.Uint16(req[2:4])
	correlationId := binary.BigEndian.Uint32(req[4:8])
	clientIdLen := binary.BigEndian.Uint16(req[8:10])

	f := 10
	t := f + int(clientIdLen)
	clientId := req[f:t]

	f = t
	taggedBuffer1 := req[f]

	f++
	compactClientIdLen := req[f]

	f++
	t = f + int(clientIdLen)
	actualClientId := req[f:t]

	f = t
	clientSoftwareVersionLen := req[f]

	f++
	t = f + int(clientSoftwareVersionLen)
	clientSoftwareVersion := req[f:t]

	return Request{
		messageSize:              size,
		apiKey:                   apiKey,
		apiVersion:               apiVersion,
		correlationId:            correlationId,
		clientIdLen:              clientIdLen,
		clientId:                 clientId,
		taggedBuffer1:            taggedBuffer1,
		compactClientIdLen:       compactClientIdLen,
		actualClientId:           actualClientId,
		clientSoftwareVersionLen: clientSoftwareVersionLen,
		clientSoftwareVersion:    clientSoftwareVersion,
	}
}
