package service

import (
	"context"
	"fmt"
	"myapp/config"
	"myapp/graph/model"
	"strings"

	"gorm.io/gorm"
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

//TeamAddMember
func TeamAddMember(ctx context.Context, input model.NewTeamHasMember) (*model.TeamHasMember, error) {
	if access, err := TeamValidateMember(ctx, input.TeamID); err != nil || !access {
		if err != nil {
			return nil, err
		}
		return nil, gqlError("Access Denied! (Not Member Of Team)", "code", "ACCESS_DENIED")
	}

	if getTeamHasMember, err := TeamHasMemberGetByTeamIDAndUserID(ctx, input.TeamID, input.UserID); err != nil || getTeamHasMember != nil {
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println(err)
			return nil, err
		}

		if getTeamHasMember != nil {
			return nil, gqlError("Already Part of Team Member", "code", "NOTHING_CHANGED")
		}
	}

	return TeamHasMemberCreate(ctx, input)
}

func TeamAddMemberByEmail(ctx context.Context, input model.NewTeamHasMemberByEmail) (*model.TeamHasMember, error) {
	if access, err := TeamValidateMember(ctx, input.TeamID); err != nil || !access {
		if err != nil {
			return nil, err
		}
		return nil, gqlError("Access Denied! (Not Member Of Team)", "code", "ACCESS_DENIED")
	}

	if getTeamHasMember, err := TeamHasMemberGetByTeamIDAndEmail(ctx, input.TeamID, input.Email); err != nil || getTeamHasMember != nil {
		if err != nil && err != gorm.ErrRecordNotFound {
			fmt.Println(err)
			return nil, err
		}

		if getTeamHasMember != nil {
			return nil, gqlError("Already Part of Team Member", "code", "NOTHING_CHANGED")
		}
	}

	user, err := UserGetByEmail(ctx, input.Email)
	if err != nil {
		fmt.Println(err)
		if err == gorm.ErrRecordNotFound {
			return nil, gqlError("Email Not Found", "code", "EMAIL_NOT_FOUND")
		}
		return nil, err
	}

	return TeamHasMemberCreate(ctx, model.NewTeamHasMember{
		TeamID: input.TeamID,
		UserID: user.ID,
	})
}

func TeamHasMemberGetByTeamIDAndEmail(ctx context.Context, teamID int, email string) (*model.TeamHasMember, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var teamHasMember model.TeamHasMember
	if err := db.Table("team_has_member").Select("team_has_member.*").Joins(`INNER JOIN "user" on "user".id = team_has_member.user_id`).Where(`"user".email = ?`, strings.ToLower(email)).Take(&teamHasMember).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &teamHasMember, nil
}

//TeamRemoveMember Remove Member
func TeamRemoveMember(ctx context.Context, input model.NewTeamHasMember) (string, error) {
	if access, err := TeamValidateMember(ctx, input.TeamID); err != nil || !access {
		if err != nil {
			return "Failed", err
		}
		return "Failed", gqlError("Access Denied! (Not Member Of Team)", "code", "ACCESS_DENIED")
	}

	return TeamHasMemberDelete(ctx, input)
}

//TeamHasMemberDelete Delete
func TeamHasMemberDelete(ctx context.Context, input model.NewTeamHasMember) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if err := db.Table("team_has_member").Where("team_id = ? AND user_id = ?", input.TeamID, input.UserID).Delete(&model.TeamHasMember{}).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//TeamHasMemberGetByTeamIDAndUserID Get By Team ID and User ID
func TeamHasMemberGetByTeamIDAndUserID(ctx context.Context, teamID int, userID int) (*model.TeamHasMember, error) {
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
