package main

import (
	"database/sql"
	"net"

	_ "github.com/mattn/go-sqlite3"
	"github.com/viniciusidacruz/grpc-http2/internal/database"
	"github.com/viniciusidacruz/grpc-http2/internal/pb"
	"github.com/viniciusidacruz/grpc-http2/internal/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("sqlite3", "./db.sqlite")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	categoryDb := database.NewCategory(db)
	categoryService := services.NewCategoryService(*categoryDb)

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
