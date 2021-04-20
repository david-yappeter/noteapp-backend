package service

import (
	"context"
	"fmt"
	"myapp/config"
	"myapp/graph/model"
	"myapp/tools"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/badoux/checkmail"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"gorm.io/gorm"
)

//UserCreate Create
func UserCreate(ctx context.Context, input model.NewUser) (*model.User, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	user := model.User{
		Name:      input.Name,
		Email:     strings.ToLower(input.Email),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		Avatar:    nil,
		Password:  tools.PasswordHash(input.Password),
	}

	if err := db.Table("user").Create(&user).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}

//UserUpdateSingleColumn Update Single Column
func UserUpdateMultipleColumnByUserID(ctx context.Context, args []updateArgs, userID int) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}
	for _, val := range args {
		data[val.Key] = val.Value
	}

	if err := db.Table("user").Where("id = ?", userID).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//UserUpdateName Update Name
func UserUpdateName(ctx context.Context, name string) (string, error) {
	tokenUser := ForContext(ctx)

	var args []updateArgs
	args = append(args, updateArgs{
		Key:   "name",
		Value: name,
	})

	if _, err := UserUpdateMultipleColumnByUserID(ctx, args, tokenUser.ID); err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//UserUpdateName Update Name
func UserUpdateAvatar(ctx context.Context, avatar *graphql.Upload) (*string, error) {
	tokenUser := ForContext(ctx)
	getUser, err := UserGetByID(ctx, tokenUser.ID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if getUser.Avatar != nil {
		GdriveDeleteFile(*getUser.Avatar)
	}

	var avatarFileID *string
	var args []updateArgs
	if avatar != nil {
		fileID, err := UploadFile(ctx, *avatar)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		args = append(args, updateArgs{
			Key:   "avatar",
			Value: fileID,
		})
		avatarFileID = &fileID
	} else {
		args = append(args, updateArgs{
			Key:   "avatar",
			Value: nil,
		})
		avatarFileID = nil
	}

	if _, err := UserUpdateMultipleColumnByUserID(ctx, args, tokenUser.ID); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return GdriveViewLink(avatarFileID), nil
}

//UserGetByID Get By ID
func UserGetByID(ctx context.Context, id int) (*model.User, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var user model.User

	if err := db.Table("user").Where("id = ?", id).Find(&user).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.Avatar = GdriveViewLink(user.Avatar)

	return &user, nil
}

//UserPaginationGetTotalData Pagination Total Data
func UserPaginationGetTotalData(ctx context.Context) (int, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var count int64

	if err := db.Table("user").Count(&count).Error; err != nil {
		fmt.Println(err)
		return 0, err
	}

	return int(count), nil
}

//UserPaginationGetTotalData Pagination Total Data
func UserPaginationGetNodes(ctx context.Context, limit *int, page *int, ascending *bool, sortBy *string) ([]*model.User, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var users []*model.User

	query := db.Table("user")
	tools.QueryMaker(query, limit, page, ascending, sortBy)

	if err := query.Find(&users).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	for index, val := range users {
		users[index].Avatar = GdriveViewLink(val.Avatar)
	}

	return users, nil
}

//UserGetByToken By Token
func UserGetByToken(ctx context.Context) (*model.User, error) {
	tokenUser := ForContext(ctx)
	if tokenUser == nil {
		return nil, &gqlerror.Error{
			Message: "Token Empty",
			Extensions: map[string]interface{}{
				"code": "TOKEN_EMPTY",
			},
		}
	}

	return UserGetByID(ctx, tokenUser.ID)
}

//UserGetByEmail Get By Email
func UserGetByEmail(ctx context.Context, email string) (*model.User, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var user model.User

	if err := db.Table("user").Where("lower(email) = ?", strings.ToLower(email)).Take(&user).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.Avatar = GdriveViewLink(user.Avatar)

	return &user, nil
}

func UserLogin(ctx context.Context, email string, password string) (*model.JwtToken, error) {
	if strings.EqualFold(email, "") {
		return nil, gqlError("Empty Email", "code", "EMPTY_EMAIL")
	}
	if strings.EqualFold(password, "") {
		return nil, gqlError("Empty Password", "code", "EMPTY_PASSWORD")
	}

	if err := checkmail.ValidateFormat(email); err != nil {
		fmt.Println(err)
		return nil, gqlError("Invalid Email", "code", "INVALID_EMAIL_FORMAT")
	}

	getUser, err := UserGetByEmail(ctx, email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gqlError("Email Not Found", "code", "EMAIL_NOT_FOUND")
		}
		return nil, err
	}

	if !tools.PasswordCompare(getUser.Password, password) {
		return nil, gqlError("Wrong Password!", "code", "WRONG_PASSWORD")
	}

	token, err := JwtTokenCreate(ctx, getUser.ID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &model.JwtToken{
		Type:  "Bearer",
		Token: token,
	}, nil
}

func UserRegister(ctx context.Context, input model.NewUser) (*model.JwtToken, error) {
	if strings.EqualFold(input.Email, "") {
		return nil, gqlError("Empty Email", "code", "EMPTY_EMAIL")
	}
	if strings.EqualFold(input.Password, "") {
		return nil, gqlError("Empty Password", "code", "EMPTY_PASSWORD")
	}

	if err := checkmail.ValidateFormat(input.Email); err != nil {
		fmt.Println(err)
		return nil, gqlError("Invalid Email", "code", "INVALID_EMAIL_FORMAT")
	}

	if input.Password != input.ConfirmPassword {
		return nil, gqlError("Password And Confirm Password Not Match!", "code", "INVALID_PASSWORD_MATCH")
	}

	getUser, err := UserGetByEmail(ctx, input.Email)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	if getUser != nil {
		return nil, gqlError("Email Already Exist", "code", "EMAIL_EXIST")
	}

	if getUser, err = UserCreate(ctx, input); err != nil {
		fmt.Println(err)
		return nil, err
	}

	token, err := JwtTokenCreate(ctx, getUser.ID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &model.JwtToken{
		Type:  "Bearer",
		Token: token,
	}, nil
}

//UserDataloaderBatchByTeamIds Dataloader
func UserDataloaderBatchByTeamIds(ctx context.Context, teamIds []int) ([][]*model.User, []error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var tempModel []*struct {
		ID        int        `json:"id" gorm:"type:int;not null;AUTO_INCREMENT"`
		Name      string     `json:"name" gorm:"type:text;not null"`
		Email     string     `json:"email" gorm:"type:text;not null"`
		Password  string     `json:"password" gorm:"type:text;not null"`
		CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp;not null"`
		UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
		Avatar    *string    `json:"avatar" gorm:"type:text;null;default:NULL"`
		TeamID    int        `json:"team_id"`
	}

	if err := db.Table("user").Select("user.*, team_has_member.team_id AS team_id").Joins("INNER JOIN team_has_member on user.id = team_has_member.user_id").Where("team_has_member.team_id IN (?)", teamIds).Find(&tempModel).Error; err != nil {
		fmt.Println(err)
		return nil, []error{err}
	}

	for index, val := range tempModel {
		tempModel[index].Avatar = GdriveViewLink(val.Avatar)
	}

	itemById := map[int][]*model.User{}
	for _, val := range tempModel {
		itemById[val.TeamID] = append(itemById[val.TeamID], &model.User{
			ID:        val.ID,
			Name:      val.Name,
			CreatedAt: val.CreatedAt,
			UpdatedAt: val.UpdatedAt,
			Email:     val.Email,
			Password:  val.Password,
			Avatar:    val.Avatar,
		})
	}

	items := make([][]*model.User, len(teamIds))
	for i, id := range teamIds {
		items[i] = itemById[id]
	}

	return items, nil
}

//UserEditPassword Edit Password
func UserEditPassword(ctx context.Context, newPassword string) (string, error) {
	if stringIsEmpty(newPassword) {
		return "Failed", gqlError("Invalid New Password", "code", "INVALID_NEW_PASSWORD")
	}

	var args []updateArgs
	args = append(args, updateArgs{
		Key:   "password",
		Value: tools.PasswordHash(newPassword),
	})
	return UserUpdateMultipleColumnByUserID(ctx, args, ForContext(ctx).ID)
}
