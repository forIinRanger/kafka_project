package app

import (
	"context"
	"grpcservice/internal/domain"
	"grpcservice/internal/domain/repository"
)

type GrpcService struct {
	repo repository.Repository
}

func NewGrpcService(rep repository.Repository) *GrpcService {
	return &GrpcService{repo: rep}
}

func (g *GrpcService) GetStats(ctx context.Context, query string) (int, error) {
	return g.repo.GetByString(ctx, query)
}

func (g *GrpcService) ProcessMessage(ctx context.Context, query string) error {
	stats := domain.CountingLetters(domain.Text{Message: query})
	return g.repo.PutStatistics(ctx, query, stats.LettersCount)
}
