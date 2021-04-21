package tests

import (
	"context"
	"fmt"
	"myapp/graph/model"

	"gorm.io/gorm"
)

//TeamHasMemberCreate Create
func (t *GormSuite) TeamHasMemberCreate(ctx context.Context, input model.NewTeamHasMember) (*model.TeamHasMember, error) {
	teamHasMember := model.TeamHasMember{
		TeamID: input.TeamID,
		UserID: input.UserID,
	}

	if err := t.tr.Table("team_has_member").Create(&teamHasMember).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &teamHasMember, nil
}

//TeamAddMember
func (t *GormSuite) TeamAddMember(ctx context.Context, userID int, input model.NewTeamHasMember) (*model.TeamHasMember, error) {
	if access, err := t.TeamValidateMember(ctx, userID, input.TeamID); err != nil || !access {
		if err != nil {
			return nil, err
		}
		return nil, gqlError("Access Denied! (Not Member Of Team)", "code", "ACCESS_DENIED")
	}

	if getTeamHasMember, err := t.TeamHasMemberGetByTeamIDAndUserID(ctx, input.TeamID, input.UserID); err != nil || getTeamHasMember != nil {
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println(err)
			return nil, err
		}

		if getTeamHasMember != nil {
			return nil, gqlError("Already Part of Team Member", "code", "NOTHING_CHANGED")
		}
	}

	return t.TeamHasMemberCreate(ctx, input)
}

//TeamRemoveMember Remove Member
func (t *GormSuite) TeamRemoveMember(ctx context.Context, userID int, input model.NewTeamHasMember) (string, error) {
	if access, err := t.TeamValidateMember(ctx, userID, input.TeamID); err != nil || !access {
		if err != nil {
			return "Failed", err
		}
		return "Failed", gqlError("Access Denied! (Not Member Of Team)", "code", "ACCESS_DENIED")
	}

	return t.TeamHasMemberDelete(ctx, input)
}

//TeamHasMemberDelete Delete
func (t *GormSuite) TeamHasMemberDelete(ctx context.Context, input model.NewTeamHasMember) (string, error) {

	if err := t.tr.Table("team_has_member").Where("team_id = ? AND user_id = ?", input.TeamID, input.UserID).Delete(&model.TeamHasMember{}).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//TeamHasMemberGetByTeamIDAndUserID Get By Team ID and User ID
func (t *GormSuite) TeamHasMemberGetByTeamIDAndUserID(ctx context.Context, teamID int, userID int) (*model.TeamHasMember, error) {

	var teamHasMember model.TeamHasMember

	if err := t.tr.Table("team_has_member").Where("user_id = ? AND team_id = ?", userID, teamID).Take(&teamHasMember).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &teamHasMember, nil
}

//TeamValidateMember Validate Member
func (t *GormSuite) TeamValidateMember(ctx context.Context, userID int, teamID int) (bool, error) {
	user := userID

	var count int64

	if err := t.tr.Table("team_has_member").Where("user_id = ? AND team_id = ?", user, teamID).Count(&count).Error; err != nil {
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
