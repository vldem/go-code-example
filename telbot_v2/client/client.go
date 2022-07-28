package main

import (
	"context"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/vldem/go-code-example/telbot_v2/internal/config"
	pb "github.com/vldem/go-code-example/telbot_v2/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

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
	case "list":
		response, err := client.UserList(ctx, &pb.UserListRequest{})
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
