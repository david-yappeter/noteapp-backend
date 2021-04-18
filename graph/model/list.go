package model

import "time"

type List struct {
	ID        int        `json:"id" gorm:"type:int;not null;AUTO_INCREMENT"`
	Name      string     `json:"name" gorm:"type:text;not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	BoardID   int        `json:"board_id" gorm:"type:int;not null"`
}
