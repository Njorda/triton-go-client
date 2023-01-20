package main

import (
	"fmt"
	"testing"

	"gocv.io/x/gocv"
	"gopkg.in/mgo.v2/bson"
)

// from fib_test.go
func BenchmarkMarshalImage(b *testing.B) {
	// run the Fib function b.N times
	mat := gocv.IMRead("./mug.jpg", gocv.IMReadColor)
	defer mat.Close()
	for n := 0; n < b.N; n++ {
		_, err := marshalImage(fmt.Sprintf("%v", n), mat)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkToBytes(b *testing.B) {
	// run the Fib function b.N times
	mat := gocv.IMRead("./mug.jpg", gocv.IMReadColor)
	defer mat.Close()
	for n := 0; n < b.N; n++ {
		mat.ToBytes()
	}
}

func BenchmarkBson(b *testing.B) {
	mat := gocv.IMRead("./mug.jpg", gocv.IMReadColor)
	writeImg := Frame{}
	writeImg.Image = mat.ToBytes()
	writeImg.Id = "1234567.jpg"
	for n := 0; n < b.N; n++ {
		_, err := bson.Marshal(writeImg)
		if err != nil {
			panic(err)
		}
	}
}
