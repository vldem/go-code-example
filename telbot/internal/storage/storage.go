// This is a storage that contains list of users
package storage

import (
	"log"
	"strconv"

	"github.com/vldem/go-code-example/telbot/internal/model"

	"github.com/pkg/errors"
)

var data map[uint]*model.User

var UserNotExists = errors.New("user does not exists")
var UserExists = errors.New("user already exists")

func init() {
	log.Println("init storage")

	data = make(map[uint]*model.User)
	u := model.NewUser("test@dummy.com", "Bob Smith", "Admin", "123456")
	if err := AddUser(u); err != nil {
		log.Panic(err)
	}
}

func ListUser() []*model.User {
	res := make([]*model.User, 0, len(data))
	for _, v := range data {
		res = append(res, v)
	}
	return res
}

func AddUser(u *model.User) error {
	if _, ok := data[u.GetId()]; ok {
		return errors.Wrap(UserExists, strconv.FormatUint(uint64(u.GetId()), 10))
	}
	data[u.GetId()] = u
	return nil
}

func UpdateUser(u *model.User) error {
	if _, ok := data[u.GetId()]; !ok {
		return errors.Wrap(UserNotExists, strconv.FormatUint(uint64(u.GetId()), 10))
	}
	data[u.GetId()] = u
	return nil
}

func DeleteUser(id uint) error {
	if _, ok := data[id]; !ok {
		return errors.Wrap(UserNotExists, strconv.FormatUint(uint64(id), 10))
	}
	delete(data, id)
	return nil
}

func GetUser(id uint) (*model.User, error) {
	if _, ok := data[id]; !ok {
		return nil, errors.Wrap(UserNotExists, strconv.FormatUint(uint64(id), 10))
	}
	return data[id], nil
}
