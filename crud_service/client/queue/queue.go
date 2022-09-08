package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/Shopify/sarama"
	"gitlab.ozon.dev/vldem/homework1/internal/config"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
	loggerPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/logger"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
)

type message struct {
	Command     string `json:"command"`
	RequestData interface{}
}

type Consumer struct {
	Id string
}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				loggerPkg.Logger.Log.Info("Data channel closed")
				return nil
			}
			if c.Id == string(msg.Key) {
				session.MarkMessage(msg, "")
				loggerPkg.Logger.Log.Info(fmt.Sprintf("response: [%v] for request Id [%v]", string(msg.Value), string(msg.Key)))
			}
			return nil
		}
	}
}

func RequestProcess(ctx context.Context, params []string) error {
	var msg []byte
	cmd := params[0]

	brokers := config.Brokers
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V2_0_0_0
	cfg.Producer.Return.Successes = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
	syncProducer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return err
	}

	switch cmd {
	case "list":
		var recPerPage, pageNum uint64
		if len(params[1:]) >= 2 {
			recPerPage, _ = strconv.ParseUint(params[1], 10, 64)
			pageNum, _ = strconv.ParseUint(params[2], 10, 64)
		}
		if recPerPage == 0 {
			recPerPage = config.DefaultRecPerPage
		}
		if pageNum == 0 {
			pageNum = config.DefaultPageNum
		}
		sortingOrder := models.SortingOrder{
			Field:      config.DefaultSortingField,
			Descending: false,
		}

		if len(params[1:]) == 4 {
			descending, _ := strconv.ParseBool(params[4])
			sortingOrder.Field = params[3]
			sortingOrder.Descending = descending
		}
		msg, err = json.Marshal(&message{
			Command: "UserList",
			RequestData: &pb.UserListRequest{
				RecPerPage: &recPerPage,
				PageNum:    &pageNum,
				Order: &pb.UserListRequest_SortingOrder{
					Field:      sortingOrder.Field,
					Descending: sortingOrder.Descending,
				},
			},
		})

		if err != nil {
			return err
		}

	}

	consumer := &Consumer{
		Id: randStringBytes(16),
	}

	par, off, err := syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: config.TopicClientRequest,
		Key:   sarama.ByteEncoder(consumer.Id),
		Value: sarama.ByteEncoder(msg),
	})
	logMsg := fmt.Sprintf("%v -> %v; %v", par, off, err)
	loggerPkg.Logger.Log.Info(logMsg)
	if err != nil {
		return err
	}

	//cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, config.ConsumerGroupUI, cfg)
	if err != nil {
		return err
	}
	// for {
	if err := client.Consume(ctx, []string{config.TopicUIResponse}, consumer); err != nil {
		loggerPkg.Logger.Log.Error(fmt.Sprintf("on consume: %v", err))
		//time.Sleep(time.Second * 10)
		return err
	}
	// }
	return nil
}

func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
