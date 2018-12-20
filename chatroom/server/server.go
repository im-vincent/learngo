package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const LOG_DIRECTORY = "./test.log"

var onlineConns = make(map[string]net.Conn)
var messageQueue = make(chan string, 1000)
var quitChan = make(chan bool)
var logFile *os.File
var logger *log.Logger

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func ProcessInfo(conn net.Conn) {
	buf := make([]byte, 1024)

	// 退出时删除onlineConns的key
	defer func(conn net.Conn) {
		addr := fmt.Sprintf("%s", conn.RemoteAddr())
		delete(onlineConns, addr)
		conn.Close()

		for i := range onlineConns {
			fmt.Println("now online conns: " + i)
		}
	}(conn)

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
		sendMessage := strings.Join(contents[1:], "#")

		// 处理一下防止里面有空格
		addr = strings.Trim(addr, " ")

		// 如果有数据就进行
		if conn, ok := onlineConns[addr]; ok {
			_, err := conn.Write([]byte(sendMessage))
			if err != nil {
				fmt.Println("online conns send failure!")
			}
		}
	} else {
		contents := strings.Split(message, "*")
		if strings.ToUpper(contents[1]) == "LIST" {
			var ips string = ""
			for i := range onlineConns {
				ips = ips + "|" + i
			}
			if conn, ok := onlineConns[contents[0]]; ok {
				_, err := conn.Write([]byte(ips))
				if err != nil {
					fmt.Println("online conns send failure!")
				}
			}
		}
	}

}

func main() {
	logFile, err := os.OpenFile(LOG_DIRECTORY, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println("log file create failure!")
		os.Exit(-1)
	}

	listenSocket, err := net.Listen("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer listenSocket.Close()

	logger = log.New(logFile, "\r\n", log.Ldate|log.Ltime|log.Llongfile)

	fmt.Println("Server is waiting ...")

	logger.Println("I am writing the logs ...")

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
