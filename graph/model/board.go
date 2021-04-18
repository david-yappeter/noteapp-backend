package model

import "time"

type Board struct {
	ID        int        `json:"id" gorm:"type:int;not null;AUTO_INCREMENT"`
	Name      string     `json:"name" gorm:"type:text;not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	TeamID    int        `json:"team_id" gorm:"type:int;not null"`
}
