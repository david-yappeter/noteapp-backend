package service

import (
	"context"
	"fmt"
	"myapp/config"
	"myapp/graph/model"
	"time"

	"gorm.io/gorm"
)

//ListCreate Create
func ListCreate(ctx context.Context, input model.NewList, prev *int) (*model.List, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	list := model.List{
		Name:      input.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		BoardID:   input.BoardID,
		Next:      nil,
		Prev:      prev,
	}

	if err := db.Table("list").Create(&list).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &list, nil
}

//ListItemGetLastNodeByListID By List ID
func ListGetLastNodeByBoardID(ctx context.Context, boardID int) (*model.List, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var list model.List

	if err := db.Table("list").Where("board_id = ? AND next IS NULL", boardID).Take(&list).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &list, nil
}

//ListItemCreateNext Create Next
func ListCreateNext(ctx context.Context, input model.NewList) (*model.List, error) {
	if access, err := BoardValidateMember(ctx, input.BoardID); err != nil || !access {
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return nil, gqlError("(Not Member Of Team or Board doesn't exist", "code", "ACCESS_DENIED")
	}

	getList, err := ListGetLastNodeByBoardID(ctx, input.BoardID)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		return nil, err
	}

	var prev *int

	if getList == nil {
		prev = nil
	} else {
		prev = &getList.ID
	}

	listItem, err := ListCreate(ctx, model.NewList{
		Name:    input.Name,
		BoardID: input.BoardID,
	}, prev)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if getList != nil {
		if _, err := ListUpdatePointer(ctx, getList.ID, &listItem.ID, getList.Prev); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return listItem, nil
}

//ListUpdatePointer Update Pointer
func ListUpdatePointer(ctx context.Context, id int, next *int, prev *int) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
		"next":       next,
		"prev":       prev,
	}

	if err := db.Table("list").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//ListGetByBoardIds by Board Ids
func ListGetByID(ctx context.Context, id int) (*model.List, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var list model.List
	if err := db.Table("list").Where("id = ?", id).Take(&list).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &list, nil
}

//ListGetByBoardIds by Board Ids
func ListGetByBoardIds(ctx context.Context, boardIds []int) ([]*model.List, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var lists []*model.List
	if err := db.Table("list").Where("board_id IN (?)", boardIds).Find(&lists).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return lists, nil
}

//ListDataloaderBatchByBoardIds Dataloader
func ListDataloaderBatchByBoardIds(ctx context.Context, boardIds []int) ([][]*model.List, []error) {
	lists, err := ListGetByBoardIds(ctx, boardIds)
	if err != nil {
		fmt.Println(err)
		return nil, []error{err}
	}

	itemById := map[int][]*model.List{}
	if len(lists) > 0 {
		for _, val := range lists {
			if val.Prev == nil {
				itemById[val.BoardID] = append([]*model.List{val}, itemById[val.BoardID]...)
			} else {
				itemById[val.BoardID] = append(itemById[val.BoardID], val)
			}
		}

		listMap := map[int]*model.List{}
		tempItemById := map[int][]*model.List{}

		for key, v := range itemById {
			var itemHead = v[0]
			for _, val := range v {
				listMap[val.ID] = val
			}

			var sortedList []*model.List
			for {
				sortedList = append(sortedList, itemHead)
				if itemHead.Next == nil {
					break
				}
				itemHead = listMap[*itemHead.Next]
			}

			tempItemById[key] = sortedList
		}

		itemById = tempItemById
	}

	items := make([][]*model.List, len(boardIds))
	for i, id := range boardIds {
		items[i] = itemById[id]
	}

	return items, nil
}

//ListItemMovePlace Move PLace
func ListMovePlace(ctx context.Context, input model.MoveList) (boardID int, err error) {
	getList, err := ListGetByID(ctx, input.ID)
	if err != nil {
		fmt.Println(err)
		return
	}

	boardID = getList.BoardID

	if getList.Prev == input.MoveBeforeID && getList.Next == input.MoveAfterID {
		boardID = 0
		err = nil
		return
	}

	var updateMap = map[int][]*int{}
	if getList.Prev != nil {
		empty := 0
		updateMap[*getList.Prev] = []*int{&empty, getList.Next}
	}
	if getList.Next != nil {
		empty := 0
		updateMap[*getList.Next] = []*int{getList.Prev, &empty}
	}

	updateMap[getList.ID] = []*int{input.MoveBeforeID, input.MoveAfterID}

	if input.MoveBeforeID != nil {
		if len(updateMap[*input.MoveBeforeID]) == 2 {
			updateMap[*input.MoveBeforeID][1] = &getList.ID
		} else {
			empty := 0
			updateMap[*input.MoveBeforeID] = []*int{&empty, &getList.ID}
		}
	}
	if input.MoveAfterID != nil {
		if len(updateMap[*input.MoveAfterID]) == 2 {
			updateMap[*input.MoveAfterID][0] = &getList.ID
		} else {
			empty := 0
			updateMap[*input.MoveAfterID] = []*int{&getList.ID, &empty}
		}
	}

	if _, err = ListUpdateMove(ctx, updateMap); err != nil {
		fmt.Println(err)
		return
	}

	return
}

//ListItemUpdateMove Update Move
func ListUpdateMove(ctx context.Context, input map[int][]*int) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	for key, val := range input {
		data := map[string]interface{}{
			"updated_at": time.Now().UTC(),
		}
		if !(val[0] != nil && *val[0] == 0) {
			data["prev"] = val[0]
		}
		if !(val[1] != nil && *val[1] == 0) {
			data["next"] = val[1]
		}

		if err := db.Table("list").Where("id = ?", key).Updates(data).Error; err != nil {
			fmt.Println(err)
			return "Failed", err
		}

	}

	return "Success", nil
}

//ListValidateMember List Validate Member
func ListValidateMember(ctx context.Context, listID int) (bool, error) {
	user := ForContext(ctx)
	if user == nil {
		fmt.Println("Not Logged In!")
		return false, gqlError("Not Logged In!", "code", "NOT_LOGGED_IN")
	}

	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var count int64

	if err := db.Table("list").Joins(
		"INNER JOIN board on list.board_id = board.id",
	).Joins(
		"INNER JOIN team on board.team_id = team.id",
	).Joins(
		"INNER JOIN team_has_member on team_has_member.team_id = team.id",
	).Where("list.id = ? and team_has_member.user_id = ?", listID, user.ID).Count(&count).Error; err != nil {
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

//ListUpdateMultipleColumnsByID Update Multiple Columns
func ListUpdateMultipleColumnsByID(ctx context.Context, id int, args []updateArgs) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}
	for _, val := range args {
		data[val.Key] = val.Value
	}

	if err := db.Table("list").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//ListUpdateName Update Name
func ListUpdateName(ctx context.Context, id int, name string) (string, error) {
	if stringIsEmpty(name) {
		return "Failed", gqlError("Invalid Name", "code", "INVALID_NAME")
	}

	var args []updateArgs
	args = append(args, updateArgs{
		Key:   "name",
		Value: name,
	})
	return ListUpdateMultipleColumnsByID(ctx, id, args)
}

//ListDeleteByID Delete By ID
func ListDeleteByID(ctx context.Context, id int) (string, error) {
	if access, err := ListValidateMember(ctx, id); err != nil || !access {
		if err != nil {
			fmt.Println(err)
			return "Failed", err
		}
		return "Failed", gqlError("Not Member Of Team or List doesn't exist", "code", "ACCESS_DENIED")
	}

	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if err := db.Exec(`
    DELETE l.*, li.* 
    FROM list as l
    INNER JOIN list_item as li on li.list_id = l.id
    WHERE l.id = ?;
    `, id).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}
