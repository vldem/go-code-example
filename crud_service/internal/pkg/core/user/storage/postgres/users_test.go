package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
)

func TestUserAdd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := setUp(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		userStorage := New(mockPool)

		queryGetUserByEmail := `SELECT u.id, u.email, u.full_name, r.name AS role, u.password FROM users AS u
JOIN roles AS r ON u.role = r.id WHERE u.email = $1`
		columns := []string{"id", "email", "full_name", "role", "password"}
		pgxRows := pgxpoolmock.NewRows(columns).ToPgxRows()
		mockPool.EXPECT().Query(context.Background(), queryGetUserByEmail, f.data.Email).Return(pgxRows, nil).Times(1)

		queryGetRolIdByName := `SELECT id FROM roles WHERE name = $1`
		columns = []string{"id"}
		pgxRows = pgxpoolmock.NewRows(columns).AddRow(uint8(1)).ToPgxRows()
		mockPool.EXPECT().Query(context.Background(), queryGetRolIdByName, f.data.Role).Return(pgxRows, nil).Times(1)

		queryAdd := `INSERT INTO users (email,full_name,role,password) VALUES( $1, $2, $3, $4) RETURNING id`
		columns = []string{"id"}
		pgxRows = pgxpoolmock.NewRows(columns).AddRow(uint(1)).ToPgxRows()
		mockPool.EXPECT().Query(context.Background(), queryAdd, f.data.Email, f.data.Name, uint8(1), f.data.Password).Return(pgxRows, nil).Times(1)

		// act
		result, err := userStorage.Add(context.Background(), models.User{
			Email:    f.data.Email,
			Name:     f.data.Name,
			Role:     f.data.Role,
			Password: f.data.Password,
		})

		// assert
		require.NoError(t, err)
		assert.Equal(t, f.data.Id, result)

	})

	t.Run("error", func(t *testing.T) {
		t.Run("user exists", func(t *testing.T) {
			// arrange
			f := setUp(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

			userStorage := New(mockPool)

			queryGetUserByEmail := `SELECT u.id, u.email, u.full_name, r.name AS role, u.password FROM users AS u
JOIN roles AS r ON u.role = r.id WHERE u.email = $1`
			columns := []string{"id", "email", "full_name", "role", "password"}
			pgxRows := pgxpoolmock.NewRows(columns).AddRow(
				f.data.Id,
				f.data.Email,
				f.data.Name,
				f.data.Role,
				f.data.Password).ToPgxRows()
			mockPool.EXPECT().Query(context.Background(), queryGetUserByEmail, f.data.Email).Return(pgxRows, nil).Times(1)

			// act
			_, err := userStorage.Add(context.Background(), models.User{
				Email:    f.data.Email,
				Name:     f.data.Name,
				Role:     f.data.Role,
				Password: f.data.Password,
			})

			// assert
			require.EqualError(t, err, fmt.Sprintf("user-id: [%v] user-email: [%v]: user already exists", f.data.Id, f.data.Email))
		})

		t.Run("wrong role", func(t *testing.T) {
			// arrange
			f := setUp(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

			userStorage := New(mockPool)
			wrongRole := "wrong role"

			queryGetUserByEmail := `SELECT u.id, u.email, u.full_name, r.name AS role, u.password FROM users AS u
JOIN roles AS r ON u.role = r.id WHERE u.email = $1`
			columns := []string{"id", "email", "full_name", "role", "password"}
			pgxRows := pgxpoolmock.NewRows(columns).ToPgxRows()
			mockPool.EXPECT().Query(context.Background(), queryGetUserByEmail, f.data.Email).Return(pgxRows, nil).Times(1)

			queryGetRolIdByName := `SELECT id FROM roles WHERE name = $1`
			columns = []string{"id"}
			pgxRows = pgxpoolmock.NewRows(columns).ToPgxRows()
			mockPool.EXPECT().Query(context.Background(), queryGetRolIdByName, wrongRole).Return(pgxRows, nil).Times(1)

			// act
			_, err := userStorage.Add(context.Background(), models.User{
				Email:    f.data.Email,
				Name:     f.data.Name,
				Role:     wrongRole,
				Password: f.data.Password,
			})

			// assert
			require.EqualError(t, err,
				fmt.Sprintf("storage.Add user-email: [%v] user-name: [%v]: storage.getRoleByName role: [%v]: role does not exists",
					f.data.Email, f.data.Name, wrongRole))
		})

		t.Run("db error", func(t *testing.T) {
			// arrange
			f := setUp(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

			userStorage := New(mockPool)

			queryGetUserByEmail := `SELECT u.id, u.email, u.full_name, r.name AS role, u.password FROM users AS u
JOIN roles AS r ON u.role = r.id WHERE u.email = $1`
			columns := []string{"id", "email", "full_name", "role", "password"}
			pgxRows := pgxpoolmock.NewRows(columns).ToPgxRows()
			mockPool.EXPECT().Query(context.Background(), queryGetUserByEmail, f.data.Email).Return(pgxRows, nil).Times(1)

			queryGetRolIdByName := `SELECT id FROM roles WHERE name = $1`
			columns = []string{"id"}
			pgxRows = pgxpoolmock.NewRows(columns).AddRow(uint8(1)).ToPgxRows()
			mockPool.EXPECT().Query(context.Background(), queryGetRolIdByName, f.data.Role).Return(pgxRows, nil).Times(1)

			queryAdd := `INSERT INTO users (email,full_name,role,password) VALUES( $1, $2, $3, $4) RETURNING id`
			columns = []string{"id"}
			pgxRows = pgxpoolmock.NewRows(columns).ToPgxRows()
			mockPool.EXPECT().
				Query(context.Background(), queryAdd, f.data.Email, f.data.Name, uint8(1), f.data.Password).
				Return(pgxRows, errors.New("db error")).Times(1)

			// act
			_, err := userStorage.Add(context.Background(), models.User{
				Email:    f.data.Email,
				Name:     f.data.Name,
				Role:     f.data.Role,
				Password: f.data.Password,
			})

			// assert
			require.EqualError(t, err,
				fmt.Sprintf("storage.Add user-email: [%v] user-name: [%v]: db error",
					f.data.Email, f.data.Name))
		})

	})
}

func TestUserGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// arrange
		f := setUp(t)

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

		userStorage := New(mockPool)

		queryUserGet := `SELECT u.id, u.email, u.full_name, r.name AS role, u.password FROM users AS u
JOIN roles AS r ON u.role = r.id WHERE u.id = $1`
		columns := []string{"id", "email", "full_name", "role", "password"}
		pgxRows := pgxpoolmock.NewRows(columns).AddRow(
			f.data.Id,
			f.data.Email,
			f.data.Name,
			f.data.Role,
			f.data.Password,
		).ToPgxRows()
		mockPool.EXPECT().Query(context.Background(), queryUserGet, f.data.Id).Return(pgxRows, nil) //pgx.ErrNoRows

		// act
		result, err := userStorage.Get(context.Background(), f.data.Id)

		// assert
		require.NoError(t, err)
		assert.Equal(t, &f.data, result)

	})

	t.Run("error", func(t *testing.T) {
		t.Run("internal error", func(t *testing.T) {
			// arrange
			f := setUp(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

			userStorage := New(mockPool)

			queryUserGet := `SELECT u.id, u.email, u.full_name, r.name AS role, u.password FROM users AS u
JOIN roles AS r ON u.role = r.id WHERE u.id = $1`

			mockPool.EXPECT().Query(context.Background(), queryUserGet, f.data.Id).Return(nil, errors.New("db error")) //pgx.ErrNoRows

			// act
			_, err := userStorage.Get(context.Background(), f.data.Id)

			// assert
			require.EqualError(t, err, fmt.Sprintf("storage.Get user-id: [%v]: db error", f.data.Id))
		})

		t.Run("user does not exists", func(t *testing.T) {
			// arrange
			f := setUp(t)

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockPool := pgxpoolmock.NewMockPgxPool(ctrl)

			userStorage := New(mockPool)

			queryUserGet := `SELECT u.id, u.email, u.full_name, r.name AS role, u.password FROM users AS u
JOIN roles AS r ON u.role = r.id WHERE u.id = $1`
			columns := []string{"id", "email", "full_name", "role", "password"}
			pgxRows := pgxpoolmock.NewRows(columns).ToPgxRows()
			mockPool.EXPECT().Query(context.Background(), queryUserGet, f.data.Id).Return(pgxRows, nil)

			// act
			_, err := userStorage.Get(context.Background(), f.data.Id)

			// assert
			require.EqualError(t, err, fmt.Sprintf("storage.Get user-id: [%v]: user does not exists", f.data.Id))
		})
	})

}
