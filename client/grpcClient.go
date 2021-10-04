package main

import (
	"fmt"

	pf "Project1/bookService"

	"google.golang.org/grpc"
)

var GrpcClient pf.BookServiceClient

func StartGrpcClient() bool {

	cc, err := grpc.Dial("localhost:8090", grpc.WithInsecure())
	if err != nil {
		return false
	}

	GrpcClient = pf.NewBookServiceClient(cc)

	fmt.Println("GRPC : ", GrpcClient)
	return true

}
