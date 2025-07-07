package main

import (
	"encoding/binary"
	"fmt"
)

type Response struct {
	messageSize   uint32
	correlationId uint32
	errCode       uint16
	arrayLen      uint8
	apiVersions   []ApiVersion
	throttleTime  uint32
	tagBuffer     uint8
}

func (res Response) calcMessageSize() uint32 {
	apiVersionSize := len(res.apiVersions) * calcVersionSize()
	return uint32(4 + 2 + 1 + apiVersionSize + 4 + 1)
}

func (res Response) toByteArray() []byte {
	arr := make([]byte, res.messageSize+4)
	binary.BigEndian.PutUint32(arr[0:4], res.messageSize)
	binary.BigEndian.PutUint32(arr[4:8], res.correlationId)
	binary.BigEndian.PutUint16(arr[8:10], res.errCode)
	arr[10] = res.arrayLen

	f := 10 + 1
	t := f

	for _, version := range res.apiVersions {
		t = f + 2
		binary.BigEndian.PutUint16(arr[f:t], version.apiKey)

		f = t
		t = f + 2
		binary.BigEndian.PutUint16(arr[f:t], version.minApiVersion)

		f = t
		t = f + 2
		binary.BigEndian.PutUint16(arr[f:t], version.maxApiVersion)

		f = t
		arr[f] = version.tagBuffer
	}

	f++
	t = f + 4
	binary.BigEndian.PutUint32(arr[f:t], res.throttleTime)

	f = t
	arr[f] = res.tagBuffer

	return arr
}

func (res Response) String() string {
	return fmt.Sprintf("Response{messageSize: %d, correlationId: %d, errCode: %d, arrayLen: %d, apiVersions: %s, throttleTime: %d}",
		res.messageSize,
		res.correlationId,
		res.errCode,
		res.arrayLen,
		res.apiVersions,
		res.throttleTime,
	)
}

type ApiVersion struct {
	apiKey        uint16
	minApiVersion uint16
	maxApiVersion uint16
	tagBuffer     uint8
}

func calcVersionSize() int {
	return 3*2 + 1
}

func (res ApiVersion) String() string {
	return fmt.Sprintf("ApiVersion{apiKey: %d, minApiVersion: %d, maxApiVersion: %d}",
		res.apiKey,
		res.minApiVersion,
		res.maxApiVersion,
	)
}

func MakeResponse(req Request) Response {
	res := Response{
		correlationId: req.correlationId,
		errCode:       0,
		apiVersions: []ApiVersion{
			{
				apiKey:        18,
				minApiVersion: 0,
				maxApiVersion: 4,
				tagBuffer:     uint8(0),
			},
		},
		throttleTime: 0,
		tagBuffer:    0,
	}
	res.messageSize = res.calcMessageSize()

	if req.apiVersion < 0 || req.apiVersion > 4 {
		res.errCode = 35
	}

	res.arrayLen = uint8(len(res.apiVersions)) + 1
	return res
}
