package model

type User struct {
	ID         int64
	Username   string
	Firstname  string
	Lastname   string
	Email      string
	Password   string
	Phone      string
	UserStatus string
	Character []Character
}
