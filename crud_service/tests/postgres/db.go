//go:build integration
// +build integration

package postgres

import (
	"context"
	"fmt"
	"sync"
	"testing"

	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/database"
	"gitlab.ozon.dev/vldem/homework1/tests/config"
)

type TDB struct {
	sync.Mutex
	DB *pgxpool.Pool
}

func NewFromEnv() *TDB {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, 5432, "test", "password", "gohw-test")

	pool, err := database.NewPostgres(context.Background(), psqlConn)
	if err != nil {
		panic(err)
	}

	return &TDB{DB: pool}
}

func (d *TDB) SetUp(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	d.Lock()
	d.Truncate(ctx)
}

func (d *TDB) TearDown() {
	defer d.Unlock()
	d.Truncate(context.Background())
}

func (d *TDB) Truncate(ctx context.Context) {
	q := "Truncate table users"
	if _, err := d.DB.Exec(ctx, q); err != nil {
		panic(err)
	}
	q = `ALTER SEQUENCE users_id_seq RESTART WITH 1`
	if _, err := d.DB.Exec(ctx, q); err != nil {
		panic(err)
	}
}
