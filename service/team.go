package service

import (
	"context"
	"fmt"
	"myapp/config"
	"myapp/graph/model"
	"time"
)

//TeamCreate Create
func TeamCreate(ctx context.Context, name string) (*model.Team, error) {
	tokenUser := ForContext(ctx)

	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	team := model.Team{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		OwnerID:   tokenUser.ID,
	}

	if err := db.Table("team").Create(&team).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	if _, err := TeamHasMemberCreate(ctx, model.NewTeamHasMember{
		TeamID: team.ID,
		UserID: tokenUser.ID,
	}); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &team, nil
}

func TeamBatchMapByUserIds(ctx context.Context, userIds []int) (map[int][]*model.Team, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var tempModel []*struct {
		ID        int        `json:"id" gorm:"type:int;not null;AUTO_INCREMENT"`
		Name      string     `json:"name" gorm:"type:text;not null"`
		CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
		UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
		OwnerID   int        `json:"owner_id" gorm:"type:int;not null"`
		UserID    int
	}

	if err := db.Table("team").Select("team.*, team_has_member.user_id as user_id").Joins("INNER JOIN team_has_member on team.id = team_has_member.team_id").Where("team_has_member.user_id IN (?)", userIds).Find(&tempModel).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	var mappedObject = map[int][]*model.Team{}
	for _, val := range tempModel {
		mappedObject[val.UserID] = append(mappedObject[val.UserID], &model.Team{
			ID:        val.ID,
			Name:      val.Name,
			CreatedAt: val.CreatedAt,
			UpdatedAt: val.UpdatedAt,
			OwnerID:   val.OwnerID,
		})
	}

	return mappedObject, nil
}
