package main

import (
	"backend/services/inventory/internal/config"
	"log"
	"net"
)

func main() {
	cfg := config.MustConfig()

	listener, err := net.Listen("tcp", cfg.GrpcPort)
	if err != nil {
		log.Fatal(err)
	}

	server := newGRPCServer(cfg)

	log.Println("inventory grpc service started", "port", cfg.GrpcPort)

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
