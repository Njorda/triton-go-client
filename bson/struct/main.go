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
	b, err := marshalImage(fmt.Sprintf("%v_test.jpg", time.Now().Unix()), writeMat)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./matrixStruct", b, os.FileMode(0644))
	if err != nil {
		panic(err)
	}
}

func marshalImage(id string, mat gocv.Mat) ([]byte, error) {
	writeImg := Frame{}
	writeImg.Image = mat.ToBytes()
	writeImg.Id = id

	b, err := bson.Marshal(writeImg)
	if err != nil {
		return nil, err
	}
	return b, nil
}
