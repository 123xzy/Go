package main

import(
	"fmt"
	"bufio"
	"net"
	"time"
	"strings"
	"log"
)

func echo(c net.Conn,shout string,delay time.Duration){
	fmt.Fprintln(c,"\t",strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c,"\t",shout)
	time.Sleep(delay)
	fmt.Fprintln(c,"\t",strings.ToLower(shout))
}

func handlerConn(c net.Conn){
	input := bufio.NewScanner(c)
	for input.Scan(){
		go echo(c,input.Text(),time.Second*2)
	}
	c.Close()
}

func main(){
	l,err := net.Listen("tcp","localhost:8000")
	if err != nil{
		log.Fatal(err)
	}
	for{
		conn,err := l.Accept()
		if err != nil{
			log.Print(err)
			continue
		}
		go handlerConn(conn)
	}
}
