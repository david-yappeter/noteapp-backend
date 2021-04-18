package model

import "time"

type User struct {
	ID        int        `json:"id" gorm:"type:int;not null;AUTO_INCREMENT"`
	Name      string     `json:"name" gorm:"type:text;not null"`
	Email     string     `json:"email" gorm:"type:text;not null"`
	Password  string     `json:"password" gorm:"type:text;not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Avatar    *string    `json:"avatar" gorm:"type:text;null;default:NULL"`
}
