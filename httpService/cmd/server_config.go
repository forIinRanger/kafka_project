package cmd

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"httpservice/internal/adapters"
	"httpservice/pkg/api"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// nolint:all
func RunServer() {

	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConfig.Producer.Return.Errors = true
	addrs := []string{os.Getenv("KAFKA_BROKERS")}
	log.Println(addrs)
	senderProd, err := adapters.NewKafkaSender(addrs, kafkaConfig)
	if err != nil {
		log.Fatalf("Error Kafka initialization: %v", err)
	}
	validatorProd := adapters.NewValidator()
	server := api.NewStrictHandler(adapters.NewMyServer(validatorProd, &senderProd), nil)
	handler := api.Handler(server)
	serv := http.Server{Addr: ":" + os.Getenv("HTTP_PORT"), Handler: handler}
	go func() {
		if err := serv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error while running: %v", err)
		}
	}()
	notificationCh := make(chan os.Signal, 1)
	signal.Notify(notificationCh, syscall.SIGINT, syscall.SIGTERM)
	<-notificationCh
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err = senderProd.Producer.Close()
	if err != nil {
		log.Println("Cannot close producer")
	}
	if err := serv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to exit: %v", err)
	}
	log.Println("server stopped successfully")

}
