package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

var final_out string

func make_out(info string) {
	fmt.Println(info)
	final_out = info
	http.HandleFunc("/out", Send)
}
func Send(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, final_out)
}
func main() {
	listener, err := net.Listen("tcp", ":8088")
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
		buf := make([]byte, 1024)
		n, err2 := conn.Read(buf) //读取对方发送的信息
		if err != nil {
			if err2 == io.EOF {
				fmt.Println("决策信息接受完毕!")
			} else {
				fmt.Println("conn.Read err", err2)
				return
			}
		}
		make_out(string(buf[:n]))
		err := http.ListenAndServeTLS(":44329", "ganduward.com_bundle.crt", "ganduward.com.key", nil)
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(5 * time.Second)
	}
}
