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
	message := grpcgoonch.Message{
		Body: "Hello from Goonch Client!",
	}

	response, err := c.SayHello(context.Background(), &message)
	if err != nil {
		log.Fatalf("Error when calling SayHello(): %v", err)
	}

	log.Printf("Response from Server: %s", response.Body)
}
