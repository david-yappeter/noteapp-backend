package tests

import (
	"context"
	"fmt"
	"myapp/graph/model"
	"time"
)

//TeamCreate Create
func (t *GormSuite) TeamCreate(ctx context.Context, userID int, name string) (*model.Team, error) {
	tokenUser := userID

	team := model.Team{
		Name:      name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		OwnerID:   tokenUser,
	}

	if err := t.tr.Table("team").Create(&team).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	if _, err := t.TeamHasMemberCreate(ctx, model.NewTeamHasMember{
		TeamID: team.ID,
		UserID: tokenUser,
	}); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &team, nil
}

//TeamDataloaderBatchByUserIds Dataloader
func (t *GormSuite) TeamDataLoaderBatchByUserIds(ctx context.Context, userIds []int) ([][]*model.Team, []error) {

	var tempModel []*struct {
		ID        int        `json:"id" gorm:"type:int;not null;AUTO_INCREMENT"`
		Name      string     `json:"name" gorm:"type:text;not null"`
		CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
		UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
		OwnerID   int        `json:"owner_id" gorm:"type:int;not null"`
		UserID    int
	}

	if err := t.tr.Table("team").Select("team.*, team_has_member.user_id as user_id").Joins("INNER JOIN team_has_member on team.id = team_has_member.team_id").Where("team_has_member.user_id IN (?)", userIds).Find(&tempModel).Error; err != nil {
		fmt.Println(err)
		return nil, []error{err}
	}

	itemById := map[int][]*model.Team{}
	for _, val := range tempModel {
		itemById[val.UserID] = append(itemById[val.UserID], &model.Team{
			ID:        val.ID,
			Name:      val.Name,
			CreatedAt: val.CreatedAt,
			UpdatedAt: val.UpdatedAt,
			OwnerID:   val.OwnerID,
		})
	}

	items := make([][]*model.Team, len(userIds))
	for i, id := range userIds {
		items[i] = itemById[id]
	}

	return items, nil
}

//TeamGetByID Get By ID
func (t *GormSuite) TeamGetByID(ctx context.Context, id int) (*model.Team, error) {

	var team model.Team
	if err := t.tr.Table("team").Where("id = ?", id).Take(&team).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &team, nil
}

func (t *GormSuite) TeamGetByIDAuthorize(ctx context.Context, userID int, id int) (*model.Team, error) {
	access, err := t.TeamValidateMember(ctx, userID, id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if !access {
		return nil, gqlError("Access Denied! (Not Member Of Team)", "code", "ACCESS_DENIED")
	}

	return t.TeamGetByID(ctx, id)
}

//TeamUpdateMultipleColumnsByID Update Multiple Columns
func (t *GormSuite) TeamUpdateMultipleColumnsByID(ctx context.Context, id int, args []updateArgs) (string, error) {

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}

	for _, val := range args {
		data[val.Key] = val.Value
	}

	if err := t.tr.Table("team").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//TeamUpdateName Update Name
func (t *GormSuite) TeamUpdateName(ctx context.Context, userID int, id int, name string) (*model.Team, error) {
	if stringIsEmpty(name) {
		return nil, gqlError("Invalid Name", "code", "INVALID_NAME")
	}

	access, err := t.TeamValidateMember(ctx, userID, id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if !access {
		return nil, gqlError("Access Denied! (Not Member Of Team)", "code", "ACCESS_DENIED")
	}

	var args []updateArgs
	args = append(args, updateArgs{
		Key:   "name",
		Value: name,
	})
	if _, err := t.TeamUpdateMultipleColumnsByID(ctx, id, args); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return t.TeamGetByID(ctx, id)
}
