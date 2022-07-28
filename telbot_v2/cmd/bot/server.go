package main

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	apiPkg "github.com/vldem/go-code-example/telbot_v2/internal/api"
	userPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user"
	pb "github.com/vldem/go-code-example/telbot_v2/pkg/api"

	"github.com/vldem/go-code-example/telbot_v2/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runGRPCServer(user userPkg.Interface) {
	listener, err := net.Listen("tcp", ":"+config.GRPCPort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServer(grpcServer, apiPkg.New(user))

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func runREST() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)
	//ctx, cancel := context.WithTimeout(ctx, shortDuration)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcherREST),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterAdminHandlerFromEndpoint(ctx, mux, ":"+config.GRPCPort, opts); err != nil {
		panic(err)
	}

	if err := http.ListenAndServe(":"+config.HTTPPort, mux); err != nil {
		panic(err)
	}
}

func headerMatcherREST(key string) (string, bool) {
	switch key {
	case "Custom":
		return key, true
	default:
		return key, false
	}
}
