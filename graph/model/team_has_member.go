package model

type TeamHasMember struct {
	ID     int `json:"id" gorm:"type:int;not null;AUTO_INCREMENT"`
	TeamID int `json:"team_id" gorm:"type:int;not null"`
	UserID int `json:"user_id" gorm:"type:int;not null"`
}
