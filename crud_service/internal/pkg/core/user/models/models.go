package models

type User struct {
	Id       uint   `db:"id"`
	Email    string `db:"email"`
	Name     string `db:"full_name"`
	Role     string `db:"role"`
	Password string `db:"password"`
}

type SortingOrder struct {
	Field      string
	Descending bool
}

type SortingField map[string]string

var sortingFields SortingField

func init() {
	sortingFields = SortingField{}
	sortingFields["id"] = "u.id"
	sortingFields["email"] = "u.email"
	sortingFields["name"] = "u.full_name"
}

func GetSortingFieldName(name string) string {
	return sortingFields[name]
}
