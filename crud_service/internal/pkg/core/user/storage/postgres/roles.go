package postgres

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
)

func (s *Storage) GetRoleIdByName(ctx context.Context, roleName string) (uint8, error) {
	query := `SELECT id FROM roles WHERE name = $1`
	rows, err := s.pool.Query(ctx, query, roleName)
	if err != nil {
		return 0, errors.Wrapf(err, "storage.getRoleByName role: [%s]", roleName)
	}
	var roleId uint8
	if err := pgxscan.ScanOne(&roleId, rows); err != nil {
		return 0, errors.Wrapf(ErrRoleNotExists, "storage.getRoleByName role: [%s]", roleName)
	}
	return roleId, nil
}
