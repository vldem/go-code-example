package auth

import (
	"crypto/md5"
	"fmt"

	"github.com/vldem/go-code-example/telbot/config"
	"github.com/vldem/go-code-example/telbot/internal/storage"

	"github.com/pkg/errors"
)

func VerifyPassword(id uint, pwd string) error {
	user, err := storage.GetUser(id)
	if err != nil {
		return errors.Wrap(err, "verify password")
	}

	pwdHash := fmt.Sprintf("%x", md5.Sum([]byte(pwd+config.Md5HashKey)))
	if (*user).GetPassword() != pwdHash {
		return fmt.Errorf("wrong old password")
	}

	return nil
}
