package model

type Character struct {
	ID            int64  `gorm:"primaryKey"`
	CharacterName string
}

