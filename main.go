package main

import (
	"context"
	"errors"
	"github.com/sdslabs/pinger/pkg/protobuf/pb"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	pb.AlterServicesServer
}

func (s *server) CreateCheckService(ctx context.Context, req *pb.CreateCheckRequest) (*pb.Response, error) {
	if req.Check != "" {
		return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
	}
	return nil, errors.New("Inavlid request")
}

func (s *server) DeleteCheckService(ctx context.Context, req *pb.DeleteCheckRequest) (*pb.Response, error) {
	if req.UserId != "" {
		return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
	}
	return nil, errors.New("Inavlid request")
}

func (s *server) UpdateCheckService(ctx context.Context, req *pb.UpdateCheckRequest) (*pb.Response, error) {
	if req.UserId != "" {
		return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
	}
	return nil, errors.New("Inavlid request")
}

func (s *server) CreatePageService(ctx context.Context, req *pb.CreatePageRequest) (*pb.Response, error) {
	if req.NewPageString != "" {
		return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
	}
	return nil, errors.New("Inavlid request")
}

func (s *server) DeletePageService(ctx context.Context, req *pb.DeletePageRequest) (*pb.Response, error) {
	if req.PageId != "" {
		return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
	}
	return nil, errors.New("Inavlid request")
}

func (s *server) UpdatePageService(ctx context.Context, req *pb.UpdatePageRequest) (*pb.Response, error) {
	if req.PageId != "" {
		return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
	}
	return nil, errors.New("Inavlid request")
}

func main() {
	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAlterServicesServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
	grpcServer.Serve(lis)
}
