package service

import (
	"context"
	"fmt"
	"myapp/config"
	"myapp/graph/model"
	"time"

	"gorm.io/gorm"
)

//BoardCreate Create
func BoardCreate(ctx context.Context, input model.NewBoard) (*model.Board, error) {
	if access, err := BoardCheckUserAccess(ctx, input.TeamID); err != nil || !access {
		if err != nil {
			return nil, err
		}
		return nil, gqlError("Not A Member Of Team", "code", "NOT_MEMBER_OF_TEAM")
	}

	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	board := model.Board{
		Name:      input.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		TeamID:    input.TeamID,
	}

	if err := db.Table("board").Create(&board).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &board, nil
}

//BoardCheckUserAccess Check User Access
func BoardCheckUserAccess(ctx context.Context, teamID int) (bool, error) {
	tokenUser := ForContext(ctx)
	if tokenUser == nil {
		return false, gqlError("Not Logged In!", "code", "NOT_LOGGED_IN")
	}

	if _, err := TeamHasMemberGetByUserIDAndTeamID(ctx, tokenUser.ID, teamID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
