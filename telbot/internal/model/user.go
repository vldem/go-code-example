// This model contains user object and methods

package model

import (
	"crypto/md5"
	"fmt"

	"github.com/vldem/go-code-example/telbot/config"
)

var lastId uint

type User struct {
	id       uint
	email    string
	name     string
	role     uint8
	password string
}

func NewUser(email, name, role, password string) *User {
	u := User{}
	lastId++
	u.SetId(lastId)
	u.SetEmail(email)
	u.SetName(name)
	u.SetRole(role)
	u.SetPwd(password)

	return &u
}

func (u *User) SetId(id uint) {
	u.id = id
}

func (u *User) SetEmail(email string) {
	u.email = email
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) SetPwd(pwd string) {
	// use MD5 has to prevent storage of raw password in the storage
	// this is not enough secure approach, but it's better then nothing
	pwdHash := md5.Sum([]byte(pwd + config.Md5HashKey))
	u.password = fmt.Sprintf("%x", pwdHash)
}

func (u *User) SetRole(role string) {
	u.role = GetRoleId(role)
}

func (u User) String() string {
	return fmt.Sprintf("%d: %s / %s / %s / %s", u.id, u.email, u.name, GetRoleName(u.role), u.password)
}

func (u User) GetId() uint {
	return u.id
}

func (u User) GetName() string {
	return u.name
}

func (u User) GetEmail() string {
	return u.email
}

func (u User) GetPassword() string {
	return u.password
}

func (u User) GetRole() uint8 {
	return u.role
}

func (u User) GetRoleName() string {
	return GetRoleName(u.role)
}
