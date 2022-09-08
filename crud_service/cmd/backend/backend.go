package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"gitlab.ozon.dev/vldem/homework1/cmd/backend/queue"
	apiPkg "gitlab.ozon.dev/vldem/homework1/internal/api/backend"
	"gitlab.ozon.dev/vldem/homework1/internal/config"
	configPkg "gitlab.ozon.dev/vldem/homework1/internal/config"
	userPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/database"
	loggerPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/logger"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	var err error
	loggerPkg.Logger.Log, err = zap.NewDevelopment()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// connection string
	psqlConn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DBname,
	)

	// connect to database
	// pool, err := pgxpool.Connect(ctx, psqlConn)
	// if err != nil {
	// 	log.Fatal("can't connect to database", err)
	// }
	// defer pool.Close()

	// if err := pool.Ping(ctx); err != nil {
	// 	log.Fatal("ping database error", err)
	// }

	pool, err := database.NewPostgres(ctx, psqlConn)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}
	defer pool.Close()

	// config connection
	config := pool.Config()
	config.MaxConnIdleTime = configPkg.MaxConnIdleTime
	config.MaxConnLifetime = configPkg.MaxConnLifetime
	config.MinConns = configPkg.MinConns
	config.MaxConns = configPkg.MaxConns

	redis := redis.NewClient(&redis.Options{
		Addr:     configPkg.RedisConfig.Addr,
		Password: configPkg.RedisConfig.Password,
		DB:       configPkg.RedisConfig.DbNum,
	})
	_, err = redis.Ping().Result()
	if err != nil {
		log.Fatal("can't connect to redis", err)
	}
	defer redis.Close()

	var user userPkg.Interface
	{
		user = userPkg.New(pool)
	}

	go runQueue(ctx, user)
	go runGRPCBackendServer(user, redis)
	//http server to show expvar
	http.ListenAndServe("127.0.0.1:8089", nil)
}

func runGRPCBackendServer(user userPkg.Interface, redis *redis.Client) {
	listener, err := net.Listen("tcp", ":"+config.GRPCPortBackend)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBackendServer(grpcServer, apiPkg.New(user, redis))

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func runQueue(ctx context.Context, user userPkg.Interface) {
	brokers := config.Brokers
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_0_0_0
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	client, err := sarama.NewConsumerGroup(brokers, config.ConsumerGroupClient, cfg)
	if err != nil {
		panic(err)
	}
	consumer := &queue.Consumer{
		P:    producer,
		User: user,
	}
	for {
		if err := client.Consume(ctx, []string{config.TopicUIRequest}, consumer); err != nil {
			loggerPkg.Logger.Log.Error(fmt.Sprintf("on consume: %v", err))
			time.Sleep(time.Second * 10)
		}
	}

}
