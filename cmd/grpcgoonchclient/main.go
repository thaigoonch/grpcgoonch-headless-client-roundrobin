package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	grpcgoonch "github.com/thaigoonch/grpcgoonch/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func doClientThings() {
	reg := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	reg.MustRegister(grpcMetrics)

	// Create an http server for prometheus
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf("0.0.0.0:%d", 9094)}

	// Start http server for prometheus
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

	for i := 0; i < 100; i++ {
		port := 9000
		host := "grpcgoonch-service"
		opts := []grpc.DialOption{
			grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
			grpc.WithTimeout(10 * time.Second),
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
