package app

import (
	"context"
	"httpservice/internal/app/sender"
	"httpservice/internal/app/validator"
	"httpservice/internal/domain"
)

type Service interface {
	HandleMessage(ctx context.Context, m domain.Message) error
}

type MyService struct {
	Validator validator.Validator
	Sender    sender.Sender
}

func NewService(v validator.Validator, s sender.Sender) MyService {
	return MyService{
		Validator: v,
		Sender:    s,
	}
}

func (s MyService) HandleMessage(ctx context.Context, m domain.Message) error {
	if ok, err := s.Validator.Validate(ctx, m.Data); !ok {
		return err
	}
	msgMap := sender.MappedMessage{Value: m.Data}
	if err := s.Sender.Send(ctx, msgMap); err != nil {
		return err
	}
	return nil
}
