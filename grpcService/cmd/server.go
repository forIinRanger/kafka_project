package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"grpcservice/internal/adapters"
	"grpcservice/internal/app"
	pb "grpcservice/proto"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//nolint:all
func RunServer() {
	//upload environment variables
	godotenv.Load("../.env")
	//DB
	dbUrl := os.Getenv("DATABASE_URL")
	fmt.Println(dbUrl)
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

	//Kafka huyafka
	kfkConsumer, err := adapters.NewKafkaConsumer(context.Background(), []string{"localhost:9092"}, ps)
	if err != nil {
		panic(err)
	}
	defer kfkConsumer.Consumer.Close()
	go func() {
		err = kfkConsumer.StartCatching()

	}()

	//grpc
	serv := app.NewGrpcService(ps)
	grpcHandler := adapters.NewGRPCHandler(*serv)
	grpcServer := grpc.NewServer()
	pb.RegisterProcessorServer(grpcServer, grpcHandler)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	grpcServer.GracefulStop()
	kfkConsumer.StopCatching()

}
