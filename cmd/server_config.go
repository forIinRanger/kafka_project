package cmd

import (
	"net/http"
	"project/internal/adapters"
	"project/pkg/api"
)

// nolint:all
func RunServer() {
	senderProd := adapters.NewKafkaSender()
	//swap to NewKafkaSender()
	validatorProd := adapters.NewValidator()
	server := api.NewStrictHandler(adapters.NewMyServer(validatorProd, &senderProd), nil)
	handler := api.Handler(server)
	http.ListenAndServe(":8080", handler)
}
