package models

type User struct {
	Id       uint
	Email    string
	Name     string
	Role     uint8
	Password string
}
