package auth

import (
	"crypto/md5"
	"fmt"

	"github.com/vldem/go-code-example/telbot_v2/internal/config"
	"github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/models"
)

func VerifyPassword(user models.User, pwd string) error {
	pwdHash := GenHashPassword(pwd)
	if user.Password != pwdHash {
		return fmt.Errorf("wrong old password")
	}
	return nil
}

func GenHashPassword(password string) string {
	// use MD5 has to prevent storage of raw password in the storage
	// this is not enough secure approach, but it's better then nothing
	pwdHash := md5.Sum([]byte(password + config.Md5HashKey))
	return fmt.Sprintf("%x", pwdHash)
}
