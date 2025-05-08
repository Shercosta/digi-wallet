package models

type User struct {
	GormModel
	Username string `json:"username"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "users"
}
