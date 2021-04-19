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

//ListDataloaderBatchByBoardIds Dataloader
func ListDataloaderBatchByBoardIds(ctx context.Context, boardIds []int) ([][]*model.List, []error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var lists []*model.List
	if err := db.Table("list").Where("board_id IN (?)", boardIds).Find(&lists).Error; err != nil {
		fmt.Println(err)
		return nil, []error{err}
	}

	itemById := map[int][]*model.List{}
	for _, val := range lists {
		itemById[val.BoardID] = append(itemById[val.BoardID], val)
	}

	items := make([][]*model.List, len(boardIds))
	for i, id := range boardIds {
		items[i] = itemById[id]
	}

	return items, nil
}
