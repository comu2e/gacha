package model
type CharacterUser struct {
	User_id         int `json:"user_id"`
	Character_id    int `json:"character_id"`
	Character_count int `json:"character_count"`
}