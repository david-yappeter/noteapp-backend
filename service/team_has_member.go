package service

import (
	"context"
	"fmt"
	"myapp/config"
	"myapp/graph/model"
)

//TeamHasMemberCreate Create
func TeamHasMemberCreate(ctx context.Context, input model.NewTeamHasMember) (*model.TeamHasMember, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	teamHasMember := model.TeamHasMember{
		TeamID: input.TeamID,
		UserID: input.UserID,
	}

	if err := db.Table("team_has_member").Create(&teamHasMember).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &teamHasMember, nil
}

func TeamHasMemberGetByUserIDAndTeamID(ctx context.Context, userID int, teamID int) (*model.TeamHasMember, error) {
    db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var teamHasMember model.TeamHasMember

	if err := db.Table("team_has_member").Where("user_id = ? AND team_id = ?", userID, teamID).Take(&teamHasMember).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &teamHasMember, nil
}