package models

type User struct {
	GormModel
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
}

func (User) TableName() string {
	return "users"
}
