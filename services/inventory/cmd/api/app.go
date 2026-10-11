package main

import (
	inventorypb "backend/pkg/gen/inventory"
	"backend/pkg/psql"
	"backend/services/inventory/internal/adapters/postgresql"
	"backend/services/inventory/internal/adapters/transport"
	"backend/services/inventory/internal/config"
	"backend/services/inventory/internal/inventory"
	"log"
	"time"

	"google.golang.org/grpc"
)

func newGRPCServer(cfg *config.Config) *grpc.Server {
	client := psql.PostgresClient{
		Host:            cfg.DBHost,
		User:            cfg.DBUser,
		Password:        cfg.DBPassword,
		Dbname:          cfg.DBDbname,
		Port:            cfg.DBPort,
		Sslmode:         cfg.DBSslmode,
		MaxOpenConns:    cfg.MaxOpenConns,
		MaxIdleConns:    cfg.MaxIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime * time.Minute,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime * time.Minute,
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
