package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	"gitlab.ozon.dev/vldem/homework1/internal/auth"
	"gitlab.ozon.dev/vldem/homework1/internal/config"
	userPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
	validatorPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/validator"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/counter"
	loggerPkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/logger"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func New(user userPkg.Interface, redis *redis.Client) *implementation {
	return &implementation{
		user:  user,
		cache: redis,
	}
}

type implementation struct {
	pb.UnimplementedBackendServer
	user  userPkg.Interface
	cache *redis.Client
}

func (i implementation) UserCreate(ctx context.Context, in *pb.BackendUserCreateRequest) (*pb.BackendUserCreateResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "backend/UserCreate")
	defer span.Finish()

	if err := validatorPkg.ValidateParameters(validatorPkg.MakeParametersToValidate([]string{
		in.GetEmail(),
		in.GetName(),
		in.GetRole(),
		in.GetPassword(),
	})); err != nil {
		span.LogKV("error", "validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := i.user.Create(ctx, models.User{
		Email:    in.GetEmail(),
		Name:     in.GetName(),
		Role:     in.GetRole(),
		Password: auth.GenHashPassword(in.GetPassword()),
	})
	if err != nil {
		span.LogKV("error", "db error")
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.BackendUserCreateResponse{
		Id: uint64(id),
	}, nil
}

func (i implementation) UserGet(ctx context.Context, in *pb.BackendUserGetRequest) (*pb.BackendUserGetResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "backend/UserGet")
	defer span.Finish()

	cacheKey := "UserGet:" + strconv.FormatUint(in.GetId(), 10)
	var user *models.User

	cacheResult, err := i.cache.Get(cacheKey).Result()
	if err == redis.Nil {
		user, err = i.user.Get(ctx, uint(in.GetId()))
		if err != nil {
			span.LogKV("error", "db error")
			return nil, status.Error(codes.Internal, err.Error())
		}

		dataToCache, err := json.Marshal(user)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if err := i.cache.Set(cacheKey, string(dataToCache), 24*time.Hour).Err(); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		counter.CacheMisInc()
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		if err := json.Unmarshal([]byte(cacheResult), &user); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		counter.CacheHitInc()
	}

	return &pb.BackendUserGetResponse{
		Id:    uint64(user.Id),
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}, nil
}

func (i implementation) UserList(ctx context.Context, in *pb.BackendUserListRequest) (*pb.BackendUserListResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "backend/UserList")
	defer span.Finish()

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
	cacheKey := "UserList:" + strconv.FormatUint(recPerPage, 10) +
		":" + strconv.FormatUint(pageNum, 10) +
		":" + in.GetOrder().GetField() +
		":" + strconv.FormatBool(in.GetOrder().GetDescending())
	var users []models.User

	cacheResult, err := i.cache.Get(cacheKey).Result()
	if err == redis.Nil {
		users, err = i.user.List(ctx, recPerPage, pageNum, sortingOrder)
		if err != nil {
			span.LogKV("error", "db error")
			return nil, status.Error(codes.Internal, err.Error())
		}
		dataToCache, err := json.Marshal(users)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if err := i.cache.Set(cacheKey, string(dataToCache), 24*time.Hour).Err(); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		counter.CacheMisInc()
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	} else {
		if err := json.Unmarshal([]byte(cacheResult), &users); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		counter.CacheHitInc()
	}

	result := make([]*pb.BackendUserListResponse_User, 0, len(users))
	for _, user := range users {
		result = append(result, &pb.BackendUserListResponse_User{
			Id:    uint64(user.Id),
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		})
	}

	return &pb.BackendUserListResponse{
		Users: result,
	}, nil
}

func (i implementation) UserUpdate(ctx context.Context, in *pb.BackendUserUpdateRequest) (*pb.BackendUserUpdateResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "backend/UserUpdate")
	defer span.Finish()

	if err := validatorPkg.ValidateParameters(validatorPkg.MakeParametersToValidate([]string{
		in.GetEmail(),
		in.GetName(),
		in.GetRole(),
		in.GetPassword(),
	})); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := i.user.Get(ctx, uint(in.GetId()))
	if err != nil {
		span.LogKV("error", "db error")
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = auth.VerifyPassword(*user, in.GetOldpassword()); err != nil {
		span.LogKV("error", "validation error")
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	user = &models.User{
		Id:       uint(in.GetId()),
		Email:    in.GetEmail(),
		Name:     in.GetName(),
		Role:     in.GetRole(),
		Password: auth.GenHashPassword(in.GetPassword()),
	}

	if err := i.user.Update(ctx, *user); err != nil {
		span.LogKV("error", "db error")
		return nil, status.Error(codes.Internal, err.Error())
	}

	cacheKey := "UserGet:" + strconv.FormatUint(in.GetId(), 10)
	if i.cache.Get(cacheKey).Err() == nil {
		//Key exists in cache. Let's update cache with new data.
		dataToCache, err := json.Marshal(user)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

		if err := i.cache.Set(cacheKey, string(dataToCache), 24*time.Hour).Err(); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}

	}

	if err := i.InvalidateCacheUserList(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.BackendUserUpdateResponse{}, nil
}

func (i implementation) UserDelete(ctx context.Context, in *pb.BackendUserDeleteRequest) (*pb.BackendUserDeleteResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "backend/UserDelete")
	defer span.Finish()

	user, err := i.user.Get(ctx, uint(in.GetId()))
	if err != nil {
		span.LogKV("error", "db error")
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err = auth.VerifyPassword(*user, in.GetPassword()); err != nil {
		span.LogKV("error", "validation error")
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	if err := i.user.Delete(ctx, uint(in.GetId())); err != nil {
		span.LogKV("error", "db error")
		return nil, status.Error(codes.Internal, err.Error())
	}

	cacheKey := "UserGet:" + strconv.FormatUint(in.GetId(), 10)
	if err := i.cache.Del(cacheKey).Err(); err != nil {
		loggerPkg.Logger.Log.Error(fmt.Sprintf("error during key deletion from cache [%v]", err))
	}

	if err := i.InvalidateCacheUserList(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.BackendUserDeleteResponse{}, nil
}

func (i implementation) UsersAdd(stream pb.Backend_UsersAddServer) error {
	ctx := stream.Context()

	span, ctx := opentracing.StartSpanFromContext(ctx, "backend/UsersAdd")
	defer span.Finish()

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Error(codes.Internal, err.Error())
		}

		id, err := i.user.Create(ctx, models.User{
			Email:    in.GetEmail(),
			Name:     in.GetName(),
			Role:     in.GetRole(),
			Password: auth.GenHashPassword(in.GetPassword()),
		})
		if err != nil {
			span.LogKV("error", "db error")
			return status.Error(codes.Internal, err.Error())
		}

		if err := stream.Send(&pb.BackendUsersAddResponse{
			Id: uint64(id),
		}); err != nil {
			return status.Error(codes.Internal, err.Error())
		}
	}
}

func (i implementation) InvalidateCacheUserList() error {
	cacheKeysPattern := "UserList:*"
	cacheKeys, err := i.cache.Keys(cacheKeysPattern).Result()
	if err != nil && err != redis.Nil {
		return status.Error(codes.Internal, err.Error())
	}
	if len(cacheKeys) > 0 {
		if err := i.cache.Del(cacheKeys...).Err(); err != nil {
			loggerPkg.Logger.Log.Error(fmt.Sprintf("error during keys deletion from cache [%v]", err))
		}
	}
	return nil
}
