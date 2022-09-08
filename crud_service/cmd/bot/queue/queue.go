package queue

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/vldem/homework1/internal/config"
	validatorPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/validator"
	loggerPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/logger"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
	"go.uber.org/zap"
)

type order struct {
	Field      string `json:"field"`
	Descending bool   `json:"descending"`
}

type messageData struct {
	Command     string `json:"command"`
	RequestData struct {
		Id          int64  `json:"id"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		Role        string `json:"role"`
		Password    string `json:"password"`
		Oldpassword string `json:"oldpassword"`
		ReqPerPage  uint64 `json:"rec_per_page"`
		PageNum     uint64 `json:"page_num"`
		Order       order  `json:"order"`
	}
}

type queueMessage struct {
	Command     string `json:"command"`
	RequestData interface{}
}
type Consumer struct {
	P sarama.SyncProducer
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
			loggerPkg.Logger.Log.Info("Done")
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				loggerPkg.Logger.Log.Info("Data channel closed")
				return nil
			}
			if err := c.ProcessQueue(session.Context(), msg.Key, msg.Value); err != nil {
				return err
			}
			session.MarkMessage(msg, "")
			return nil
		}
	}
}

func (c *Consumer) ProcessQueue(ctx context.Context, key []byte, message []byte) error {
	var requestData messageData

	if err := json.Unmarshal(message, &requestData); err != nil {
		panic(err)
		//return errors.Wrapf(err, "error during unmarshal request")
	}
	switch requestData.Command {
	case "UserList":
		loggerPkg.Logger.Log.Info("UI service",
			zap.String("message key", string(key)),
			zap.String("command", requestData.Command),
			zap.Uint64("rec_per_page", requestData.RequestData.ReqPerPage),
			zap.Uint64("page_num", requestData.RequestData.PageNum),
			zap.String("sorting field", requestData.RequestData.Order.Field),
			zap.Bool("role", requestData.RequestData.Order.Descending),
		)

		if requestData.RequestData.Order.Field != "" {
			err := validatorPkg.ValidateSortingField(requestData.RequestData.Order.Field)
			if err != nil {
				_, _, err := c.P.SendMessage(&sarama.ProducerMessage{
					Topic: config.TopicUIResponse,
					Key:   sarama.ByteEncoder(key),
					Value: sarama.ByteEncoder("invalid argument"),
				})
				if err != nil {
					return errors.Wrapf(err, "sending message to %v topic", config.TopicUIResponse)
				}
				return nil
			}
		}

		if requestData.RequestData.Order.Field == "" {
			requestData.RequestData.Order.Field = config.DefaultSortingField
		}

		if requestData.RequestData.ReqPerPage == 0 {
			requestData.RequestData.ReqPerPage = config.DefaultRecPerPage
		}
		if requestData.RequestData.PageNum == 0 {
			requestData.RequestData.PageNum = config.DefaultPageNum
		}

		msg, err := json.Marshal(&queueMessage{
			Command: "UserList",
			RequestData: &pb.BackendUserListRequest{
				RecPerPage: &requestData.RequestData.ReqPerPage,
				PageNum:    &requestData.RequestData.PageNum,
				Order: &pb.BackendUserListRequest_SortingOrder{
					Field:      requestData.RequestData.Order.Field,
					Descending: requestData.RequestData.Order.Descending,
				},
			},
		})
		if err != nil {
			return errors.Wrap(err, "marshaling message")
		}

		_, _, err = c.P.SendMessage(&sarama.ProducerMessage{
			Topic: config.TopicUIRequest,
			Key:   sarama.ByteEncoder(key),
			Value: sarama.ByteEncoder(msg),
		})
		if err != nil {
			return errors.Wrapf(err, "sending message to %v topic", config.TopicUIRequest)
		}

	}
	return nil
}
