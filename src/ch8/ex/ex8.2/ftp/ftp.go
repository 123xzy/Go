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

type Ftp_Conn struct{
	conn net.Conn
	cwd string
	home string
	exit bool
}

func (ftp_conn *Ftp_Conn) Write(content []byte) error{
	var length uint32
	length = uint32(len(content))
	if length == 0{
		return binary.Write(ftp_conn.Con,binary.LittleEndain,&length)
	}

	length = length + uint32(binary.Size(length))
	err := binary.write(ftp_conn.con,binary.littleendain,&length)

	if err != nil{
		return err
	}

	err = binary.write(ftp_conn.con,binary.littleendain,content)
	if err != nil{
		return err
	}

	return nil
}

// transfer string to byte[]
func str2bytes(s string) (b []byte){
	*(*string)(unsafe.Pointer(&b)) = s
	*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&b)) + 2 * unsafe.Sizeof(&b))) = len(s)
	return
}

func bytes2str(b []byte)string{
	return *(*string)(unsafe.Pointer(&b))
}
