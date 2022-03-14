package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	grpcgoonch "github.com/thaigoonch/grpcgoonch/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func doClientThings() {
	for {
		port := 9000
		host := "grpcgoonch-service"
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		}
		conn, err := grpc.Dial(fmt.Sprintf("dns:///%s:%d", host, port), opts...)
		if err != nil {
			grpclog.Fatalf("Could not connect on port %d: %v", port, err)
		}
		defer conn.Close()

		c := grpcgoonch.NewServiceClient(conn)

		text := "encrypt me"
		key := []byte("#89er@jdks$jmf_d")
		request := grpcgoonch.Request{
			Text: text,
			Key:  key,
		}

		response, err := c.CryptoRequest(context.Background(), &request)
		if err != nil {
			grpclog.Fatalf("Error when calling CryptoRequest(): %v", err)
		}

		log.Printf("Response from Goonch Server: %s", response.Result)
	}
}

func main() {
	start := time.Now()
	defer func() {
		fmt.Println("Execution Time: ", time.Since(start))
	}()
	wg := sync.WaitGroup{}

	for i := 0; i < 11; i++ {
		wg.Add(1)
		go func() {
			doClientThings()
			wg.Done()
		}()
	}
	wg.Wait()
}
