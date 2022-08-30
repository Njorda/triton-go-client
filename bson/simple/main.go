package main

import (
	"fmt"
	"os"

	"gopkg.in/mgo.v2/bson"
)

type doc struct {
	A int
	B []int
}

type images struct {
	Data []byte
}

func main() {
	fmt.Println("Hello")

	f, err := os.Open("my_file.txt")
	if err != nil {
		panic(err)
	}
	b1 := make([]byte, 10000000)
	n, err := f.Read(b1)
	if err != nil {
		panic(err)
	}
	fmt.Println("bytes read: ", n)
	test := doc{}

	err = bson.Unmarshal(b1, &test)
	if err != nil {
		panic(err)
	}
	fmt.Println(test)

}
