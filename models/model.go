package models

type User struct {
	UserID   int `gorm:"primaryKey;autoIncrement"`
	Name     string
	Username string
	Password string
	Status   string
}
