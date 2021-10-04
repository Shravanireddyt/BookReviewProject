package main

import (
	"fmt"
	"sync"
)

var WaitGroup sync.WaitGroup

func main() {

	ret := StartHttpServer()

	if ret == true {

		fmt.Printf("Http init Done...")

		ret = StartGrpcClient()
		if ret != true {
			fmt.Printf("Grpc init failed...")
		} else {
			HttpServer.ListenAndServe()
			WaitGroup.Add(1)
		}
		fmt.Printf("Grpc init Done...")
	}
	WaitGroup.Wait()
}
