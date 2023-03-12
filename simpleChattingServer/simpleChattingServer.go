package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func log(msg string, err error) {
	fmt.Println(msg, err)
	os.Exit(-1)
}

func work(sock net.Conn) {
	for {
		var buf [1024]byte
		n, err := sock.Read(buf[:])
		if err != nil {
			log("读数据出错：", err)
		}
		fmt.Printf("%v说：%v", sock.RemoteAddr().String(), string(buf[:n]))

		var str string
		reader := bufio.NewReader(os.Stdin)
		str, err = reader.ReadString('\n')
		if err != nil {
			fmt.Println("发送失败：", err)
			os.Exit(-1)
		}
		sock.Write([]byte(str))
	}
}

func main() {
	listenSock, err := net.Listen("tcp", ":8088")
	if err != nil {
		log("监听失败:", err)
	}
	for {
		connectSock, err := listenSock.Accept()
		defer connectSock.Close()
		if err != nil {
			log("接受请求失败:", err)
		}
		go work(connectSock)
	}
}
