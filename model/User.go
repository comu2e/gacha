package model

type User struct {
	ID         int64
	Name       string
	Firstname  string
	Lastname   string
	Email      string
	Password   string
	Phone      string
	UserStatus bool
	XToken     string
	//Character []Character
}
