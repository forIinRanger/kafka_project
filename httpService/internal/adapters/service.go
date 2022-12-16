package adapters

import (
	"context"
	"httpservice/internal/app"
	"httpservice/internal/app/sender"
	"httpservice/internal/app/validator"
	"httpservice/internal/domain"
	"httpservice/pkg/api"
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
