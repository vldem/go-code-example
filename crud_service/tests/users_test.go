//go:build integration
// +build integration

package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	pb "gitlab.ozon.dev/vldem/homework1/pkg/api"
	"gitlab.ozon.dev/vldem/homework1/tests/fixtures"
)

func TestUserCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		defer Db.TearDown()

		//act
		user := fixtures.User().Email("test01@dummy.com").Name("Bob Smith").Role("Admin").Password("123456").P()
		resp, err := BackendClient.UserCreate(context.Background(), &pb.BackendUserCreateRequest{
			Email:    user.Email,
			Name:     user.Name,
			Role:     user.Role,
			Password: user.Password,
		})

		//assert
		// log.Printf("error: [%v]", err)

		assert.Nil(t, err)
		assert.Equal(t, uint64(1), resp.GetId())

	})

	t.Run("error", func(t *testing.T) {
		//arrange
		Db.SetUp(t)
		defer Db.TearDown()

		//act
		user := fixtures.User().Email("test01dummy.com").Name("Bob Smith").Role("Admin").Password("123456").P()
		_, err := BackendClient.UserCreate(context.Background(), &pb.BackendUserCreateRequest{
			Email:    user.Email,
			Name:     user.Name,
			Role:     user.Role,
			Password: user.Password,
		})

		//assert
		//log.Printf("error: [%v]", err)

		require.EqualError(t, err, fmt.Sprintf("rpc error: code = InvalidArgument desc = bad email <%v>", user.Email))

	})
}
