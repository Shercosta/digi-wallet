package models

type User struct {
	GormModel
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Level    int    `json:"level" gorm:"default:1;not null"`
	Balance  Balance `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}
