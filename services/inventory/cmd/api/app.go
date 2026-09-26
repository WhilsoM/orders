package main

import (
	inventorypb "backend/pkg/gen/inventory"
	"backend/pkg/psql"
	"backend/services/inventory/internal/adapters/postgresql"
	"backend/services/inventory/internal/adapters/transport"
	"backend/services/inventory/internal/inventory"
	"log"

	"google.golang.org/grpc"
)

func newGRPCServer() *grpc.Server {
	client := psql.PostgresClient{
		Host:     "localhost",
		User:     "postgres",
		Password: "mysecretpassword",
		Dbname:   "orders",
		Port:     "5432",
		Sslmode:  "disable",
	}

	db, err := psql.NewPostgres(client)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("database connected to inventory")

	repo := postgresql.NewInventoryRepository(db)
	svc := inventory.NewInventoryService(repo)
	handler := transport.NewGRPCServer(svc)

	server := grpc.NewServer()

	inventorypb.RegisterInventoryServiceServer(
		server,
		handler,
	)

	return server
}
