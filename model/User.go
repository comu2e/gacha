package model

type User struct {
	ID         int64 `json:"id"`
	Name       string `json:"name"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Phone      string `json:"phone"`
	UserStatus bool `json:"user_status"`
	XToken     string `json:"x_token"`
	//Character []Character
}
