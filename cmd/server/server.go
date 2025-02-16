package main

import (
	"fmt"
	"log"
	"net"

	"github.com/shellkah/averse/cachepb"
	"github.com/shellkah/averse/config"
	"github.com/shellkah/averse/internal"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	fmt.Printf("Loaded config: %+v\n", cfg)

	address := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	s := internal.NewServer(cfg.Cache.Capacity)
	cachepb.RegisterCacheServiceServer(grpcServer, s)

	log.Println("gRPC server is listening on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
