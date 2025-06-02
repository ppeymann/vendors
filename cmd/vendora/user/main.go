package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ppeymann/vendors.git/env"
	grpcservices "github.com/ppeymann/vendors.git/grpc-services"
	userpb "github.com/ppeymann/vendors.git/proto/user"
	"github.com/ppeymann/vendors.git/repository"
	"google.golang.org/grpc"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	lis, err := net.Listen("tcp", env.GetEnv("USER_PORT", ":50051"))

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	host := env.GetEnv("DB_HOST", "localhost")
	port := env.GetEnv("DB_PORT", "5432")
	user := env.GetEnv("DB_USER", "postgres")
	password := env.GetEnv("DB_PASSWORD", "password")
	dbname := env.GetEnv("DB_NAME", "postgres")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(pg.Open(dsn), &gorm.Config{SkipDefaultTransaction: true})
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewUserRepo(db, "postgres")
	err = repo.Migrate()
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, grpcservices.NewUserService(repo))

	log.Printf("gRPC server listening on %s", env.GetEnv("USER_PORT", ":50051"))

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
