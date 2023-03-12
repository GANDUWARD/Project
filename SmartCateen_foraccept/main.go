package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
)

type Accepter struct {
	Port             *net.TCPAddr
	Final_string     string
	External_Storage string
}

func (a *Accepter) make_out(info string) {
	fmt.Println(info)
	a.Final_string = info
}
func (a *Accepter) check_if_over(l *int, outcome *string) error {
	f, err := os.Open(a.External_Storage)
	if err != nil {
		fmt.Println("外存打开错误:", err)
	}
	var number int = 0
	reader := bufio.NewReader(f)
	for {
		//多次循环读取
		_, err := reader.ReadString('\n') //以回车为分割依据将字符串读取
		if err == io.EOF {
			break
		}
		number++
	}
	//读取完毕后，如果外存数量达到要求行数就返回真
	if number == *l {
		*outcome = "finished!"
		return nil
	} else {
		return err
	}
}
func (a *Accepter) check_if_empty(send string) bool {
	return send == ""
}
func (a *Accepter) Port_setter(arg *int, send *string) error {
	*arg = 0
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8088")
	if err != nil {
		fmt.Println("ResolveTCPAddr err=", err)
	}
	a.Port = tcpAddr
	*send = "finished"
	return nil
}
func (a *Accepter) Data_writing(buf []byte) {
	f, err := os.Create(a.External_Storage)
	if err != nil {
		log.Panicln(err)
	}
	//写入文件
	_, err = f.WriteString(string(buf))
	if err != nil && err != io.EOF {
		log.Panicln(err)
	}
	f.Close()
}
func main() {
	a := new(Accepter) //默认端口8088
	a = &Accepter{&net.TCPAddr{}, "", "./out.TXT"}
	//注册服务
	rpc.Register(a)
	//开始监听
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8088")
	if err != nil {
		fmt.Println("ResolveTCPAddr err=", err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println("监听错误！", err)
	}
	defer listener.Close()
	conn, err1 := listener.Accept()
	if err1 != nil {
		fmt.Println("Accept error:", err1)
	}
	defer conn.Close()
	for {
		//定时更新
		time.Sleep(5 * time.Second)
		buf := make([]byte, 1024)
		n, err2 := conn.Read(buf) //读取对方发送的信息
		if err != nil {
			if err2 == io.EOF {
				fmt.Println("决策信息接受完毕!")
			} else {
				log.Panicln("conn.Read err", err2)
				return
			}
		}
		a.make_out(string(buf[:n]))
		if a.check_if_empty(string(buf[:n])) {
			continue
		} //发送为空就不进行重置和写入
		//创建文件,并写入
		go a.Data_writing(buf)
		time.Sleep(2 * time.Second)
		rpc.ServeConn(conn) //启用rpc服务
	}
}
