// This model contains user object and methods

package model

import (
	"github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/models"
	storagePkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/storage"
	localStoragePkg "github.com/vldem/go-code-example/telbot_v2/internal/pkg/core/user/storage/local"
)

type Interface interface {
	Create(user models.User) error
	Update(user models.User) error
	Delete(id uint) error
	Get(id uint) (*models.User, error)
	List() []models.User
}

type core struct {
	storage storagePkg.Interface
}

func New() Interface {
	return &core{
		storage: localStoragePkg.New(),
	}
}

func (c *core) Create(user models.User) error {
	return c.storage.Add(user)
}

func (c *core) Update(user models.User) error {
	return c.storage.Update(user)
}

func (c *core) Delete(id uint) error {
	return c.storage.Delete(id)
}

func (c *core) Get(id uint) (*models.User, error) {
	result, err := c.storage.Get(id)
	return result, err
}

func (c *core) List() []models.User {
	return c.storage.List()
}

// func NewUser(email, name, role, password string) *User {
// 	u := User{}
// 	lastId++
// 	u.SetId(lastId)
// 	u.SetEmail(email)
// 	u.SetName(name)
// 	u.SetRole(role)
// 	u.SetPwd(password)

// 	return &u
// }

// func (u *User) SetId(id uint) {
// 	u.id = id
// }

// func (u *User) SetEmail(email string) {
// 	u.email = email
// }

// func (u *User) SetName(name string) {
// 	u.name = name
// }

// func (u *User) SetPwd(pwd string) {
// 	// use MD5 has to prevent storage of raw password in the storage
// 	// this is not enough secure approach, but it's better then nothing
// 	pwdHash := md5.Sum([]byte(pwd + config.Md5HashKey))
// 	u.password = fmt.Sprintf("%x", pwdHash)
// }

// func (u *User) SetRole(role string) {
// 	u.role = GetRoleId(role)
// }

// func (u User) String() string {
// 	return fmt.Sprintf("%d: %s / %s / %s / %s", u.id, u.email, u.name, GetRoleName(u.role), u.password)
// }

// func (u User) GetId() uint {
// 	return u.id
// }

// func (u User) GetName() string {
// 	return u.name
// }

// func (u User) GetEmail() string {
// 	return u.email
// }

// func (u User) GetPassword() string {
// 	return u.password
// }

// func (u User) GetRole() uint8 {
// 	return u.role
// }

// func (u User) GetRoleName() string {
// 	return GetRoleName(u.role)
// }
