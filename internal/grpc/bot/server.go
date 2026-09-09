package grpcbot

import (
	"context"
	"log"
	"net"

	"github.com/lbenedar/fresta/proto/grpc_bots"
	"github.com/lbenedar/fresta/proto/grpc_health"
	"google.golang.org/grpc"
)

type serverAPI struct {
	grpc_bots.UnimplementedBotsGatewayServer
	grpc_health.UnimplementedHealthServer
	// server GatewayServer
}

type GatewayServer interface {
	// GetPlatforms(ctx context.Context, in *PlatformStatusRequest, opts ...grpc.CallOption) (*PlatformStatusResponse, error)
	// GetStatus(ctx context.Context, in *GetStatusRequest, opts ...grpc.CallOption) (*StatusResult, error)
	// GetWorlds(ctx context.Context, in *GetWorldsRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[WorldsResult], error)
}

func CreateServer(port string) (*grpc.Server, net.Listener, error) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
		return nil, nil, err
	}

	grpcServer := grpc.NewServer()
	grpc_bots.RegisterBotsGatewayServer(grpcServer, &serverAPI{})
	grpc_health.RegisterHealthServer(grpcServer, &serverAPI{})
	return grpcServer, listener, err
}

func (s *serverAPI) GetStatus(ctx context.Context, statusReq *grpc_bots.GetStatusRequest) (*grpc_bots.StatusResult, error) {
	result := grpc_bots.StatusResult{Success: false}

	if statusReq.GetType() == grpc_bots.StatusType_URL_DATA {
		result.Success = true
		result.Data = "127.0.0.1"
		return &result, nil
	}

	result.ErrorMessage = "Not Found"
	return &result, nil
}

func (s *serverAPI) Check(ctx context.Context, in *grpc_health.HealthCheckRequest) (*grpc_health.HealthCheckResponse, error) {
	return &grpc_health.HealthCheckResponse{Status: grpc_health.HealthCheckResponse_SERVING}, nil
}
