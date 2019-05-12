package client

import(
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"ch8/ex/ex8.2/ftp"
)

type FtpClient struct{
	ftp.Ftp_Conn
}

func (*ftp_conn *FtpClient)WriteCommand(cmdId uint8,args []string) error{
	if cmdId == ftp.Commands["put"]{
		return ftp_conn.WritePut(cmdId,args[0])
	}

	var length uint32
	argstr := string.Join(args,"")
	length = uint32(binary.Size(length) + binary.Size(cmdId)) + uint32(len(argstr))

	err := binary.Write(ftp_conn,binary.LittleEndian,length)
	if err != nil{
		return err
	}


	err := binary.Write(ftp_conn,binary.LittleEndian,cmdId)
	if err != nil{
		return err
	}

	err := binary.Write(ftp_conn,binary.LittleEndian,ftp.str2bytes(argstr))
	if err != nil{
		return err
	}

	return nil
}

func (ftpCon *FtpClient)WritePut(cmdId uint8,filePath string) error{
	filePath = strings.Replace(filePath,"\\","/",-1)
	f,err := os.Open(filePath)
	if err  != nil{
		return err
	}
	defer f.Close()

	var length uint32
	fileName := ftp.str2bytes(path.Base(filePath)
	length = uint32(binary.Size(length) + binary.Size(cmdId)) + uint32(len(fileName))

	err = binary.Write(ftpCon.Con,binary.LittleEndian,length)
	if err != nil{
		return err
	}

	err = binary.Write(ftpCon.Con,binary.LittleEndian,cmdId)
	if err != nil{
		return err
	}

	err = binary.Write(ftpCon.Con,binary.LittleEndian,fileName)
	if err != nil{
		return err
	}

	//send file size
	fileInfo,err := f.Stat()
	if err != nil{
		return err
	}
	if fileInfo.IsDir(){
		return errors.New("Put not support dir,please try putdir\n")
	}else{
		err = binary.Write(ftpCon.Con,binary.LittleEndian,fileInfo.Size())
		if err != nil{
			return err
		}
	}

	//send file content
	buf := make([]byte,4096)
	bufReader := bufio.NewReader(f)
	for{
		n,err := bufReader.Read(buf)
		if err != nil{
			if err == io.EOF{
				break
			}
			return err
		}

		err = binary.Write(ftpCon.Con,binary.LittleEndian,buf[0:])
		if err != nil{
			return err
		}

	}
	return nil
}


