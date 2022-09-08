package queue

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/vldem/homework1/internal/config"
	userPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
	loggerPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/logger"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	P    sarama.SyncProducer
	User userPkg.Interface
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
			if err := c.ProcessQueue(msg.Key, msg.Value); err != nil {
				return err
			}
			session.MarkMessage(msg, "")
			return nil
		}
	}
}

func (c *Consumer) ProcessQueue(key []byte, message []byte) error {
	var requestData messageData

	if err := json.Unmarshal(message, &requestData); err != nil {
		panic(err)
		//return errors.Wrapf(err, "error during unmarshal request")
	}
	switch requestData.Command {
	case "UserList":
		loggerPkg.Logger.Log.Info("Backend service",
			zap.String("message key", string(key)),
			zap.String("command", requestData.Command),
			zap.Uint64("rec_per_page", requestData.RequestData.ReqPerPage),
			zap.Uint64("page_num", requestData.RequestData.PageNum),
			zap.String("sorting field", requestData.RequestData.Order.Field),
			zap.Bool("role", requestData.RequestData.Order.Descending),
		)

		sortingOrder := models.SortingOrder{
			Field:      requestData.RequestData.Order.Field,
			Descending: requestData.RequestData.Order.Descending,
		}
		var msg []byte
		users, err := c.User.List(context.Background(), requestData.RequestData.ReqPerPage, requestData.RequestData.PageNum, sortingOrder)
		if err != nil {
			msg = []byte(status.Error(codes.Internal, err.Error()).Error())
		} else {
			result := make([]*pb.BackendUserListResponse_User, 0, len(users))
			for _, user := range users {
				result = append(result, &pb.BackendUserListResponse_User{
					Id:    uint64(user.Id),
					Email: user.Email,
					Name:  user.Name,
					Role:  user.Role,
				})
			}
			msg, err = json.Marshal(result)
			if err != nil {
				panic(err)
			}
		}
		_, _, err = c.P.SendMessage(&sarama.ProducerMessage{
			Topic: config.TopicUIResponse,
			Key:   sarama.ByteEncoder(key),
			Value: sarama.ByteEncoder(msg),
		})
		if err != nil {
			return errors.Wrapf(err, "sending message to %v topic", config.TopicUIResponse)
		}

	}
	return nil
}
