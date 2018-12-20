package main

import (
	"fmt"
	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Student struct {
	Name string `json:"student_name"`
	Age  int
}

func main() {
	// 对数组类型的json
	ints := [5]int{1, 2, 3, 4, 5}
	s, err := json.Marshal(&ints)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(s))

	// 对map类型操作
	m := make(map[string]float64)
	m["zhangsan"] = 100.4
	s, err = json.Marshal(m)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(s))

	// 对对象进行操作
	student := Student{"zhangsan", 26}
	s, err = json.Marshal(student)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(s))

	//	对s进行解码
	var s4 interface{}
	err = json.Unmarshal(s, &s4)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", s4)

}
