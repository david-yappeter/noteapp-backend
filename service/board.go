package service

import (
	"context"
	"fmt"
	"myapp/config"
	"myapp/graph/model"
	"time"
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
	access, err := TeamValidateMember(ctx, teamID)
	if err != nil {
		return false, err
	}

	return access, nil
}

//BoardDataloaderBatchByTeamIds Dataloader
func BoardDataloaderBatchByTeamIds(ctx context.Context, teamIds []int) ([][]*model.Board, []error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var boards []*model.Board
	if err := db.Table("board").Where("team_id IN (?)", teamIds).Find(&boards).Error; err != nil {
		fmt.Println(err)
		return nil, []error{err}
	}

	itemById := map[int][]*model.Board{}
	for _, val := range boards {
		itemById[val.TeamID] = append(itemById[val.TeamID], val)
	}

	items := make([][]*model.Board, len(teamIds))
	for i, id := range teamIds {
		items[i] = itemById[id]
	}

	return items, nil
}
