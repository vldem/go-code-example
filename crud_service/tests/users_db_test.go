//go:build integration
// +build integration

package tests

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	internalModel "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"
	repository "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/storage/postgres"
)

func TestCreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		defer Db.TearDown()

		userRepo := repository.New(Db.DB)

		// act
		res, err := userRepo.Add(context.Background(), internalModel.User{
			Email:    "test01@dummy.com",
			Name:     "Bob Smith",
			Role:     "Admin",
			Password: "123456",
		})
		// assert
		assert.NoError(t, err)
		assert.Equal(t, uint(1), res)
	})
}
