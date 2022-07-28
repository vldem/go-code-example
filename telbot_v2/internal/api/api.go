package api

import (
	"context"
	"time"

	"github.com/vldem/go-code-example/telbot_v2/internal/auth"
	userPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user"
	"github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/models"
	validatorPkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/validator"
	pb "github.com/vldem/go-code-example/telbot_v2/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const shortDuration = 10 * time.Millisecond

func New(user userPkg.Interface) pb.AdminServer {
	return &implementation{
		user: user,
	}
}

type implementation struct {
	pb.UnimplementedAdminServer
	user userPkg.Interface
}

func (i implementation) UserCreate(ctx context.Context, in *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {

	if err := validatorPkg.ValidateParameters(validatorPkg.MakeParametersToValidate([]string{
		in.GetEmail(),
		in.GetName(),
		in.GetRole(),
		in.GetPassword(),
	})); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx, cancel := context.WithTimeout(ctx, shortDuration)
	defer cancel()
	sleepCh := make(chan struct{})
	var err error

	go func(ch chan struct{}) {
		err = i.user.Create(models.User{
			Email:    in.GetEmail(),
			Name:     in.GetName(),
			Role:     models.GetRoleId(in.GetRole()),
			Password: auth.GenHashPassword(in.GetPassword()),
		})
		sleepCh <- struct{}{}
	}(sleepCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err() // prints "context deadline exceeded"
	case <-sleepCh:
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UserCreateResponse{}, nil
}

func (i implementation) UserGet(ctx context.Context, in *pb.UserGetRequest) (*pb.UserGetResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, shortDuration)
	defer cancel()
	sleepCh := make(chan struct{})

	var user *models.User
	var err error

	go func(ch chan struct{}) {
		user, err = i.user.Get(uint(in.GetId()))
		ch <- struct{}{}
	}(sleepCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err() // prints "context deadline exceeded"
	case <-sleepCh:
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UserGetResponse{
		Id:    uint64(user.Id),
		Email: user.Email,
		Name:  user.Name,
		Role:  models.GetRoleName(user.Role),
	}, nil
}

func (i implementation) UserList(ctx context.Context, in *pb.UserListRequest) (*pb.UserListResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, shortDuration)
	defer cancel()

	// ctxData, ok := metadata.FromIncomingContext(ctx)
	// if ok {
	// 	log.Println(ctxData.Get("custom"))
	// }

	sleepCh := make(chan struct{})
	var users []models.User

	go func(ch chan struct{}) {
		users = i.user.List()
		//time.Sleep(1 * time.Second)
		ch <- struct{}{}
	}(sleepCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err() // prints "context deadline exceeded"
	case <-sleepCh:
	}

	result := make([]*pb.UserListResponse_User, 0, len(users))
	for _, user := range users {
		result = append(result, &pb.UserListResponse_User{
			Id:    uint64(user.Id),
			Email: user.Email,
			Name:  user.Name,
			Role:  models.GetRoleName(user.Role),
		})
	}

	return &pb.UserListResponse{
		Users: result,
	}, nil
}

func (i implementation) UserUpdate(ctx context.Context, in *pb.UserUpdateRequest) (*pb.UserUpdateResponse, error) {
	if err := validatorPkg.ValidateParameters(validatorPkg.MakeParametersToValidate([]string{
		in.GetEmail(),
		in.GetName(),
		in.GetRole(),
		in.GetPassword(),
	})); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx, cancel := context.WithTimeout(ctx, shortDuration)
	defer cancel()
	sleepCh := make(chan struct{})
	var user *models.User
	var err error

	go func(ch chan struct{}) {
		user, err = i.user.Get(uint(in.GetId()))
		sleepCh <- struct{}{}
	}(sleepCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err() // prints "context deadline exceeded"
	case <-sleepCh:
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = auth.VerifyPassword(*user, in.GetOldpassword()); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	go func(ch chan struct{}) {
		err = i.user.Update(models.User{
			Id:       uint(in.GetId()),
			Email:    in.GetEmail(),
			Name:     in.GetName(),
			Role:     models.GetRoleId(in.GetRole()),
			Password: auth.GenHashPassword(in.GetPassword()),
		})
		sleepCh <- struct{}{}
	}(sleepCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err() // prints "context deadline exceeded"
	case <-sleepCh:
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UserUpdateResponse{}, nil
}

func (i implementation) UserDelete(ctx context.Context, in *pb.UserDeleteRequest) (*pb.UserDeleteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, shortDuration)
	defer cancel()
	sleepCh := make(chan struct{})
	var user *models.User
	var err error

	go func(ch chan struct{}) {
		user, err = i.user.Get(uint(in.GetId()))
		sleepCh <- struct{}{}
	}(sleepCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err() // prints "context deadline exceeded"
	case <-sleepCh:
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = auth.VerifyPassword(*user, in.GetPassword()); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	go func(ch chan struct{}) {
		// time.Sleep(1 * time.Second)
		err = i.user.Delete(uint(in.GetId()))
		sleepCh <- struct{}{}
	}(sleepCh)

	select {
	case <-ctx.Done():
		return nil, ctx.Err() // prints "context deadline exceeded"
	case <-sleepCh:
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UserDeleteResponse{}, nil
}
