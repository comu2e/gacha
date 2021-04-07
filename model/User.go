package model

type User struct {
	ID         int64  `gorm:"primaryKey"`
	Username   string `gorm:"unique"`
	Firstname  string
	Lastname   string
	Email      string
	Password   string
	Phone      string
	UserStatus string
	Character []Character
}
