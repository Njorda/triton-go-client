package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gocv.io/x/gocv"
	"gopkg.in/mgo.v2/bson"
)

type Frame struct {
	Image []byte `bson:"image"`
	Id    string `bson:"id"`
}

func main() {
	writeMat := gocv.IMRead("./mug.jpg", gocv.IMReadColor)
	defer writeMat.Close()

	writeImg := Frame{}
	writeImg.Image = writeMat.ToBytes()
	writeImg.Id = fmt.Sprintf("%v_test.jpg", time.Now().Unix())

	b, err := bson.Marshal(writeImg)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./matrixStruct", b, os.FileMode(0644))
	if err != nil {
		panic(err)
	}
}
