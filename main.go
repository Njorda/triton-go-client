package main

import (
	"fmt"
	"log"

	triton "github.com/Njorda/trition-go-client/grpc-client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	url = "localhost:8001"
)

func main() {
	fmt.Println("Here we go")

	// Connect to gRPC server
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Couldn't connect to endpoint %s: %v", url, err)
	}
	defer conn.Close()

	// Set up the triton client
	client := triton.NewGRPCInferenceServiceClient(conn)

	// Check here if it is up
	serverLiveResponse := ServerLiveRequest(client)
	fmt.Printf("Triton Health - Live: %v\n", serverLiveResponse.Live)

	// Check here if it is ready ..
	serverReadyResponse := ServerReadyRequest(client)
	fmt.Printf("Triton Health - Ready: %v\n", serverReadyResponse.Ready)

	modelMetadataResponse := ModelMetadataRequest(client, FLAGS.ModelName, "")
	fmt.Println(modelMetadataResponse)

	inputData0 := make([]int32, inputSize)
	inputData1 := make([]int32, inputSize)
	for i := 0; i < inputSize; i++ {
		inputData0[i] = int32(i)
		inputData1[i] = 1
	}
	inputs := [][]int32{inputData0, inputData1}
	rawInput := Preprocess(inputs)

	/* We use a simple model that takes 2 input tensors of 16 integers
	each and returns 2 output tensors of 16 integers each. One
	output tensor is the element-wise sum of the inputs and one
	output is the element-wise difference. */
	inferResponse := ModelInferRequest(client, rawInput, FLAGS.ModelName, FLAGS.ModelVersion)

}
