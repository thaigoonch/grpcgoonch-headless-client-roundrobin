package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	grpcgoonch "github.com/thaigoonch/grpcgoonch-headless/service"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

var (
	port        = 9000
	reg         = prometheus.NewRegistry()
	reqsMetrics = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "grpcgoonchheadlessclientroundrobin_requests_sent_total",
		Help: "The number of records sent from grpcgoonch-headless-client-roundrobin",
	})
)

func init() {
	reg.MustRegister(reqsMetrics)
	_, err := reg.Gather()
	if err != nil {
		log.Fatalf("Prometheus metric registration error: %v", err)
	}
}

func main() {
	pusher := push.New("http://prometheus-pushgateway:9091", "grpcgoonchheadlessclientroundrobin").Gatherer(reg)

	host := "grpcgoonch-headless-service"
	opts := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
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

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 200; i++ {
				response, err := c.CryptoRequest(context.Background(), &request)
				if err != nil {
					grpclog.Fatalf("Error when calling CryptoRequest(): %v", err)
				}
				log.Printf("Response from Goonch Server: %s", response.Result)
				reqsMetrics.Inc()
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if err := pusher.Add(); err != nil {
		log.Printf("Could not push to Pushgateway: %v", err)
	}
}
