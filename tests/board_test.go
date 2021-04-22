package tests

import (
	"context"
	"fmt"
	"myapp/graph/model"
	"time"
)

//BoardCreate Create
func (t *GormSuite) BoardCreate(ctx context.Context, userID int, input model.NewBoard) (*model.Board, error) {
	if access, err := t.BoardCheckUserAccess(ctx, userID, input.TeamID); err != nil || !access {
		if err != nil {
			return nil, err
		}
		return nil, gqlError("Not A Member Of Team or Team doesn't exist", "code", "NOT_MEMBER_OF_TEAM")
	}

	board := model.Board{
		Name:      input.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		TeamID:    input.TeamID,
	}

	if err := t.tr.Table("board").Create(&board).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &board, nil
}

//BoardCheckUserAccess Check User Access
func (t *GormSuite) BoardCheckUserAccess(ctx context.Context, userID int, teamID int) (bool, error) {
	access, err := t.TeamValidateMember(ctx, userID, teamID)
	if err != nil {
		return false, err
	}

	return access, nil
}

//BoardDataloaderBatchByTeamIds Dataloader
func (t *GormSuite) BoardDataloaderBatchByTeamIds(ctx context.Context, teamIds []int) ([][]*model.Board, []error) {

	var boards []*model.Board
	if err := t.tr.Table("board").Where("team_id IN (?)", teamIds).Find(&boards).Error; err != nil {
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
func (t *GormSuite) BoardValidateMember(ctx context.Context, userID int, boardID int) (bool, error) {
	user := userID

	var count int64

	if err := t.tr.Table("board").Joins(
		"INNER JOIN team on board.team_id = team.id",
	).Joins(
		"INNER JOIN team_has_member on team_has_member.team_id = team.id",
	).Where("board.id = ? and team_has_member.user_id = ?", boardID, user).Count(&count).Error; err != nil {
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
func (t *GormSuite) BoardUpdateMultipleColumnsByID(ctx context.Context, id int, args []updateArgs) (string, error) {
	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}
	for _, val := range args {
		data[val.Key] = val.Value
	}

	if err := t.tr.Table("board").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//BoardUpdateName Update Name
func (t *GormSuite) BoardUpdateName(ctx context.Context, id int, name string) (string, error) {
	if stringIsEmpty(name) {
		return "Failed", gqlError("Invalid Name", "code", "INVALID_NAME")
	}

	var args []updateArgs
	args = append(args, updateArgs{
		Key:   "name",
		Value: name,
	})
	return t.BoardUpdateMultipleColumnsByID(ctx, id, args)
}

func (t *GormSuite) BoardGetByID(ctx context.Context, userID int, id int) (*model.Board, error) {
	if access, err := t.BoardValidateMember(ctx, userID, id); err != nil || !access {
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return nil, gqlError("(Not Member Of Team or Board doesn't exist", "code", "ACCESS_DENIED")
	}

	var board model.Board
	if err := t.tr.Table("board").Where("id = ?", id).Take(&board).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &board, nil
}
