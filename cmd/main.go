package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
	service "github.com/lightlink/user-service/internal/user/delivery/grpc"
	"github.com/lightlink/user-service/internal/user/repository/postgres"
	"github.com/lightlink/user-service/internal/user/usecase"
	proto "github.com/lightlink/user-service/protogen/user"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
		log.Fatalf("Ошибка при поднятии listener'a: %v", err)
	}

	grpcServer := grpc.NewServer()

	postgresDSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	postgresConnect, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		log.Fatalf("Ошибка при подключении к БД: %v", err)
	}

	defer func() {
		if err = postgresConnect.Close(); err != nil {
			panic(err)
		}
	}()

	userRepository := postgres.NewUserPostgresRepository(postgresConnect)
	userUsecase := usecase.NewUserUsecase(userRepository)
	userService := service.NewUserService(userUsecase)
	proto.RegisterUserServiceServer(grpcServer, userService)

	fmt.Println("gRPC сервер запущен на порту :8081")
	log.Fatal(grpcServer.Serve(listener))
}
