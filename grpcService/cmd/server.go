package cmd

import (
	"context"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"grpcservice/internal/adapters"
	"grpcservice/internal/app"
	pb "grpcservice/proto"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//nolint:all
func RunServer() {
	//upload environment variables
	if os.Getenv("DOCKER_ENV") == "" {
		// Мы НЕ в Docker, загрузи .env файл
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found")
		}
	}
	//DB
	dbUrl := os.Getenv("DATABASE_URL")
	log.Println(dbUrl)
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		panic(err)
	}
	ps := adapters.NewPostgresRepo(db)
	serv := app.NewGrpcService(ps)
	//Kafka huyafka
	kfkConsumer, err := adapters.NewKafkaConsumer(context.Background(), []string{"kafka:9092"}, serv)
	if err != nil {
		panic(err)
	}
	defer kfkConsumer.Consumer.Close()
	go func() {
		err = kfkConsumer.StartCatching()

	}()

	//grpc

	grpcHandler := adapters.NewGRPCHandler(*serv)
	grpcServer := grpc.NewServer()
	pb.RegisterProcessorServer(grpcServer, grpcHandler)
	lis, err := net.Listen("tcp", ":"+os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("Error trying to connect grpc port: %w", err)
	}
	go func() {
		if err = grpcServer.Serve(lis); lis != nil {
			log.Fatalf("Error while running: %w", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	grpcServer.GracefulStop()
	kfkConsumer.StopCatching()

}
