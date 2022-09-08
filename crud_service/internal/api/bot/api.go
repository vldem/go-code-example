package api

import (
	"context"
	"io"
	"log"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/vldem/homework1/internal/config"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
	validatorPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/validator"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/counter"
	loggerPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/logger"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func New(client pb.BackendClient) pb.AdminServer {
	return &implementation{
		client: client,
	}
}

type implementation struct {
	client pb.BackendClient
	pb.UnimplementedAdminServer
}

func (i implementation) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ui/UserCreate")
	defer span.Finish()
	span.LogKV("user_email", in.GetEmail())
	span.LogKV("user_name", in.GetName())
	span.LogKV("user_role", in.GetRole())
	span.LogKV("user_password", in.GetPassword())

	counter.InRequestInc()
	if err := validatorPkg.ValidateParameters(validatorPkg.MakeParametersToValidate([]string{
		in.GetEmail(),
		in.GetName(),
		in.GetRole(),
		in.GetPassword(),
	})); err != nil {
		counter.ErrorCounterInc()
		span.LogKV("error", "validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	counter.OutRequestInc()
	out, err := i.client.UserCreate(ctx, &pb.BackendUserCreateRequest{
		Email:    in.GetEmail(),
		Name:     in.GetName(),
		Role:     in.GetRole(),
		Password: in.GetPassword(),
	})
	if err != nil {
		counter.FailedRequestInc()
		counter.ErrorCounterInc()
		span.LogKV("error", "error from backend service")
		return nil, status.Error(codes.Internal, err.Error())
	}
	counter.SuccessRequestInc()

	return &pb.UserCreateResponse{
		Id: out.GetId(),
	}, nil
}

func (i implementation) UsersAdd(ctx context.Context, in *pb.UsersAddRequest) (*pb.UsersAddResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ui/UsersAdd")
	defer span.Finish()

	counter.InRequestInc()
	data := in.Users
	if len(data) == 0 {
		span.LogKV("error", "invalid argument: data is empty")
		return nil, status.Error(codes.InvalidArgument, "data is empty")
	}
	for i := 0; i < len(data); i++ {
		if err := validatorPkg.ValidateParameters(validatorPkg.MakeParametersToValidate([]string{
			data[i].Email,
			data[i].Name,
			data[i].Role,
			data[i].Password,
		})); err != nil {
			counter.ErrorCounterInc()
			span.LogKV("error", "validation error")
			loggerPkg.Logger.Log.Debug("Validation filed",
				zap.Int("iteration:", i),
				zap.String("email:", data[i].Email),
				zap.String("name:", data[i].Name),
				zap.String("role:", data[i].Role),
				zap.String("password:", data[i].Password),
			)
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	stream, err := i.client.UsersAdd(ctx)
	if err != nil {
		counter.ErrorCounterInc()
		span.LogKV("error", "error during client connection")
		loggerPkg.Logger.Log.Error(err.Error())
		return nil, status.Error(codes.Internal, err.Error())
	}
	waitCh := make(chan struct{})
	result := make([]uint64, 0, len(data))

	var recvErr error

	go func() {
		for {
			inBackend, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitCh)
				return
			}
			if err != nil {
				counter.ErrorCounterInc()
				loggerPkg.Logger.Log.Error(err.Error())
				span.LogKV("error", "error from backend service")
				recvErr = status.Error(codes.Internal, err.Error())
				close(waitCh)
				return
			}
			result = append(result, inBackend.GetId())
			log.Printf("Got message %d", inBackend.GetId())
		}
	}()
	for _, row := range data {
		counter.OutRequestInc()
		if err := stream.Send(&pb.BackendUsersAddRequest{
			Email:    row.Email,
			Name:     row.Name,
			Role:     row.Role,
			Password: row.Password,
		}); err != nil {
			loggerPkg.Logger.Log.Error(err.Error(),
				zap.String("email:", row.Email),
				zap.String("name:", row.Name),
				zap.String("role:", row.Role),
				zap.String("password:", row.Password),
			)
			counter.ErrorCounterInc()
			span.LogKV("error", "error sending data to backend service")
			return nil, status.Error(codes.Internal, err.Error())
		}
		counter.OutRequestInc()
	}
	stream.CloseSend()
	<-waitCh

	if recvErr != nil {
		return nil, recvErr
	}

	counter.SuccessRequestInc()
	return &pb.UsersAddResponse{
		Ids: result,
	}, nil
}

func (i implementation) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ui/UserGet")
	defer span.Finish()

	//validate user's input
	counter.InRequestInc()
	if err := validatorPkg.ValidateUserId(strconv.FormatUint(in.GetId(), 10)); err != nil {
		counter.ErrorCounterInc()
		span.LogKV("error", "validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	counter.OutRequestInc()
	out, err := i.client.UserGet(ctx, &pb.BackendUserGetRequest{
		Id: in.GetId(),
	})
	if err != nil {
		counter.ErrorCounterInc()
		counter.FailedRequestInc()
		span.LogKV("error", "error from backend service")
		return nil, status.Error(codes.Internal, err.Error())
	}
	counter.SuccessRequestInc()

	return &pb.UserGetResponse{
		Id:    out.GetId(),
		Email: out.GetEmail(),
		Name:  out.GetName(),
		Role:  out.GetRole(),
	}, nil
}

func (i implementation) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ui/UserList")
	defer span.Finish()

	counter.InRequestInc()

	if in.GetOrder().GetField() != "" {
		err := validatorPkg.ValidateSortingField(in.GetOrder().GetField())
		if err != nil {
			span.LogKV("error", "validation error")
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	ctxData, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Println(ctxData.Get("custom"))
	}

	var recPerPage, pageNum uint64
	recPerPage = in.GetRecPerPage()
	pageNum = in.GetPageNum()

	sortingOrder := models.SortingOrder{
		Field:      in.GetOrder().GetField(),
		Descending: in.GetOrder().GetDescending(),
	}
	if sortingOrder.Field == "" {
		sortingOrder.Field = config.DefaultSortingField
	}

	if recPerPage == 0 {
		recPerPage = config.DefaultRecPerPage
	}
	if pageNum == 0 {
		pageNum = config.DefaultPageNum
	}

	counter.OutRequestInc()
	users, err := i.client.UserList(ctx, &pb.BackendUserListRequest{
		RecPerPage: &recPerPage,
		PageNum:    &pageNum,
		Order: &pb.BackendUserListRequest_SortingOrder{
			Field:      sortingOrder.Field,
			Descending: sortingOrder.Descending,
		},
	})
	if err != nil {
		counter.ErrorCounterInc()
		counter.FailedRequestInc()
		span.LogKV("error", "error from backend service")
		return nil, status.Error(codes.Internal, err.Error())
	}

	counter.SuccessRequestInc()

	result := make([]*pb.UserListResponse_User, 0, len(users.Users))
	for _, user := range users.Users {
		result = append(result, &pb.UserListResponse_User{
			Id:    uint64(user.GetId()),
			Email: user.GetEmail(),
			Name:  user.GetName(),
			Role:  user.GetRole(),
		})
	}

	return &pb.UserListResponse{
		Users: result,
	}, nil

}

func (i implementation) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ui/UserUpdate")
	defer span.Finish()

	// validate user's input
	counter.InRequestInc()
	if err := validatorPkg.ValidateUserId(strconv.FormatUint(in.GetId(), 10)); err != nil {
		counter.ErrorCounterInc()
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validatorPkg.ValidatePassword(in.GetOldpassword()); err != nil {
		counter.ErrorCounterInc()
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validatorPkg.ValidateParameters(validatorPkg.MakeParametersToValidate([]string{
		in.GetEmail(),
		in.GetName(),
		in.GetRole(),
		in.GetPassword(),
	})); err != nil {
		counter.ErrorCounterInc()
		span.LogKV("error", "validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	counter.OutRequestInc()
	if _, err := i.client.UserUpdate(ctx, &pb.BackendUserUpdateRequest{
		Id:          in.GetId(),
		Email:       in.GetEmail(),
		Name:        in.GetName(),
		Role:        in.GetRole(),
		Password:    in.GetPassword(),
		Oldpassword: in.GetOldpassword(),
	}); err != nil {
		counter.FailedRequestInc()
		counter.ErrorCounterInc()
		span.LogKV("error", "error from backend service")
		return nil, status.Error(codes.Internal, err.Error())
	}
	counter.SuccessRequestInc()
	return &pb.UserUpdateResponse{}, nil
}

func (i implementation) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ui/UserUpdate")
	defer span.Finish()

	// validate user's input
	counter.InRequestInc()
	if err := validatorPkg.ValidateUserId(strconv.FormatUint(in.GetId(), 10)); err != nil {
		counter.ErrorCounterInc()
		span.LogKV("error", "validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := validatorPkg.ValidatePassword(in.GetPassword()); err != nil {
		counter.ErrorCounterInc()
		span.LogKV("error", "validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	counter.OutRequestInc()
	if _, err := i.client.UserDelete(ctx, &pb.BackendUserDeleteRequest{
		Id:       in.GetId(),
		Password: in.GetPassword(),
	}); err != nil {
		counter.ErrorCounterInc()
		counter.FailedRequestInc()
		span.LogKV("error", "error from backend service")
		return nil, status.Error(codes.Internal, err.Error())
	}
	counter.SuccessRequestInc()
	return &pb.UserDeleteResponse{}, nil
}
