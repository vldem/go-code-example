//go:build integration
// +build integration

package fixtures

import "gitlab.ozon.dev/vldem/homework1/internal/pkg/core/user/models"

type UserBuilder struct {
	instance *models.User
}

func User() *UserBuilder {
	return &UserBuilder{
		instance: &models.User{},
	}
}

func (b *UserBuilder) Id(v uint) *UserBuilder {
	b.instance.Id = v
	return b
}

func (b *UserBuilder) Email(v string) *UserBuilder {
	b.instance.Email = v
	return b
}

func (b *UserBuilder) Name(v string) *UserBuilder {
	b.instance.Name = v
	return b
}

func (b *UserBuilder) Role(v string) *UserBuilder {
	b.instance.Role = v
	return b
}

func (b *UserBuilder) Password(v string) *UserBuilder {
	b.instance.Password = v
	return b
}

func (b *UserBuilder) P() *models.User {
	return b.instance
}

func (b *UserBuilder) V() models.User {
	return *b.instance
}
