package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type dataGram struct {
	Exponent string
	Time     string
	Advice   string
}

func dataSend(jdata []byte, conn net.Conn) error {
	//朝指定ip发送数据
	_, err := conn.Write(jdata)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
func dialdst(address string, exp string, adtime string) {
	//时间为当前时间，数据为python接口传入的数据
	data := new(dataGram)
	data.Time = time.Now().Format("2006-1-7-10-41")
	data.Advice = adtime
	data.Exponent = exp
	d_json, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
	}
	defer conn.Close()
	err = dataSend(d_json, conn)
	if err != nil {
		fmt.Println(err)
	}
}
func main() {
	add := "43.143.215.213:8088"
	f, err := os.Open("out.TXT") //输出文档以out.TXT形式保存至本项目目录
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	buf := make([]byte, 1024)
	for {
		n, err5 := f.Read(buf) //从文件读取内容
		if err5 != nil {
			if err5 == io.EOF {
				fmt.Println("文件读取完毕")
			} else {
				fmt.Println("f.Read err:", err5)
			}
			return
		}
		if n == 0 {
			fmt.Println("空文件发送！")
			return
		}
		exp := string(buf[:bytes.IndexByte(buf[:], '\t')])                                    //根据制表符截取拥挤指数
		adt := string(buf[bytes.IndexByte(buf[:], '\t')+1 : bytes.IndexByte(buf[:], '\n')-1]) //根据回车截取建议时间
		dialdst(add, exp, adt)
	}
}
