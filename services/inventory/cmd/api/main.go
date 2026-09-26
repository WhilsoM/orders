package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	server := newGRPCServer()

	log.Println("inventory grpc service started", "port", ":50051")

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
