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
	groupRepo "github.com/lightlink/user-service/internal/group/repository/grpc"
	service "github.com/lightlink/user-service/internal/user/delivery/grpc"
	userRepo "github.com/lightlink/user-service/internal/user/repository/postgres"
	userUC "github.com/lightlink/user-service/internal/user/usecase"
	protoGroup "github.com/lightlink/user-service/protogen/group"
	protoUser "github.com/lightlink/user-service/protogen/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	defer func() {
		if err = postgresConnect.Close(); err != nil {
			panic(err)
		}
	}()

	userRepository := userRepo.NewUserPostgresRepository(postgresConnect)

	userUsecase := userUC.NewUserUsecase(userRepository)
	userService := service.NewUserService(userUsecase)

	protoUser.RegisterUserServiceServer(grpcServer, userService)

	fmt.Println("gRPC сервер запущен на порту :8081")
	log.Fatal(grpcServer.Serve(listener))
}

func startHTTP() {
	router := mux.NewRouter()

	postgresConnect, err := connectToDB()
	if err != nil {
		log.Fatalf("Ошибка при подключении к БД: %v", err)
	}
	defer func() {
		if err = postgresConnect.Close(); err != nil {
			panic(err)
		}
	}()

	client, _ := grpc.Dial(
		fmt.Sprintf("%s:%s", os.Getenv("GROUP_SERVICE_HOST"), os.Getenv("GROUP_SERVICE_PORT")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	groupServiceClient := protoGroup.NewGroupServiceClient(client)

	userRepository := userRepo.NewUserPostgresRepository(postgresConnect)
	friendshipRepository := friendshipRepo.NewFriendshipPostgresRepository(postgresConnect)
	groupRepository := groupRepo.NewGroupGrpcRepository(&groupServiceClient)

	friendshipUsecase := friendshipUC.NewFriendshipUsecase(userRepository, friendshipRepository, groupRepository)
	friendshipHandler := friendshipDelivery.NewFriendshipHandler(friendshipUsecase)

	router.HandleFunc("/api/friend-request", friendshipHandler.SendFriendRequest).Methods("POST")
	router.HandleFunc("/api/accept-friend-request", friendshipHandler.AcceptFriendRequest).Methods("POST")
	router.HandleFunc("/api/decline-friend-request", friendshipHandler.DeclineFriendRequest).Methods("POST")
	router.HandleFunc("/api/pending-requests", friendshipHandler.GetPendingRequests).Methods("GET")
	router.HandleFunc("/api/friend-list", friendshipHandler.GetFriendList).Methods("GET")

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
