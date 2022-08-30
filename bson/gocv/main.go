package main

import (
	"io/ioutil"
	"os"

	"gocv.io/x/gocv"
	"gopkg.in/mgo.v2/bson"
)

type images struct {
	Data []byte
}

func main() {
	mat := gocv.IMRead("./mug.jpg", gocv.IMReadColor)
	defer mat.Close()
	img := images{}
	img.Data = mat.ToBytes()
	b, err := bson.Marshal(img)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("matrix.txt", b, os.FileMode(0644))
	if err != nil {
		panic(err)
	}
}
