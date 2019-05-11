package ftp

import(
	"encoding/binary"
	"net"
	"unsafe"
)

var commands = map[string]uint8{
	"cd": uint8(1),
	"ls": uint8(2),
	"exit":uint8(3),
	"mkdir":uint8(4),
	"put":uint8(5),
	"get":uint8(6),
}

type ftp_conn struct{
	conn net.Conn
	cwd string
	home string
	exit bool
}

func (
