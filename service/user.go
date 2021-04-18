package service

import (
	"context"
	"fmt"
	"myapp/config"
	"myapp/graph/model"
	"myapp/tools"
	"time"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

//UserCreate Create
func UserCreate(ctx context.Context, input model.NewUser) (*model.User, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	user := model.User{
		Name:      input.Name,
		Email:     input.Email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		Avatar:    nil,
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

	if err := db.Table("user").Updates(data).Where("id = ?", userID).Error; err != nil {
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
func UserUpdateAvatar(ctx context.Context, avatar *string) (string, error) {
	tokenUser := ForContext(ctx)

	var args []updateArgs
	args = append(args, updateArgs{
		Key:   "avatar",
		Value: avatar,
	})

	if _, err := UserUpdateMultipleColumnByUserID(ctx, args, tokenUser.ID); err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
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
