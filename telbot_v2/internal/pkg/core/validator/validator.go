// This package contains validators for user inputs in order to prevent saving of suspicious data in data storage

package validator

import (
	"fmt"
	"regexp"

	"github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/models"

	"github.com/pkg/errors"
)

var ErrInvalidData = errors.New("invalid data")

type ValidationHandler func(string) error

type Validator struct {
	rules map[string]ValidationHandler
}

var validator Validator

func (v *Validator) RegisterHandler(validatorName string, f ValidationHandler) {
	v.rules[validatorName] = f
}

func init() {
	validator = Validator{make(map[string]ValidationHandler)}
	validator.RegisterHandler("email", ValidateEmail)
	validator.RegisterHandler("name", ValidateName)
	validator.RegisterHandler("role", ValidateRole)
	validator.RegisterHandler("password", ValidatePassword)
}

func MakeParametersToValidate(params []string) map[string]string {
	paramList := make(map[string]string, len(params))
	paramList["email"] = params[0]
	paramList["name"] = params[1]
	paramList["role"] = params[2]
	paramList["password"] = params[3]
	return paramList
}

func ValidateEmail(email string) error {

	matched, err := regexp.MatchString(`^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*\.\w{2,4}$`, email)

	if err != nil {
		return errors.Wrap(err, "email validator")
	}

	if !matched {
		return fmt.Errorf("bad email <%v>", email)
	}

	return nil
}

func ValidatePassword(pwd string) error {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9]{6,15}$`, pwd)

	if err != nil {
		return errors.Wrap(err, "password validator")
	}

	if !matched {
		return fmt.Errorf("bad password <%v> (should contain <a-zA-Z0-9> length 6-15 symbols) ", pwd)
	}

	return nil
}

func ValidateName(name string) error {
	matched, err := regexp.MatchString(`^[a-zA-Z0-9 \-\_]{2,30}$`, name)

	if err != nil {
		return errors.Wrap(err, "name validator")
	}

	if !matched {
		return fmt.Errorf("bad name <%v> (should contain <a-zA-Z0-9 -_> length 2-30 symbols)", name)
	}

	return nil
}

func ValidateRole(role string) error {
	if models.GetRoleId(role) == 0 {
		return fmt.Errorf("bad role <%v>", role)
	}

	return nil
}

func ValidateParameters(data map[string]string) error {
	if len(data) == 0 {
		return fmt.Errorf("empty parameters to validate")
	}

	for key, val := range data {
		if validatorFunc, ok := validator.rules[key]; ok {
			if err := validatorFunc(val); err != nil {
				return err
			}
		}
	}

	return nil
}
