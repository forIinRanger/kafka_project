package cmd

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"project/internal/adapters"
	"project/pkg/api"
	"syscall"
	"time"
)

// nolint:all
func RunServer() {
	senderProd := adapters.NewSenderNaive()
	//swap to NewKafkaSender()
	validatorProd := adapters.NewValidator()
	server := api.NewStrictHandler(adapters.NewMyServer(validatorProd, &senderProd), nil)
	handler := api.Handler(server)
	serv := http.Server{Addr: ":8080", Handler: handler}
	go func() {
		if err := serv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error while running: %w", err)
		}
	}()
	notificationCh := make(chan os.Signal, 1)
	signal.Notify(notificationCh, syscall.SIGINT, syscall.SIGTERM)
	<-notificationCh
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := serv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to exit: %w", err)
	}
	log.Println("server stopped successfully")

}
