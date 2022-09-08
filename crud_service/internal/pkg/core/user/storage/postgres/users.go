package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/driftprogramming/pgxpoolmock"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
	storagePkg "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/storage"
)

const poolSize = 10

var ErrUserNotExists = errors.New("user does not exists")
var ErrRoleNotExists = errors.New("role does not exists")
var ErrUserExists = errors.New("user already exists")

type Storage struct {
	pool pgxpoolmock.PgxPool //*pgxpool.Pool
}

func New(pool pgxpoolmock.PgxPool) storagePkg.Interface {
	return &Storage{
		pool: pool,
	}
}

func (s *Storage) List(ctx context.Context, recPerPage, pageNum uint64, sortingOrder models.SortingOrder) ([]models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage/List")
	defer span.Finish()

	limit := recPerPage
	offset := (pageNum - 1) * limit

	sortingField := models.GetSortingFieldName(sortingOrder.Field)
	descending := ""
	if sortingOrder.Descending {
		descending = "DESC"
	}

	query := fmt.Sprintf("SELECT u.id, u.email, u.full_name, r.name AS role FROM users AS u	JOIN roles AS r ON u.role = r.id ORDER BY %s %s LIMIT $1 OFFSET $2", sortingField, descending)

	var result []models.User
	if err := pgxscan.Select(ctx, s.pool, &result, query, limit, offset); err != nil {
		span.LogKV("error", "sql error")
		return nil, errors.Wrap(err, "storage.List: select")
	}
	return result, nil
}

func (s *Storage) Add(ctx context.Context, user models.User) (uint, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage/Add")
	defer span.Finish()

	foundUser, err := s.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, ErrUserNotExists) {
		return 0, errors.Wrapf(err, "storage.Add user-email: [%s] user-name: [%s]", user.Email, user.Name)
	}
	if foundUser != nil {
		return 0, errors.Wrapf(ErrUserExists, "user-id: [%s] user-email: [%s]", strconv.FormatUint(uint64(foundUser.Id), 10), foundUser.Email)
	}

	roleId, err := s.GetRoleIdByName(ctx, user.Role)
	if err != nil {
		return 0, errors.Wrapf(err, "storage.Add user-email: [%s] user-name: [%s]", user.Email, user.Name)
	}

	query := `INSERT INTO users (email,full_name,role,password) VALUES( $1, $2, $3, $4) RETURNING id`

	//row := s.pool.QueryRow(ctx, query, user.Email, user.Name, roleId, user.Password)
	rows, err := s.pool.Query(ctx, query, user.Email, user.Name, roleId, user.Password)
	if err != nil {
		span.LogKV("error", "sql error")
		return 0, errors.Wrapf(err, "storage.Add user-email: [%s] user-name: [%s]", user.Email, user.Name)
	}
	var id uint
	if err := pgxscan.ScanOne(&id, rows); err != nil {
		//if err := row.Scan(&id); err != nil {
		span.LogKV("error", "scanone error")
		return 0, errors.Wrapf(err, "storage.Add user-email: [%s] user-name: [%s]", user.Email, user.Name)
	}

	return id, nil
}

func (s *Storage) Update(ctx context.Context, user models.User) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage/Update")
	defer span.Finish()

	// checks that new email does not belong some other user
	foundUser, _ := s.GetUserByEmail(ctx, user.Email)
	if foundUser != nil && foundUser.Id != user.Id {
		return errors.Wrapf(ErrUserExists, "user-id: [%s] user-email: [%s]", strconv.FormatUint(uint64(foundUser.Id), 10), foundUser.Email)
	}

	// checks that the user to update exists
	foundUser, _ = s.Get(ctx, user.Id)
	if foundUser == nil {
		return errors.Wrapf(ErrUserNotExists, "user-id: [%s]", strconv.FormatUint(uint64(foundUser.Id), 10))
	}

	roleId, err := s.GetRoleIdByName(ctx, user.Role)
	if err != nil {
		return errors.Wrapf(err, "storage.Add user-email: [%s] user-name: [%s]", user.Email, user.Name)
	}

	query := `UPDATE users SET email = $2, full_name = $3, role = $4, password = $5 WHERE id = $1`

	result, err := s.pool.Exec(ctx, query, user.Id, user.Email, user.Name, roleId, user.Password)
	if err != nil {
		span.LogKV("error", "sql error")
		return errors.Wrapf(err, "storage.Update user-id: [%s]  ", strconv.FormatUint(uint64(user.Id), 10))
	}

	if result.RowsAffected() == 0 {
		return errors.Errorf("cannot update user. storage.Update user-id: [%s]  ", strconv.FormatUint(uint64(user.Id), 10))
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, id uint) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage/Delete")
	defer span.Finish()

	query := `DELETE FROM users WHERE id = $1`

	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		span.LogKV("error", "sql error")
		return errors.Wrapf(err, "storage.Delete user-id: [%s]  ", strconv.FormatUint(uint64(id), 10))
	}

	if result.RowsAffected() == 0 {
		return errors.Errorf("cannot delete user. storage.Update user-id: [%s]  ", strconv.FormatUint(uint64(id), 10))
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, id uint) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage/Get")
	defer span.Finish()

	query := `SELECT u.id, u.email, u.full_name, r.name AS role, u.password FROM users AS u
JOIN roles AS r ON u.role = r.id WHERE u.id = $1`
	rows, err := s.pool.Query(ctx, query, id)
	if err != nil {
		span.LogKV("error", "sql error")
		return nil, errors.Wrapf(err, "storage.Get user-id: [%s]", strconv.FormatUint(uint64(id), 10))
	}
	var user models.User
	if err := pgxscan.ScanOne(&user, rows); err != nil {
		return nil, errors.Wrapf(ErrUserNotExists, "storage.Get user-id: [%s]", strconv.FormatUint(uint64(id), 10))
	}
	return &user, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "storage/GetUserByEmail")
	defer span.Finish()

	query := `SELECT u.id, u.email, u.full_name, r.name AS role, u.password FROM users AS u
JOIN roles AS r ON u.role = r.id WHERE u.email = $1`
	rows, err := s.pool.Query(ctx, query, email)
	if err != nil {
		span.LogKV("error", "sql error")
		return nil, errors.Wrapf(err, "storage.getUserbyEmail user-email: [%s]", email)
	}
	var user models.User
	if err := pgxscan.ScanOne(&user, rows); err != nil {
		if pgxscan.NotFound(err) {
			return nil, errors.Wrapf(ErrUserNotExists, "storage.getUserbyEmail user-email: [%s]", email)
		}
		span.LogKV("error", "scanone error")
		return nil, errors.Wrapf(err, "storage.getUserbyEmail user-email: [%s]", email)
	}
	return &user, nil
}
