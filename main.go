package main

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/sdslabs/pinger/pkg/protobuf/pb"
)

type server struct {
	pb.AlterServicesServer
}

func (s *server) CreateCheckService(ctx context.Context, req *pb.CreateCheckRequest) (*pb.Response, error) {
	return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
}

func (s *server) DeleteCheckService(ctx context.Context, req *pb.DeleteCheckRequest) (*pb.Response, error) {
	return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
}

func (s *server) UpdateCheckService(ctx context.Context, req *pb.UpdateCheckRequest) (*pb.Response, error) {
	return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
}

func (s *server) CreatePageService(ctx context.Context, req *pb.CreatePageRequest) (*pb.Response, error) {
	return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
}

func (s *server) DeletePageService(ctx context.Context, req *pb.DeletePageRequest) (*pb.Response, error) {
	return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
}

func (s *server) UpdatePageService(ctx context.Context, req *pb.UpdatePageRequest) (*pb.Response, error) {
	return &pb.Response{Response: "gRPC Code 0 : OK "}, nil
}

func main() {
	lis, err := net.Listen("tcp", "127.0.0.1:5000")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAlterServicesServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
