package main

import (
	"fmt"
	"log"
	"net/http"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	grpcgoonch "github.com/thaigoonch/grpcgoonch/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/grpclog"
)

var (
	port     = 9000
	promPort = 9095
)

func doClientThings(grpcMetrics *grpc_prometheus.ClientMetrics) {
	for i := 0; i < 100; i++ {
		host := "grpcgoonch-service"
		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithBalancerName("round_robin"),
			grpc.WithUnaryInterceptor(grpcMetrics.UnaryClientInterceptor()),
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
	reg := prometheus.NewRegistry()
	grpcMetrics := grpc_prometheus.NewClientMetrics()
	reg.MustRegister(grpcMetrics)

	// Create an http server for prometheus
	httpServer := &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr:    fmt.Sprintf(":%d", promPort)}

	// Start http server for prometheus
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Error starting http server: %v", err)
		}
	}()

	for i := 0; i < 11; i++ {
		doClientThings(grpcMetrics)
	}
}
