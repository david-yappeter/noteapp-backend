package tests

import (
	"context"
	"fmt"
	"myapp/graph/model"
	"myapp/service"
	"myapp/tools"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/badoux/checkmail"
	"gorm.io/gorm"
)

//UserCreate Create
func (t *GormSuite) UserCreate(ctx context.Context, input model.NewUser) (*model.User, error) {

	user := model.User{
		Name:      input.Name,
		Email:     strings.ToLower(input.Email),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		Avatar:    nil,
		Password:  tools.PasswordHash(input.Password),
	}

	if err := t.tr.Table("user").Create(&user).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &user, nil
}

//UserUpdateSingleColumn Update Single Column
func (t *GormSuite) UserUpdateMultipleColumnByUserID(ctx context.Context, args []updateArgs, userID int) (string, error) {
	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}
	for _, val := range args {
		data[val.Key] = val.Value
	}

	if err := t.tr.Table("user").Where("id = ?", userID).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//UserUpdateName Update Name
func (t *GormSuite) UserUpdateName(ctx context.Context, userID int, name string) (string, error) {
	tokenUser := userID

	var args []updateArgs
	args = append(args, updateArgs{
		Key:   "name",
		Value: name,
	})

	if _, err := t.UserUpdateMultipleColumnByUserID(ctx, args, tokenUser); err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//UserUpdateName Update Name
func (t *GormSuite) UserUpdateAvatar(ctx context.Context, userID int, avatar *graphql.Upload) (*string, error) {
	tokenUser := userID

	var avatarFileID *string
	var args []updateArgs
	if avatar != nil {
		fileID := "test"
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

	if _, err := t.UserUpdateMultipleColumnByUserID(ctx, args, tokenUser); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return service.GdriveViewLink(avatarFileID), nil
}

//UserGetByID Get By ID
func (t *GormSuite) UserGetByID(ctx context.Context, id int) (*model.User, error) {

	var user model.User

	if err := t.tr.Table("user").Where("id = ?", id).Find(&user).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.Avatar = service.GdriveViewLink(user.Avatar)

	return &user, nil
}

//UserPaginationGetTotalData Pagination Total Data
func (t *GormSuite) UserPaginationGetTotalData(ctx context.Context) (int, error) {

	var count int64

	if err := t.tr.Table("user").Count(&count).Error; err != nil {
		fmt.Println(err)
		return 0, err
	}

	return int(count), nil
}

//UserPaginationGetTotalData Pagination Total Data
func (t *GormSuite) UserPaginationGetNodes(ctx context.Context, limit *int, page *int, ascending *bool, sortBy *string) ([]*model.User, error) {

	var users []*model.User

	query := t.tr.Table("user")
	tools.QueryMaker(query, limit, page, ascending, sortBy)

	if err := query.Find(&users).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	for index, val := range users {
		users[index].Avatar = service.GdriveViewLink(val.Avatar)
	}

	return users, nil
}

//UserGetByToken By Token
func (t *GormSuite) UserGetByToken(ctx context.Context, userID int) (*model.User, error) {
	tokenUser := userID

	return t.UserGetByID(ctx, tokenUser)
}

//UserGetByEmail Get By Email
func (t *GormSuite) UserGetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := t.tr.Table("user").Where("lower(email) = ?", strings.ToLower(email)).Take(&user).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	user.Avatar = service.GdriveViewLink(user.Avatar)

	return &user, nil
}

func (t *GormSuite) UserLogin(ctx context.Context, email string, password string) (*model.JwtToken, error) {
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

	getUser, err := t.UserGetByEmail(ctx, email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gqlError("Email Not Found", "code", "EMAIL_NOT_FOUND")
		}
		return nil, err
	}

	if !tools.PasswordCompare(getUser.Password, password) {
		return nil, gqlError("Wrong Password!", "code", "WRONG_PASSWORD")
	}

	token, err := service.JwtTokenCreate(ctx, getUser.ID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &model.JwtToken{
		Type:  "Bearer",
		Token: token,
	}, nil
}

func (t *GormSuite) UserRegister(ctx context.Context, input model.NewUser) (*model.JwtToken, error) {
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

	getUser, err := t.UserGetByEmail(ctx, input.Email)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}

	if getUser != nil {
		return nil, gqlError("Email Already Exist", "code", "EMAIL_EXIST")
	}

	if getUser, err = t.UserCreate(ctx, input); err != nil {
		fmt.Println(err)
		return nil, err
	}

	token, err := service.JwtTokenCreate(ctx, getUser.ID)
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
func (t *GormSuite) UserDataloaderBatchByTeamIds(ctx context.Context, teamIds []int) ([][]*model.User, []error) {
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

	if err := t.tr.Table("user").Select("user.*, team_has_member.team_id AS team_id").Joins("INNER JOIN team_has_member on user.id = team_has_member.user_id").Where("team_has_member.team_id IN (?)", teamIds).Find(&tempModel).Error; err != nil {
		fmt.Println(err)
		return nil, []error{err}
	}

	for index, val := range tempModel {
		tempModel[index].Avatar = service.GdriveViewLink(val.Avatar)
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
func (t *GormSuite) UserEditPassword(ctx context.Context, userID int, newPassword string) (string, error) {
	if stringIsEmpty(newPassword) {
		return "Failed", gqlError("Invalid New Password", "code", "INVALID_NEW_PASSWORD")
	}

	var args []updateArgs
	args = append(args, updateArgs{
		Key:   "password",
		Value: tools.PasswordHash(newPassword),
	})
	return t.UserUpdateMultipleColumnByUserID(ctx, args, userID)
}
