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
		return nil, gqlError("Not A Member Of Team or Team doesn't exist", "code", "NOT_MEMBER_OF_TEAM")
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

//BoardValidateMember Validate Member
func BoardValidateMember(ctx context.Context, boardID int) (bool, error) {
	user := ForContext(ctx)
	if user == nil {
		fmt.Println("Not Logged In!")
		return false, gqlError("Not Logged In!", "code", "NOT_LOGGED_IN")
	}

	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var count int64

	if err := db.Table("board").Joins(
		"INNER JOIN team on board.team_id = team.id",
	).Joins(
		"INNER JOIN team_has_member on team_has_member.team_id = team.id",
	).Where("board.id = ? and team_has_member.user_id = ?", boardID, user.ID).Count(&count).Error; err != nil {
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

//BoardUpdateMultipleColumnsByID Update Columns
func BoardUpdateMultipleColumnsByID(ctx context.Context, id int, args []updateArgs) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}
	for _, val := range args {
		data[val.Key] = val.Value
	}

	if err := db.Table("board").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//BoardUpdateName Update Name
func BoardUpdateName(ctx context.Context, id int, name string) (string, error) {
	if stringIsEmpty(name) {
		return "Failed", gqlError("Invalid Name", "code", "INVALID_NAME")
	}

	var args []updateArgs
	args = append(args, updateArgs{
		Key:   "name",
		Value: name,
	})
	return BoardUpdateMultipleColumnsByID(ctx, id, args)
}

//BoardGetByID By ID
func BoardGetByID(ctx context.Context, id int) (*model.Board, error) {
	if access, err := BoardValidateMember(ctx, id); err != nil || !access {
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return nil, gqlError("(Not Member Of Team or Board doesn't exist", "code", "ACCESS_DENIED")
	}

	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var board model.Board
	if err := db.Table("board").Where("id = ?", id).Take(&board).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &board, nil
}

//BoardDeleteByID Delete By ID
func BoardDeleteByID(ctx context.Context, id int) (string, error) {
	if access, err := BoardValidateMember(ctx, id); err != nil || !access {
		if err != nil {
			fmt.Println(err)
			return "Failed", err
		}
		return "Failed", gqlError("(Not Member Of Team or Board doesn't exist", "code", "ACCESS_DENIED")
	}

	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if err := db.Exec(`
    DELETE b.*, l.*, li.* 
    FROM board as b
    INNER JOIN list as l on l.board_id = b.id
    INNER JOIN list_item as li on li.list_id = l.id
    WHERE b.id = ?;
    `, id).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}
