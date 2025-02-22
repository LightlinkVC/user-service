package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	friendshipDelivery "github.com/lightlink/user-service/internal/friendship/delivery/http"
	friendshipRepo "github.com/lightlink/user-service/internal/friendship/repository/postgres"
	friendshipUC "github.com/lightlink/user-service/internal/friendship/usecase"
	service "github.com/lightlink/user-service/internal/user/delivery/grpc"
	userRepo "github.com/lightlink/user-service/internal/user/repository/postgres"
	userUC "github.com/lightlink/user-service/internal/user/usecase"
	proto "github.com/lightlink/user-service/protogen/user"
	"google.golang.org/grpc"
)

func main() {
	// === Запускаем gRPC ===
	go startGRPC()

	// === Запускаем HTTP ===
	startHTTP()
}

func startGRPC() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("Ошибка при поднятии gRPC listener'a: %v", err)
	}

	grpcServer := grpc.NewServer()

	postgresConnect, err := connectToDB()
	if err != nil {
		log.Fatalf("Ошибка при подключении к БД: %v", err)
	}
	defer postgresConnect.Close()

	userRepository := userRepo.NewUserPostgresRepository(postgresConnect)

	userUsecase := userUC.NewUserUsecase(userRepository)
	userService := service.NewUserService(userUsecase)

	proto.RegisterUserServiceServer(grpcServer, userService)

	fmt.Println("gRPC сервер запущен на порту :8081")
	log.Fatal(grpcServer.Serve(listener))
}

func startHTTP() {
	router := mux.NewRouter()

	postgresConnect, err := connectToDB()
	if err != nil {
		log.Fatalf("Ошибка при подключении к БД: %v", err)
	}
	defer postgresConnect.Close()

	userRepository := userRepo.NewUserPostgresRepository(postgresConnect)
	friendshipRepository := friendshipRepo.NewFriendshipPostgresRepository(postgresConnect)

	friendshipUsecase := friendshipUC.NewFriendshipUsecase(userRepository, friendshipRepository)
	friendshipHandler := friendshipDelivery.NewFriendshipHandler(friendshipUsecase)

	router.HandleFunc("/api/friend-request", friendshipHandler.SendFriendRequest).Methods("POST")

	fmt.Println("HTTP сервер запущен на порту :8083")
	log.Fatal(http.ListenAndServe(":8083", router))
}

func connectToDB() (*sql.DB, error) {
	postgresDSN := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", postgresDSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}
