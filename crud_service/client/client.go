package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"gitlab.ozon.dev/vldem/homework1/client/queue"
	"gitlab.ozon.dev/vldem/homework1/internal/config"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
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
		log.Fatal(err)
	}

	args := os.Args[1:]
	var params []string
	var cmd string = "list"
	if len(args) > 0 {
		params = strings.Split(args[0], ";")
		cmd = params[0]
	}

	conns, err := grpc.Dial(":"+config.GRPCPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewAdminClient(conns)

	ctx := context.Background()

	switch cmd {
	case "queue":
		if err := queue.RequestProcess(ctx, params[1:]); err != nil {
			log.Fatal(err)
		}
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

		response, err := client.UserList(ctx, &pb.UserListRequest{
			RecPerPage: &recPerPage,
			PageNum:    &pageNum,
			Order: &pb.UserListRequest_SortingOrder{
				Field:      sortingOrder.Field,
				Descending: sortingOrder.Descending,
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("response: [%v]", response)
	case "get":
		id, _ := strconv.ParseUint(params[1], 10, 64)
		response, err := client.UserGet(ctx, &pb.UserGetRequest{
			Id: id,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("response: [%v]", response)
	case "add":
		if len(params) < 5 {
			log.Fatal(errors.New("invalid arguments"))
		}
		response, err := client.UserCreate(ctx, &pb.UserCreateRequest{
			Email:    params[1],
			Name:     params[2],
			Role:     params[3],
			Password: params[4],
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("response: [%v]", response)
	case "addList":
		if len(params) < 2 {
			log.Fatal(errors.New("invalid arguments"))
		}
		type userType struct {
			Email    string `json:"email"`
			Name     string `json:"name"`
			Role     string `json:"role"`
			Password string `json:"password"`
		}
		var user []userType
		if err := json.Unmarshal([]byte(params[1]), &user); err != nil {
			log.Fatal(err)
		}

		input := make([]*pb.UsersAddRequest_User, 0, len(user))
		for _, row := range user {
			input = append(input, &pb.UsersAddRequest_User{
				Email:    row.Email,
				Name:     row.Name,
				Role:     row.Role,
				Password: row.Password,
			})
		}

		response, err := client.UsersAdd(ctx, &pb.UsersAddRequest{
			Users: input,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("response: [%v]", response)
	case "update":
		if len(params) < 7 {
			log.Fatal(errors.New("invalid arguments"))
		}
		id, _ := strconv.ParseUint(params[1], 10, 64)
		response, err := client.UserUpdate(ctx, &pb.UserUpdateRequest{
			Id:          id,
			Email:       params[2],
			Name:        params[3],
			Role:        params[4],
			Password:    params[5],
			Oldpassword: params[6],
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("response: [%v]", response)
	case "delete":
		if len(params) < 3 {
			log.Fatal(errors.New("invalid arguments"))
		}
		id, _ := strconv.ParseUint(params[1], 10, 64)
		response, err := client.UserDelete(ctx, &pb.UserDeleteRequest{
			Id:       id,
			Password: params[2],
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("response: [%v]", response)
	}

}
