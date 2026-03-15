package server

import (
	"context"

	pb "github.com/HelenaBlack/anti-bruteforce/api/gen"
	"github.com/HelenaBlack/anti-bruteforce/internal/app"
)

type GRPCServer struct {
	pb.UnimplementedAntibruteforceServer
	service *app.AntiBruteforceService
}

func NewGRPCServer(service *app.AntiBruteforceService) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	ok, err := s.service.Check(ctx, req.Login, req.Password, req.Ip)
	if err != nil {
		return nil, err
	}
	return &pb.CheckResponse{Ok: ok}, nil
}

func (s *GRPCServer) Reset(ctx context.Context, req *pb.ResetRequest) (*pb.ResetResponse, error) {
	err := s.service.Reset(ctx, req.Login, req.Ip)
	if err != nil {
		return nil, err
	}
	return &pb.ResetResponse{Ok: true}, nil
}

func (s *GRPCServer) AddToBlacklist(ctx context.Context, req *pb.SubnetRequest) (*pb.SubnetResponse, error) {
	err := s.service.AddToBlacklist(ctx, req.Subnet)
	if err != nil {
		return nil, err
	}
	return &pb.SubnetResponse{Ok: true}, nil
}

func (s *GRPCServer) RemoveFromBlacklist(ctx context.Context, req *pb.SubnetRequest) (*pb.SubnetResponse, error) {
	err := s.service.RemoveFromBlacklist(ctx, req.Subnet)
	if err != nil {
		return nil, err
	}
	return &pb.SubnetResponse{Ok: true}, nil
}

func (s *GRPCServer) AddToWhitelist(ctx context.Context, req *pb.SubnetRequest) (*pb.SubnetResponse, error) {
	err := s.service.AddToWhitelist(ctx, req.Subnet)
	if err != nil {
		return nil, err
	}
	return &pb.SubnetResponse{Ok: true}, nil
}

func (s *GRPCServer) RemoveFromWhitelist(ctx context.Context, req *pb.SubnetRequest) (*pb.SubnetResponse, error) {
	err := s.service.RemoveFromWhitelist(ctx, req.Subnet)
	if err != nil {
		return nil, err
	}
	return &pb.SubnetResponse{Ok: true}, nil
}
