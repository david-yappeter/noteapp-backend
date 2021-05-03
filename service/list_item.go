package service

import (
	"context"
	"fmt"
	"myapp/config"
	"myapp/graph/model"
	"time"

	"gorm.io/gorm"
)

//ListItemCreate Create
func ListItemCreate(ctx context.Context, input model.NewListItem, prev *int) (*model.ListItem, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	listItem := model.ListItem{
		Name:      input.Name,
		ListID:    &input.ListID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		Next:      nil,
		Prev:      prev,
	}

	if err := db.Table("list_item").Create(&listItem).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &listItem, nil
}

//ListItemGetLastNodeByListID By List ID
func ListItemGetLastNodeByListID(ctx context.Context, listID int) (*model.ListItem, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var listItem model.ListItem

	if err := db.Table("list_item").Where("list_id = ? AND next IS NULL", listID).Take(&listItem).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &listItem, nil
}

//ListItemCreateNext Create Next
func ListItemCreateNext(ctx context.Context, input model.NewListItem) (*model.ListItem, error) {
	if access, err := ListValidateMember(ctx, input.ListID); err != nil || !access {
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return nil, gqlError("Not Member Of Team or List doesn't exist", "code", "ACCESS_DENIED")
	}

	getListItem, err := ListItemGetLastNodeByListID(ctx, input.ListID)
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		return nil, err
	}

	var prev *int

	if getListItem == nil {
		prev = nil
	} else {
		prev = &getListItem.ID
	}

	listItem, err := ListItemCreate(ctx, model.NewListItem{
		Name:   input.Name,
		ListID: input.ListID,
	}, prev)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if getListItem != nil {
		if _, err := ListItemUpdatePointer(ctx, getListItem.ID, &listItem.ID, getListItem.Prev); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return listItem, nil
}

//ListItemUpdatePointer Update Pointer
func ListItemUpdatePointer(ctx context.Context, id int, next *int, prev *int) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
		"next":       next,
		"prev":       prev,
	}

	if err := db.Table("list_item").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//ListItemUpdatePointer Update Pointer
func ListItemUpdatePointerAndListID(ctx context.Context, id int, listID int, next *int, prev *int) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
		"list_id":    listID,
		"next":       next,
		"prev":       prev,
	}

	if err := db.Table("list_item").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//ListItemUpdatePointerValue Update Pointer Value (Ignore Null)
func ListItemUpdatePointerValue(ctx context.Context, id int, next *int, prev *int) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}

	if next != nil {
		data["next"] = next
	}
	if prev != nil {
		data["prev"] = prev
	}

	if err := db.Table("list_item").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

func ListItemUpdatePointerNull(ctx context.Context, id int, next bool, prev bool, listID bool) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}

	if next {
		data["next"] = nil
	}
	if prev {
		data["prev"] = nil
	}
	if listID {
		data["list_id"] = nil
	}

	if err := db.Table("list_item").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

func ListItemUpdatePointerNextNull(ctx context.Context, id int) (string, error) {
	return ListItemUpdatePointerNull(ctx, id, true, false, false)
}

func ListItemUpdatePointerPrevNull(ctx context.Context, id int) (string, error) {
	return ListItemUpdatePointerNull(ctx, id, false, true, false)
}

//ListItemGetByID Get By ID
func ListItemGetByID(ctx context.Context, id int) (*model.ListItem, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var listItem model.ListItem
	if err := db.Table("list_item").Where("id = ?", id).Take(&listItem).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &listItem, nil
}

//ListItemDataloaderBatchByListIds Dataloader
func ListItemDataloaderBatchByListIds(ctx context.Context, listIds []int) ([][]*model.ListItem, []error) {
	listItems, err := ListItemGetByListIds(ctx, listIds)
	if err != nil {
		fmt.Println(err)
		return nil, []error{err}
	}

	itemById := map[int][]*model.ListItem{}
	if len(listItems) > 0 {
		for _, val := range listItems {
			if val.Prev == nil {
				itemById[*val.ListID] = append([]*model.ListItem{val}, itemById[*val.ListID]...)
			} else {
				itemById[*val.ListID] = append(itemById[*val.ListID], val)
			}
		}

		listItemMapping := map[int]*model.ListItem{}
		tempItemById := map[int][]*model.ListItem{}

		for key, v := range itemById {
			var itemHead = v[0]
			for _, val := range v {
				listItemMapping[val.ID] = val
			}

			index := 0
			var sortedListItem []*model.ListItem
			for {
				sortedListItem = append(sortedListItem, itemHead)
				if itemHead.Next == nil {
					break
				}
				itemHead = listItemMapping[*itemHead.Next]
				index++
				if index > 200 {
					panic("stop")
				}
			}

			tempItemById[key] = sortedListItem
		}

		itemById = tempItemById
	}

	items := make([][]*model.ListItem, len(listIds))
	for i, id := range listIds {
		items[i] = itemById[id]
	}

	return items, nil
}

//ListItemGetByListIds By List Ids
func ListItemGetByListIds(ctx context.Context, listIds []int) ([]*model.ListItem, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var listItems []*model.ListItem
	if err := db.Table("list_item").Where("list_id IN (?)", listIds).Find(&listItems).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return listItems, nil
}

//ListItemMapGetByListIds Map
func ListItemMapGetByListIds(ctx context.Context, listIds []int) (map[int][]*model.ListItem, error) {
	listItems, err := ListItemGetByListIds(ctx, listIds)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var mappedObject = map[int][]*model.ListItem{}
	for _, val := range listItems {
		mappedObject[*val.ListID] = append(mappedObject[*val.ListID], val)
	}

	return mappedObject, nil
}

//ListItemMovePlace Move PLace
func ListItemMovePlace(ctx context.Context, input model.MoveListItem) (string, error) {
	if access, err := ListItemValidateMember(ctx, input.ID); err != nil || !access {
		if err != nil {
			fmt.Println(err)
			return "Failed", err
		}
		return "Failed", gqlError("Not Member Of Team or List Item doesn't exist", "code", "ACCESS_DENIED")
	}

	getListItem, err := ListItemGetByID(ctx, input.ID)
	if err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	// fmt.Println("Before", input.DestinationIndex)

	// if *getListItem.ListID == input.DestinationListID {
	// 	getListItemByList, err := ListItemGetByListIds(ctx, []int{input.DestinationListID})
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		return "Failed", err
	// 	}
	// 	var index int
	// 	for key, val := range getListItemByList {
	// 		if *val == *getListItem {
	// 			index = key
	// 		}
	// 	}

	// 	if index < input.DestinationIndex {
	// 		input.DestinationIndex--
	// 	}
	// }

	// fmt.Println("After", input.DestinationIndex)

	if getListItem.Next == nil && getListItem.Prev == nil {
	} else if getListItem.Next != nil && getListItem.Prev != nil {
		if _, err = ListItemUpdatePointerValue(ctx, *getListItem.Next, nil, getListItem.Prev); err != nil {
			fmt.Println(err)
			return "Failed", err
		}
		if _, err = ListItemUpdatePointerValue(ctx, *getListItem.Prev, getListItem.Next, nil); err != nil {
			fmt.Println(err)
			return "Failed", err
		}
	} else {
		var err error
		if getListItem.Next != nil {
			_, err = ListItemUpdatePointerPrevNull(ctx, *getListItem.Next)
		} else {
			_, err = ListItemUpdatePointerNextNull(ctx, *getListItem.Prev)
		}
		if err != nil {
			fmt.Println(err)
			return "Failed", err
		}
	}

	if resp, err := ListItemUpdatePointerNull(ctx, input.ID, true, true, true); err != nil {
		return resp, err
	}

	var errBatch []error
	getListItemByListDataloader, errBatch := ListItemDataloaderBatchByListIds(ctx, []int{input.DestinationListID})
	if len(errBatch) > 0 && errBatch[0] != nil {
		fmt.Println(errBatch[0])
		return "Failed", errBatch[0]
	}

	getListItemByList := getListItemByListDataloader[0]

	fmt.Println("List Item By ListID")
	for _, val := range getListItemByList {
		fmt.Printf("%+v\n", val)
	}

	lens := len(getListItemByList)

	if input.DestinationIndex == 0 {
		if lens == 0 {
			if resp, err := ListItemUpdatePointerAndListID(ctx, input.ID, input.DestinationListID, nil, nil); err != nil {
				return resp, err
			}
		} else {
			if resp, err := ListItemUpdatePointerValue(ctx, getListItemByList[0].ID, nil, &input.ID); err != nil {
				return resp, err
			}
			if resp, err := ListItemUpdatePointerAndListID(ctx, input.ID, input.DestinationListID, &getListItemByList[0].ID, nil); err != nil {
				return resp, err
			}
		}
	} else if input.DestinationIndex >= lens {
		if resp, err := ListItemUpdatePointerValue(ctx, getListItemByList[lens-1].ID, &input.ID, nil); err != nil {
			return resp, err
		}
		if resp, err := ListItemUpdatePointerAndListID(ctx, input.ID, input.DestinationListID, nil, &getListItemByList[lens-1].ID); err != nil {
			return resp, err
		}
	} else {
		if resp, err := ListItemUpdatePointerValue(ctx, getListItemByList[input.DestinationIndex-1].ID, &input.ID, nil); err != nil {
			return resp, err
		}
		if resp, err := ListItemUpdatePointerValue(ctx, getListItemByList[input.DestinationIndex].ID, nil, &input.ID); err != nil {
			return resp, err
		}
		if resp, err := ListItemUpdatePointerAndListID(ctx, input.ID, input.DestinationListID, &getListItemByList[input.DestinationIndex].ID, &getListItemByList[input.DestinationIndex-1].ID); err != nil {
			return resp, err
		}
	}

	return "Success", nil
}

// //ListItemUpdateMove Update Move
// func ListItemUpdateMove(ctx context.Context, input map[int][]*int) (string, error) {
// 	db := config.ConnectGorm()
// 	sqlDB, _ := db.DB()
// 	defer sqlDB.Close()

// 	for key, val := range input {
// 		data := map[string]interface{}{
// 			"updated_at": time.Now().UTC(),
// 		}
// 		if !(val[0] != nil && *val[0] == 0) {
// 			data["prev"] = val[0]
// 		}
// 		if !(val[1] != nil && *val[1] == 0) {
// 			data["next"] = val[1]
// 		}
// 		if len(val) == 3 {
// 			data["list_id"] = val[2]
// 		}

// 		if err := db.Table("list_item").Where("id = ?", key).Updates(data).Error; err != nil {
// 			fmt.Println(err)
// 			return "Failed", err
// 		}

// 	}

// 	return "Success", nil
// }

//ListUpdateMultipleColumnsByID Update Multiple Columns
func ListItemUpdateMultipleColumnsByID(ctx context.Context, id int, args []updateArgs) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}
	for _, val := range args {
		data[val.Key] = val.Value
	}

	if err := db.Table("list_item").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//ListUpdateName Update Name
func ListItemUpdateName(ctx context.Context, id int, name string) (string, error) {
	if stringIsEmpty(name) {
		return "Failed", gqlError("Invalid Name", "code", "INVALID_NAME")
	}

	var args []updateArgs
	args = append(args, updateArgs{
		Key:   "name",
		Value: name,
	})
	return ListItemUpdateMultipleColumnsByID(ctx, id, args)
}

func ListItemValidateMember(ctx context.Context, listItemID int) (bool, error) {
	user := ForContext(ctx)
	if user == nil {
		fmt.Println("Not Logged In!")
		return false, gqlError("Not Logged In!", "code", "NOT_LOGGED_IN")
	}

	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	var count int64

	if err := db.Table("list_item").Joins(
		"INNER JOIN list on list_item.list_id = list.id",
	).Joins(
		"INNER JOIN board on list.board_id = board.id",
	).Joins(
		"INNER JOIN team on board.team_id = team.id",
	).Joins(
		"INNER JOIN team_has_member on team_has_member.team_id = team.id",
	).Where("list_item.id = ? and team_has_member.user_id = ?", listItemID, user.ID).Count(&count).Error; err != nil {
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

func ListItemDeleteByID(ctx context.Context, id int) (string, error) {
	if access, err := ListItemValidateMember(ctx, id); err != nil || !access {
		if err != nil {
			fmt.Println(err)
			return "Failed", err
		}
		return "Failed", gqlError("Not Member Of Team or List Item doesn't exist", "code", "ACCESS_DENIED")
	}

	getListItem, err := ListItemGetByID(ctx, id)
	if err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	if getListItem.Next == nil && getListItem.Prev == nil {
	} else if getListItem.Next != nil && getListItem.Prev != nil {
		if _, err = ListItemUpdatePointerValue(ctx, *getListItem.Next, nil, getListItem.Prev); err != nil {
			fmt.Println(err)
			return "Failed", err
		}
		if _, err = ListItemUpdatePointerValue(ctx, *getListItem.Prev, getListItem.Next, nil); err != nil {
			fmt.Println(err)
			return "Failed", err
		}
	} else {
		var err error
		if getListItem.Next != nil {
			_, err = ListItemUpdatePointerPrevNull(ctx, *getListItem.Next)
		} else {
			_, err = ListItemUpdatePointerNextNull(ctx, *getListItem.Prev)
		}
		if err != nil {
			fmt.Println(err)
			return "Failed", err
		}
	}

	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if err := db.Table("list_item").Where("id = ?", id).Delete(&model.ListItem{}).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

// func ListItemGetByIds(ctx context.Context, ids []*int) ([]*model.ListItem, error) {
// 	db := config.ConnectGorm()
// 	sqlDB, _ := db.DB()
// 	defer sqlDB.Close()

// 	var listItems []*model.ListItem
// 	if err := db.Table("list_item").Where("id IN (?)", ids).Find(&listItems).Error; err != nil {
// 		fmt.Println(err)
// 		return nil, err
// 	}

// 	return listItems, nil
// }

//ListItemDeleteByBoardID Delete By Board ID
func ListItemDeleteByBoardID(ctx context.Context, boardID int) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if err := db.Exec("DELETE li.* FROM list_item as li INNER JOIN list as l on li.list_id = l.id INNER JOIN board as b on b.id = l.board_id WHERE b.id = ?", boardID).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

func ListItemDeleteByListID(ctx context.Context, listID int) (string, error) {
	db := config.ConnectGorm()
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if err := db.Exec("DELETE li.* FROM list_item as li INNER JOIN list as l on li.list_id = l.id WHERE l.id = ?", listID).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}
