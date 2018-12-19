package main

import (
	"fmt"
	"net"
	"strings"
)

var onlineConns = make(map[string]net.Conn)
var messageQueue = make(chan string, 1000)
var quitChan = make(chan bool)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func ProcessInfo(conn net.Conn) {
	buf := make([]byte, 1024)
	defer conn.Close()

	for {
		numOfBytes, err := conn.Read(buf)
		if err != nil {
			break
		}

		if numOfBytes != 0 {
			message := string(buf[:numOfBytes])
			messageQueue <- message
		}

	}
}

func ConsumeMessage() {
	for {
		select {
		case message := <-messageQueue:
			// 对消息进行解析
			doProcessMessage(message)
		case <-quitChan:
			break
		}
	}
}

func doProcessMessage(message string) {
	// 127.0.0.1:63329#你好,我是63332
	// 按照#号进行拆分addr、sendMessage
	contents := strings.Split(message, "#")
	if len(contents) > 1 {
		addr := contents[0]
		sendMessage := contents[1]

		// 处理一下防止里面有空格
		addr = strings.Trim(addr, " ")

		// 如果有数据就进行
		if conn, ok := onlineConns[addr]; ok {
			_, err := conn.Write([]byte(sendMessage))
			if err != nil {
				fmt.Println("online conns send failure!")
			}
		}
	}

}

func main() {
	listenSocket, err := net.Listen("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer listenSocket.Close()

	fmt.Println("Server is waiting ...")

	go ConsumeMessage()

	for {
		conn, err := listenSocket.Accept()
		CheckError(err)

		// 将conn存储到onlineConns映射表
		// 把打印出来的当做字符串在赋值
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		onlineConns[addr] = conn
		for i := range onlineConns {
			fmt.Println(i)
		}
		go ProcessInfo(conn)
	}
}
