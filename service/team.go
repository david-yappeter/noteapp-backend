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

//TeamDataloaderBatchByUserIds Dataloader
func TeamDataLoaderBatchByUserIds(ctx context.Context, userIds []int) ([][]*model.Team, []error) {
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
func TeamGetByID(ctx context.Context, id int) (*model.Team, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var team model.Team
	if err := db.Table("team").Where("id = ?", id).Take(&team).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &team, nil
}

func TeamGetByIDAuthorize(ctx context.Context, id int) (*model.Team, error) {
	access, err := TeamValidateMember(ctx, id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if !access {
		return nil, gqlError("Access Denied! (Not Member Of Team)", "code", "ACCESS_DENIED")
	}

	return TeamGetByID(ctx, id)
}

//TeamUpdateMultipleColumnsByID Update Multiple Columns
func TeamUpdateMultipleColumnsByID(ctx context.Context, id int, args []updateArgs) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}

	for _, val := range args {
		data[val.Key] = val.Value
	}

	if err := db.Table("team").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//TeamUpdateName Update Name
func TeamUpdateName(ctx context.Context, id int, name string) (*model.Team, error) {
	if stringIsEmpty(name) {
		return nil, gqlError("Invalid Name", "code", "INVALID_NAME")
	}

	access, err := TeamValidateMember(ctx, id)
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
	if _, err := TeamUpdateMultipleColumnsByID(ctx, id, args); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return TeamGetByID(ctx, id)
}

//BoardDeleteByID Delete By ID
func TeamDeleteByID(ctx context.Context, id int) (string, error) {
	if access, err := TeamValidateMember(ctx, id); err != nil || !access {
		if err != nil {
			fmt.Println(err)
			return "Failed", err
		}
		return "Failed", gqlError("(Not Member Of Team or Team doesn't exist", "code", "ACCESS_DENIED")
	}

	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if err := db.Exec(`
    DELETE t.*, b.*, l.*, li.* 
    FROM team as t
    INNER JOIN board as b on b.team_id = t.id
    INNER JOIN list as l on l.board_id = b.id
    INNER JOIN list_item as li on li.list_id = l.id
    WHERE t.id = ?;
    `, id).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}
