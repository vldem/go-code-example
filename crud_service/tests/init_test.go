//go:build integration
// +build integration

package tests

import (
	"time"

	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
	"gitlab.ozon.dev/vldem/homework1/tests/config"
	"gitlab.ozon.dev/vldem/homework1/tests/postgres"
	"google.golang.org/grpc"
)

var (
	BackendClient pb.BackendClient
	Db            *postgres.TDB
)

func init() {
	cfg, err := config.FromEnv()

	conn, err := grpc.Dial(cfg.Host, grpc.WithInsecure(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		panic(err)
	}
	BackendClient = pb.NewBackendClient(conn)

	Db = postgres.NewFromEnv()
}
