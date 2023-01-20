package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

type Data struct {
	Data []byte `json:"data" bson:"data"`
	Name string `json:"Name" bson:"Name"`
}

func main() {

	mat := gocv.IMRead("mug.jpg", gocv.IMReadGrayScale)
	defer mat.Close()
	out := []uint8{}
	for j := 0; j < mat.Rows(); j++ {
		for i := 0; i < mat.Cols(); i++ {
			out = append(out, mat.GetUCharAt(j, i))
		}
	}

	fmt.Println(len(out), mat.Rows()*mat.Cols())

	var inputBytes0 []byte
	//var inputBytes1 []byte
	// Temp variable to hold our converted int32 -> []byte
	//bs := make([]byte, 4)
	for i := 0; i < len(out); i++ {
		inputBytes0 = append(inputBytes0, out[i])
		//binary.LittleEndian.PutUint16(bs, uint16(out[i]))
		//inputBytes0 = append(inputBytes0, bs...)
		//binary.LittleEndian.PutUint16(bs, uint16(out[i]))
		//inputBytes1 = append(inputBytes1, bs...)
	}
	// [3,512,512] -> [-1][(CHANNEL1, ROW1, COLUMN1),(CHANNEL1, ROW1, COLUMN2) ]

	// buf, err := gocv.IMEncode(gocv.JPEGFileExt, mat)
	// if err != nil {
	// 	panic(err)
	// }

	// defer buf.Close()
	// result := make([]byte, buf.Len())
	// copy(result, buf.GetBytes())

	// d := Data{Data: result, Name: "test"}

	// // will make it to a binary json
	// // This we can send over the wired.
	// _, err = bson.Marshal(d)
	// if err != nil {
	// 	panic(err)
	// }

	// Ask for help how to flatten the go cv mat in the best way, check with Michael of Simon about this.

}
