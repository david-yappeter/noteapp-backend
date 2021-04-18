package service

import (
	"context"
	"fmt"
	"myapp/config"
	"myapp/graph/model"
	"time"
)

//ListCreate Create
func ListCreate(ctx context.Context, input model.NewList) (*model.List, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	list := model.List{
		Name:      input.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		BoardID:   input.BoardID,
	}

	if err := db.Table("list").Create(&list).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &list, nil
}
