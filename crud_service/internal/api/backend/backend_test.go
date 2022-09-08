package backend

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/vldem/homework1/internal/auth"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
)

func TestUserCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().
			Create(f.Ctx, models.User{
				Email:    f.data.Email,
				Name:     f.data.Name,
				Role:     f.data.Role,
				Password: auth.GenHashPassword(f.data.Password),
			}).Return(uint(1), nil).Times(1)

		// act
		resp, err := f.service.UserCreate(f.Ctx, &pb.BackendUserCreateRequest{
			Email:    f.data.Email,
			Name:     f.data.Name,
			Role:     f.data.Role,
			Password: f.data.Password,
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, uint64(f.data.Id), resp.GetId())
	})

	t.Run("error", func(t *testing.T) {
		t.Run("invalid argument", func(t *testing.T) {
			// arrange
			f := userSetUp(t)

			// act
			_, err := f.service.UserCreate(f.Ctx, &pb.BackendUserCreateRequest{
				Email:    "",
				Name:     "Bob",
				Role:     "User",
				Password: "123456",
			})

			// assert
			require.EqualError(t, err, "rpc error: code = InvalidArgument desc = bad email <>")
		})

		t.Run("internal error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			f.userRepo.EXPECT().
				Create(f.Ctx, models.User{
					Email:    f.data.Email,
					Name:     f.data.Name,
					Role:     f.data.Role,
					Password: auth.GenHashPassword(f.data.Password),
				}).Return(uint(0), errors.New("db error")).Times(1)
			// act
			_, err := f.service.UserCreate(f.Ctx, &pb.BackendUserCreateRequest{
				Email:    f.data.Email,
				Name:     f.data.Name,
				Role:     f.data.Role,
				Password: f.data.Password,
			})

			// assert
			require.EqualError(t, err, "rpc error: code = Internal desc = db error")
		})
	})

}

func TestUserGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().
			Get(f.Ctx, f.data.Id).Return(&models.User{
			Id:    f.data.Id,
			Email: f.data.Email,
			Name:  f.data.Name,
			Role:  f.data.Role,
		}, nil).Times(1)

		// act
		resp, err := f.service.UserGet(f.Ctx, &pb.BackendUserGetRequest{
			Id: uint64(f.data.Id),
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, uint64(f.data.Id), resp.GetId())
		assert.Equal(t, f.data.Email, resp.GetEmail())
		assert.Equal(t, f.data.Name, resp.GetName())
		assert.Equal(t, f.data.Role, resp.GetRole())
	})

	t.Run("error", func(t *testing.T) {
		// arrange
		f := userSetUp(t)

		f.userRepo.EXPECT().
			Get(f.Ctx, f.data.Id).Return(nil, errors.New("db error")).Times(1)

		// act
		_, err := f.service.UserGet(f.Ctx, &pb.BackendUserGetRequest{
			Id: uint64(f.data.Id),
		})

		// assert
		require.EqualError(t, err, "rpc error: code = Internal desc = db error")
	})

}

func TestUserList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userListSetUp(t)
		f.userRepo.EXPECT().
			List(f.Ctx, f.recPerPage, f.pageNum, models.SortingOrder{
				Field:      f.order.Field,
				Descending: f.order.Descending,
			}).Return(f.data, nil).Times(1)

		// act
		resp, err := f.service.UserList(f.Ctx, &pb.BackendUserListRequest{
			RecPerPage: &f.recPerPage,
			PageNum:    &f.pageNum,
			Order: &pb.BackendUserListRequest_SortingOrder{
				Field:      f.order.Field,
				Descending: f.order.Descending,
			},
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, f.list, resp.GetUsers())
	})

	t.Run("error", func(t *testing.T) {
		t.Run("invalid argument", func(t *testing.T) {
			// arrange
			f := userListSetUp(t)

			// act
			_, err := f.service.UserList(f.Ctx, &pb.BackendUserListRequest{
				RecPerPage: &f.recPerPage,
				PageNum:    &f.pageNum,
				Order: &pb.BackendUserListRequest_SortingOrder{
					Field:      "SomeWrongField",
					Descending: f.order.Descending,
				},
			})

			// assert
			require.EqualError(t, err, "rpc error: code = InvalidArgument desc = bad sorting field <SomeWrongField>")
		})

		t.Run("internal error", func(t *testing.T) {
			// arrange
			f := userListSetUp(t)
			f.userRepo.EXPECT().
				List(f.Ctx, f.recPerPage, f.pageNum, models.SortingOrder{
					Field:      f.order.Field,
					Descending: f.order.Descending,
				}).Return(nil, errors.New("db error")).Times(1)

			// act
			_, err := f.service.UserList(f.Ctx, &pb.BackendUserListRequest{
				RecPerPage: &f.recPerPage,
				PageNum:    &f.pageNum,
				Order: &pb.BackendUserListRequest_SortingOrder{
					Field:      f.order.Field,
					Descending: f.order.Descending,
				},
			})

			// assert
			require.EqualError(t, err, "rpc error: code = Internal desc = db error")
		})
	})

}

func TestUserUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().
			Get(f.Ctx, f.data.Id).Return(&models.User{
			Id:       f.data.Id,
			Email:    f.data.Email,
			Name:     f.data.Name,
			Role:     f.data.Role,
			Password: auth.GenHashPassword(f.data.Password),
		}, nil).Times(1)

		f.userRepo.EXPECT().
			Update(f.Ctx, models.User{
				Id:       f.data.Id,
				Email:    f.data.Email,
				Name:     f.data.Name,
				Role:     f.data.Role,
				Password: auth.GenHashPassword(f.data.Password),
			}).Return(nil).Times(1)

		// act
		resp, err := f.service.UserUpdate(f.Ctx, &pb.BackendUserUpdateRequest{
			Id:          uint64(f.data.Id),
			Email:       f.data.Email,
			Name:        f.data.Name,
			Role:        f.data.Role,
			Password:    f.data.Password,
			Oldpassword: f.data.Password,
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, &pb.BackendUserUpdateResponse{}, resp)
	})

	t.Run("error", func(t *testing.T) {
		t.Run("invalid argument", func(t *testing.T) {
			// arrange
			f := userSetUp(t)

			// act
			_, err := f.service.UserCreate(f.Ctx, &pb.BackendUserCreateRequest{
				Email:    "",
				Name:     "Bob",
				Role:     "User",
				Password: "123456",
			})

			// assert
			require.EqualError(t, err, "rpc error: code = InvalidArgument desc = bad email <>")
		})

		t.Run("permission denied", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			f.userRepo.EXPECT().
				Get(f.Ctx, f.data.Id).Return(&models.User{
				Id:       f.data.Id,
				Email:    f.data.Email,
				Name:     f.data.Name,
				Role:     f.data.Role,
				Password: auth.GenHashPassword(f.data.Password),
			}, nil).Times(1)

			// act
			_, err := f.service.UserUpdate(f.Ctx, &pb.BackendUserUpdateRequest{
				Id:          uint64(f.data.Id),
				Email:       f.data.Email,
				Name:        f.data.Name,
				Role:        f.data.Role,
				Password:    f.data.Password,
				Oldpassword: "6543221",
			})

			// assert
			require.EqualError(t, err, "rpc error: code = PermissionDenied desc = wrong old password")
		})

		t.Run("internal error", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			f.userRepo.EXPECT().
				Get(f.Ctx, f.data.Id).Return(&models.User{
				Id:       f.data.Id,
				Email:    f.data.Email,
				Name:     f.data.Name,
				Role:     f.data.Role,
				Password: auth.GenHashPassword(f.data.Password),
			}, nil).Times(1)
			f.userRepo.EXPECT().
				Update(f.Ctx, models.User{
					Id:       f.data.Id,
					Email:    f.data.Email,
					Name:     f.data.Name,
					Role:     f.data.Role,
					Password: auth.GenHashPassword(f.data.Password),
				}).Return(errors.New("db error")).Times(1)

			// act
			_, err := f.service.UserUpdate(f.Ctx, &pb.BackendUserUpdateRequest{
				Id:          uint64(f.data.Id),
				Email:       f.data.Email,
				Name:        f.data.Name,
				Role:        f.data.Role,
				Password:    f.data.Password,
				Oldpassword: f.data.Password,
			})

			// assert
			require.EqualError(t, err, "rpc error: code = Internal desc = db error")
		})

	})

}

func TestUserDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := userSetUp(t)
		f.userRepo.EXPECT().
			Get(f.Ctx, f.data.Id).Return(&models.User{
			Id:       f.data.Id,
			Email:    f.data.Email,
			Name:     f.data.Name,
			Role:     f.data.Role,
			Password: auth.GenHashPassword(f.data.Password),
		}, nil).Times(1)
		f.userRepo.EXPECT().
			Delete(f.Ctx, f.data.Id).Return(nil).Times(1)

		// act
		resp, err := f.service.UserDelete(f.Ctx, &pb.BackendUserDeleteRequest{
			Id:       uint64(f.data.Id),
			Password: f.data.Password,
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, &pb.BackendUserDeleteResponse{}, resp)
	})

	t.Run("error", func(t *testing.T) {
		t.Run("permission denied", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			f.userRepo.EXPECT().
				Get(f.Ctx, f.data.Id).Return(&models.User{
				Id:       f.data.Id,
				Email:    f.data.Email,
				Name:     f.data.Name,
				Role:     f.data.Role,
				Password: auth.GenHashPassword(f.data.Password),
			}, nil).Times(1)

			// act
			_, err := f.service.UserDelete(f.Ctx, &pb.BackendUserDeleteRequest{
				Id:       uint64(f.data.Id),
				Password: "654321",
			})

			// assert
			require.EqualError(t, err, "rpc error: code = PermissionDenied desc = wrong old password")
		})

		t.Run("user not found", func(t *testing.T) {
			// arrange
			f := userSetUp(t)
			f.userRepo.EXPECT().
				Get(f.Ctx, uint(2)).Return(nil, errors.New("user not found")).Times(1)

			// act
			_, err := f.service.UserDelete(f.Ctx, &pb.BackendUserDeleteRequest{
				Id:       uint64(2),
				Password: f.data.Password,
			})

			// assert
			require.EqualError(t, err, "rpc error: code = Internal desc = user not found")
		})

	})

}
