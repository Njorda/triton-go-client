package main

import (
	"io"
	"io/ioutil"
	"os"

	"gocv.io/x/gocv"
	"gopkg.in/mgo.v2/bson"
)

type image struct {
	Data []byte
}

func main() {
	writeMat := gocv.IMRead("./mug.jpg", gocv.IMReadColor)
	defer writeMat.Close()

	writeImg := image{}
	writeImg.Data = writeMat.ToBytes()
	b, err := bson.Marshal(writeImg)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./matrix", b, os.FileMode(0644))
	if err != nil {
		panic(err)
	}

	f, err := os.Open("./matrix")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	readImg := image{}
	b, err = io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	err = bson.Unmarshal(b, &readImg)
	if err != nil {
		panic(err)
	}

	readMat, err := gocv.NewMatFromBytes(writeMat.Rows(), writeMat.Cols(), gocv.MatTypeCV8UC3, readImg.Data)
	ok := gocv.IMWrite("./mugFromGo.jpg", readMat)
	if !ok {
		panic("not ok")
	}
}
