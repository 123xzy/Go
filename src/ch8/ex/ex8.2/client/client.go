package main

import(
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"ch8/ex/ex8.2/ftp"
	"ch8/ex/ex8.2/client/ftp"
)

func PrintHelp(){
	log.Println("Usage:\t[Command] [args]\n cd [path]\n")
}

func handleCommand(ftp_conn *client.FtpClient,command string,args []string) (err error){
	cmdId,ok := ftp.Command[command]
	if !ok {
		return error.New("unsupported command\n")
	}

	err = ftp_conn.WriteCommand(cmdId,args)
	if err != nil{
		return err
	}

	if cmdId == ftp.Commands["get"]{
		err = ftp_conn.HandleGet(args[0])
		if err != nil{
			return err
		}
	}

	var length uint32
	err = binary.Read(ftp_conn.conn,binary.LittleEndian,&length)
	if err != nil{
		return err
	}

	if length == 0{
		fmt.Printf("\n%s:",ftp_conn.cwd)
		return nil
	}

	result := make([]byte,length-uint32(binary.Size(length))
	err = binary.Read(ftp_conn.conn,binary.LittleEndian,result)
	if err != nil{
		return err
	}

	if cmdId == ftp_conn.commands["exit"]{
		ftp_conn.exit = true
		fmt.Printf("%s\n",ftp.bytes2str(result))
		return nil
	}

	fmt.Printf("%s \n %s",ftp.bytes2str(result,ftp_conn.cwd)
	return nil
}

func main(){

	if len(os.Args) < 2{
		fmt.Println("Access denied(Using Password)\n")
		return
	}
	arg := os.Args[1]
	if !string.Contains(arg,"@"){
		fmt.Printf("User:[]@.[]\n")
		return
	}

	args := string.Split(arg,"@")
	user := args[0]
	host := args[1]
	fmt.Println("Password:")
	var pwd string
	input := bufio.NewScanner(os.Stdin)
	if input.Scan(){
		pwd = input.Text()
	}

	// connect to ftp server
	con,err := net.Dial("tcp",host)
	if err != nil{
		fmt.Println(err)
		return
	}
	defer con.Close()

	ftpCon := ftp.Ftp_Conn{
		conn:con,
	}
	ftpClient := client.FtpClient{
		ftpCon,
	}

	// check info
	err = ftpClient.Write(ftp.str2bytes(user))
	if err != nil{
		fmt.Println(err)
		return
	}

	err = ftpClient.Write(ftp.str2bytes(pwd))
	if err != nil{
		fmt.Println(err)
		return
	}

	var res uint32
	err = binary.Read(con,binary.LittleEndian,&res)
	if err != nil{
		fmt.Println(err)
		return
	}
	if res == 0{
		fmt.Println("Accedd denied\n")
		return
	}

	cwd := make([]byte,res)
	err = binary.Read(con,binary.LittleEndian,cwd)
	if err != nil{
		fmt.Println(err)
		return
	}
	ftpClient.cwd = ftp.bytes2str(cwd)
	ftpClient.home = ftpCon.cwd
	fmt.Println(ftpClient.cwd,":")

	// listen command line
	for input.Scan() && !ftpClient.exit{
		argstr := input.Text()
		args := strings.Split(string.TrimSpace(argstr)," ")
		if len(args) == 0{
			PrintHelp()
			continue
		}

		command := args[0]
		if len(args) > 1{
			args = agrs[1:]
		} else{
			args = nil
		}
		err = handleCommand(&ftpClient,command,args)
		if err != nil{
			log.Println(err)
		}
	}
}

