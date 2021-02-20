package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	v2 "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/krish512/ambassador_logservice_es/pkg/elastic"
	l "github.com/krish512/ambassador_logservice_es/pkg/logger"
	"github.com/krish512/ambassador_logservice_es/pkg/sink"
)

func main() {

	// Initialize resources
	logger := l.InitLogger()
	elastic.InitElasticsearch()

	grpcServer := grpc.NewServer()
	v2.RegisterAccessLogServiceServer(grpcServer, sink.New())

	l, err := net.Listen("tcp", "0.0.0.0:9001")
	if err != nil {
		logger.Error("Failed to open port", zap.Error(err))
	}

	// Handling SIGTERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		grpcServer.GracefulStop()
		os.Exit(1)
	}()

	// Start server
	logger.Info("Listening on tcp://0.0.0.0:9001")
	grpcServer.Serve(l)

}
