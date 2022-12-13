package adapters

import (
	"context"
	"project/internal/app"
	"project/internal/app/sender"
	"project/internal/app/validator"
	"project/internal/domain"
	"project/pkg/api"
)

type MyServer struct {
	service app.MyService
}

func NewMyServer(v validator.Validator, s sender.Sender) MyServer {
	service := app.NewService(v, s)
	return MyServer{service: service}
}

// nolint:all
func (serv MyServer) PostTask(ctx context.Context, request api.PostTaskRequestObject) (api.PostTaskResponseObject, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	msg := domain.Message{Data: request.Body.Data}
	err := serv.service.HandleMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return api.PostTask200JSONResponse{Message: &request.Body.Data}, nil
}
