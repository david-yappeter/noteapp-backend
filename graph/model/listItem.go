package model

import "time"

type ListItem struct {
	ID        int        `json:"id" gorm:"type:int;not null;AUTO_INCREMENT"`
	Name      string     `json:"name" gorm:"type:text;not null"`
	ListID    *int       `json:"list_id" gorm:"type:int;null;default:NULL"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Next      *int       `json:"next" gorm:"type:int;null;default:NULL"`
	Prev      *int       `json:"prev" gorm:"type:int;null;default:NULL"`
}
