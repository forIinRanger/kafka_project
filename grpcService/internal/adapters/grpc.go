package adapters

import (
	"context"
	"grpcservice/internal/app"
	pb "grpcservice/proto"
)

type GRPCHandler struct {
	pb.UnimplementedProcessorServer
	serv app.GrpcService
}

func NewGRPCHandler(serv app.GrpcService) *GRPCHandler {
	return &GRPCHandler{
		serv: serv,
	}
}
func (gr *GRPCHandler) ProcessData(ctx context.Context, req *pb.ProcessRequest) (*pb.ProcessResponse, error) {
	err := gr.serv.ProcessMessage(ctx, req.Data)
	if err != nil {
		return nil, err
	}
	stats, err := gr.serv.GetStats(ctx, req.Data)
	if err != nil {
		return nil, err
	}
	return &pb.ProcessResponse{LettersCount: int64(stats)}, nil

}
