package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	//args
	fmt.Println("test os args", os.Args, len(os.Args))
	//flag
	abool := flag.Bool("bool", false, "a bool")
	times := flag.Int("times", 1, "run times")
	port := flag.String("port", ":6060", "run port")

	var bookName string
	flag.StringVar(&bookName, "book", "风声", "book name")

	flag.Parse()
	fmt.Println("abool:", *abool)
	fmt.Println("times:", *times)
	fmt.Println("port:", *port)
	fmt.Println("bookName:", bookName)
}
