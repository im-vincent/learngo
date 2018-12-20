package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func MessageSend(conn net.Conn) {
	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		readLine, _, _ := reader.ReadLine()
		input = string(readLine)

		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		}

		_, err := conn.Write([]byte(input))
		if err != nil {
			conn.Close()
			fmt.Println("client connect failure: " + err.Error())
			break
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	CheckError(err)
	defer conn.Close()

	go MessageSend(conn)

	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("您已经退出 欢迎再次使用！")
			os.Exit(0)
		}
		fmt.Println("receive server message content: " + string(buf))
	}
	fmt.Println("Client program end!")
}
