// This is Demidov Vladislav's telegram bot
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	_ "net/http/pprof"

	"github.com/Shopify/sarama"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.ozon.dev/vldem/homework1/cmd/bot/queue"
	apiPkg "gitlab.ozon.dev/vldem/homework1/internal/api/bot"
	"gitlab.ozon.dev/vldem/homework1/internal/config"
	botPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/bot"
	cmdAddPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/bot/command/add"
	cmdDeletePkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/bot/command/delete"
	cmdHelpPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/bot/command/help"
	cmdListPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/bot/command/list"
	cmdUpdatePkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/bot/command/update"
	loggerPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/logger"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var err error
	loggerPkg.Logger.Log, err = zap.NewDevelopment()

	if err != nil {
		log.Fatal("cannot initialize logger")
	}
	defer loggerPkg.Logger.Log.Sync()
	loggerPkg.Logger.Log.Info("Application started")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conns, err := grpc.Dial(":"+config.GRPCPortBackend, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		loggerPkg.Logger.Log.Fatal(err.Error())
	}

	client := pb.NewBackendClient(conns)

	var bot botPkg.Interface
	{
		bot = botPkg.MustNew()

		commandAdd := cmdAddPkg.New(client)
		bot.RegisterHandler(commandAdd)

		commandUpdate := cmdUpdatePkg.New(client)
		bot.RegisterHandler(commandUpdate)

		commandDelete := cmdDeletePkg.New(client)
		bot.RegisterHandler(commandDelete)

		commandList := cmdListPkg.New(client)
		bot.RegisterHandler(commandList)

		commandHelp := cmdHelpPkg.New(map[string]string{
			commandAdd.Name():    commandAdd.Description(),
			commandUpdate.Name(): commandUpdate.Description(),
			commandDelete.Name(): commandDelete.Description(),
			commandList.Name():   commandList.Description(),
		})
		bot.RegisterHandler(commandHelp)
	}
	go runBot(ctx, bot)
	go runGRPCServer(client)
	go runREST(ctx)
	go runQueue(ctx)
	http.ListenAndServe("127.0.0.1:8088", nil)
}

func runBot(ctx context.Context, bot botPkg.Interface) {
	if err := bot.Run(ctx); err != nil {
		log.Panic(err)
	}
}

func runGRPCServer(client pb.BackendClient) {
	listener, err := net.Listen("tcp", ":"+config.GRPCPort)
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServer(grpcServer, apiPkg.New(client))

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func runREST(ctx context.Context) {

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

func runQueue(ctx context.Context) {
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
		P: producer,
	}
	for {
		if err := client.Consume(ctx, []string{config.TopicClientRequest}, consumer); err != nil {
			loggerPkg.Logger.Log.Error(fmt.Sprintf("on consume: %v", err))
			time.Sleep(time.Second * 10)
		}
	}

}
