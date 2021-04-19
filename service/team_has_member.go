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

//TeamValidateMember Validate Member
func TeamValidateMember(ctx context.Context, teamID int) (bool, error) {
	user := ForContext(ctx)
	if user == nil {
		fmt.Println("Not Logged In!")
		return false, gqlError("Not Logged In!", "code", "NOT_LOGGED_IN")
	}

	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var count int64

	if err := db.Table("team_has_member").Where("user_id = ? AND team_id = ?", user.ID, teamID).Count(&count).Error; err != nil {
		fmt.Println(err)
		return false, err
	}

	if count == 0 {
		return false, nil
	} else if count == 1 {
		return true, nil
	}

	fmt.Println("Unhandled Data")
	return false, gqlError("Unhandled Case", "code", "UNHANDLED_CASE")
}
