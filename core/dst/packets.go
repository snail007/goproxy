// Copyright 2014 The DST Authors. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package dst

import (
	"encoding/binary"
	"fmt"
	"net"
)

const dstHeaderLen = 12

type packetType int8

const (
	typeHandshake packetType = 0x0
	typeData                 = 0x1
	typeAck                  = 0x2
	typeNegAck               = 0x3
	typeShutdown             = 0x4
)

func (t packetType) String() string {
	switch t {
	case typeData:
		return "data"
	case typeHandshake:
		return "handshake"
	case typeAck:
		return "ack"
	case typeNegAck:
		return "negAck"
	case typeShutdown:
		return "shutdown"
	default:
		return "unknown"
	}
}

type connectionID uint32

func (c connectionID) String() string {
	return fmt.Sprintf("Ci%08x", uint32(c))
}

type sequenceNo uint32

func (s sequenceNo) String() string {
	return fmt.Sprintf("Sq%d", uint32(s))
}

type timestamp uint32

func (t timestamp) String() string {
	return fmt.Sprintf("Ts%d", uint32(t))
}

const (
	flagRequest  = 1 << 0 // This packet is a handshake request
	flagResponse = 1 << 1 // This packet is a handshake response
	flagCookie   = 1 << 2 // This packet contains a coookie challenge
)

type header struct {
	packetType packetType   // 4 bits
	flags      uint8        // 4 bits
	connID     connectionID // 24 bits
	sequenceNo sequenceNo
	timestamp  timestamp
}

func (h header) marshal(bs []byte) {
	binary.BigEndian.PutUint32(bs, uint32(h.connID&0xffffff))
	bs[0] = h.flags | uint8(h.packetType)<<4
	binary.BigEndian.PutUint32(bs[4:], uint32(h.sequenceNo))
	binary.BigEndian.PutUint32(bs[8:], uint32(h.timestamp))
}

func unmarshalHeader(bs []byte) header {
	var h header
	h.packetType = packetType(bs[0] >> 4)
	h.flags = bs[0] & 0xf
	h.connID = connectionID(binary.BigEndian.Uint32(bs) & 0xffffff)
	h.sequenceNo = sequenceNo(binary.BigEndian.Uint32(bs[4:]))
	h.timestamp = timestamp(binary.BigEndian.Uint32(bs[8:]))
	return h
}

func (h header) String() string {
	return fmt.Sprintf("header{type=%s flags=0x%x connID=%v seq=%v time=%v}", h.packetType, h.flags, h.connID, h.sequenceNo, h.timestamp)
}

type handshakeData struct {
	packetSize uint32
	connID     connectionID
	cookie     uint32
}

func (h handshakeData) marshalInto(data []byte) {
	binary.BigEndian.PutUint32(data[0:], h.packetSize)
	binary.BigEndian.PutUint32(data[4:], uint32(h.connID))
	binary.BigEndian.PutUint32(data[8:], h.cookie)
}

func (h handshakeData) marshal() []byte {
	var data [12]byte
	h.marshalInto(data[:])
	return data[:]
}

func unmarshalHandshakeData(data []byte) handshakeData {
	var h handshakeData
	h.packetSize = binary.BigEndian.Uint32(data[0:])
	h.connID = connectionID(binary.BigEndian.Uint32(data[4:]))
	h.cookie = binary.BigEndian.Uint32(data[8:])
	return h
}

func (h handshakeData) String() string {
	return fmt.Sprintf("handshake{size=%d connID=%v cookie=0x%08x}", h.packetSize, h.connID, h.cookie)
}

type packet struct {
	src  connectionID
	dst  net.Addr
	hdr  header
	data []byte
}

func (p packet) String() string {
	var dst string
	if p.dst != nil {
		dst = "dst=" + p.dst.String() + " "
	}
	switch p.hdr.packetType {
	case typeHandshake:
		return fmt.Sprintf("%spacket{src=%v %v %v}", dst, p.src, p.hdr, unmarshalHandshakeData(p.data))
	default:
		return fmt.Sprintf("%spacket{src=%v %v data[:%d]}", dst, p.src, p.hdr, len(p.data))
	}
}

func (p packet) LessSeq(seq sequenceNo) bool {
	diff := seq - p.hdr.sequenceNo
	if diff == 0 {
		return false
	}
	return diff < 1<<31
}

func (a packet) Less(b packet) bool {
	return a.LessSeq(b.hdr.sequenceNo)
}
