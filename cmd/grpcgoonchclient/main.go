package main

import (
	"fmt"
	"log"

	grpcgoonch "github.com/thaigoonch/grpcgoonch/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn

	port := 9000
	conn, err := grpc.Dial(fmt.Sprintf(":%d", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect on port %d: %v", port, err)
	}
	defer conn.Close()
	c := grpcgoonch.NewServiceClient(conn)
	request := grpcgoonch.Request{
		Text: "Hello from Goonch Client!",
		Key:  "key",
	}

	response, err := c.CryptoRequest(context.Background(), &request)
	if err != nil {
		log.Fatalf("Error when calling CryptoRequest(): %v", err)
	}

	log.Printf("Response from Goonch Server: %s", response.Result)
}
