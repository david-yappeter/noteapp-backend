package tests

import (
	"context"
	"fmt"
	"myapp/graph/model"
	"time"

	"gorm.io/gorm"
)

//ListItemCreate Create
func (t *GormSuite) ListItemCreate(ctx context.Context, input model.NewListItem, prev *int) (*model.ListItem, error) {
	listItem := model.ListItem{
		Name:      input.Name,
		ListID:    &input.ListID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: nil,
		Next:      nil,
		Prev:      prev,
	}

	if err := t.tr.Table("list_item").Create(&listItem).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &listItem, nil
}

//ListItemGetLastNodeByListID By List ID
func (t *GormSuite) ListItemGetLastNodeByListID(ctx context.Context, listID int) (*model.ListItem, error) {

	var listItem model.ListItem

	if err := t.tr.Table("list_item").Where("list_id = ? AND next IS NULL", listID).Take(&listItem).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &listItem, nil
}

//ListItemCreateNext Create Next
func (t *GormSuite) ListItemCreateNext(ctx context.Context, userID int, input model.NewListItem) (*model.ListItem, error) {
	if access, err := t.ListValidateMember(ctx, userID, input.ListID); err != nil || !access {
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return nil, gqlError("Not Member Of Team or List doesn't exist", "code", "ACCESS_DENIED")
	}

	getListItem, err := t.ListItemGetLastNodeByListID(ctx, input.ListID)
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

	listItem, err := t.ListItemCreate(ctx, model.NewListItem{
		Name:   input.Name,
		ListID: input.ListID,
	}, prev)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if getListItem != nil {
		if _, err := t.ListItemUpdatePointer(ctx, getListItem.ID, &listItem.ID, getListItem.Prev); err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return listItem, nil
}

//ListItemUpdatePointer Update Pointer
func (t *GormSuite) ListItemUpdatePointer(ctx context.Context, id int, next *int, prev *int) (string, error) {

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
		"next":       next,
		"prev":       prev,
	}

	if err := t.tr.Table("list_item").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//ListItemGetByID Get By ID
func (t *GormSuite) ListItemGetByID(ctx context.Context, id int) (*model.ListItem, error) {

	var listItem model.ListItem
	if err := t.tr.Table("list_item").Where("id = ?", id).Take(&listItem).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &listItem, nil
}

//ListItemDataloaderBatchByListIds Dataloader
func (t *GormSuite) ListItemDataloaderBatchByListIds(ctx context.Context, listIds []int) ([][]*model.ListItem, []error) {
	listItems, err := t.ListItemGetByListIds(ctx, listIds)
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

			var sortedListItem []*model.ListItem
			for {
				sortedListItem = append(sortedListItem, itemHead)
				if itemHead.Next == nil {
					break
				}
				itemHead = listItemMapping[*itemHead.Next]
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
func (t *GormSuite) ListItemGetByListIds(ctx context.Context, listIds []int) ([]*model.ListItem, error) {

	var listItems []*model.ListItem
	if err := t.tr.Table("list_item").Where("list_id IN (?)", listIds).Find(&listItems).Error; err != nil {
		fmt.Println(err)
		return nil, err
	}

	return listItems, nil
}

//ListItemMapGetByListIds Map
func (t *GormSuite) ListItemMapGetByListIds(ctx context.Context, listIds []int) (map[int][]*model.ListItem, error) {
	listItems, err := t.ListItemGetByListIds(ctx, listIds)
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
func (t *GormSuite) ListItemMovePlace(ctx context.Context, userID int, input model.MoveListItem) (string, error) {
	if access, err := t.ListItemValidateMember(ctx, userID, input.ID); err != nil || !access {
		if err != nil {
			fmt.Println(err)
			return "Failed", err
		}
		return "Failed", gqlError("Not Member Of Team or List Item doesn't exist", "code", "ACCESS_DENIED")
	}

	getListItem, err := t.ListItemGetByID(ctx, input.ID)
	if err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	if getListItem.Next == nil && getListItem.Prev == nil {
	} else if getListItem.Next != nil && getListItem.Prev != nil {
		if _, err = t.ListItemUpdatePointerValue(ctx, *getListItem.Next, nil, getListItem.Prev); err != nil {
			fmt.Println(err)
			return "Failed", err
		}
		if _, err = t.ListItemUpdatePointerValue(ctx, *getListItem.Prev, getListItem.Next, nil); err != nil {
			fmt.Println(err)
			return "Failed", err
		}
	} else {
		var err error
		if getListItem.Next != nil {
			_, err = t.ListItemUpdatePointerPrevNull(ctx, *getListItem.Next)
		} else {
			_, err = t.ListItemUpdatePointerNextNull(ctx, *getListItem.Prev)
		}
		if err != nil {
			fmt.Println(err)
			return "Failed", err
		}
	}

	getListItemByList, err := t.ListItemGetByListIds(ctx, []int{input.DestinationListID})
	if err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	lens := len(getListItemByList)

	if input.DestinationIndex >= lens {
		if resp, err := t.ListItemUpdatePointerValue(ctx, getListItemByList[lens-1].ID, &input.ID, nil); err != nil {
			return resp, err
		}
		if resp, err := t.ListItemUpdatePointerAndListID(ctx, input.ID, input.DestinationListID, nil, &getListItemByList[lens-1].ID); err != nil {
			return resp, err
		}
	} else if input.DestinationIndex == 0 {
		if lens == 0 {
			if resp, err := t.ListItemUpdatePointerAndListID(ctx, input.ID, input.DestinationListID, nil, nil); err != nil {
				return resp, err
			}
		} else {
			if resp, err := t.ListItemUpdatePointerValue(ctx, getListItemByList[lens-1].ID, nil, &input.ID); err != nil {
				return resp, err
			}
			if resp, err := t.ListItemUpdatePointerAndListID(ctx, input.ID, input.DestinationListID, &getListItemByList[lens-1].ID, nil); err != nil {
				return resp, err
			}
		}
	} else {
		if resp, err := t.ListItemUpdatePointerValue(ctx, getListItemByList[lens-1].ID, &input.ID, nil); err != nil {
			return resp, err
		}
		if resp, err := t.ListItemUpdatePointerValue(ctx, getListItemByList[lens].ID, nil, &input.ID); err != nil {
			return resp, err
		}
		if resp, err := t.ListItemUpdatePointerAndListID(ctx, input.ID, input.DestinationListID, &getListItemByList[lens-1].ID, &getListItemByList[lens].ID); err != nil {
			return resp, err
		}
	}

	return "Success", nil
}

//ListItemUpdateMove Update Move
func (t *GormSuite) ListItemUpdateMove(ctx context.Context, input map[int][]*int) (string, error) {

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
		if len(val) == 3 {
			data["list_id"] = val[2]
		}

		if err := t.tr.Table("list_item").Where("id = ?", key).Updates(data).Error; err != nil {
			fmt.Println(err)
			return "Failed", err
		}

	}

	return "Success", nil
}

//ListUpdateMultipleColumnsByID Update Multiple Columns
func (t *GormSuite) ListItemUpdateMultipleColumnsByID(ctx context.Context, id int, args []updateArgs) (string, error) {

	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}
	for _, val := range args {
		data[val.Key] = val.Value
	}

	if err := t.tr.Table("list_item").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//ListUpdateName Update Name
func (t *GormSuite) ListItemUpdateName(ctx context.Context, id int, name string) (string, error) {
	if stringIsEmpty(name) {
		return "Failed", gqlError("Invalid Name", "code", "INVALID_NAME")
	}

	var args []updateArgs
	args = append(args, updateArgs{
		Key:   "name",
		Value: name,
	})
	return t.ListItemUpdateMultipleColumnsByID(ctx, id, args)
}

func (t *GormSuite) ListItemValidateMember(ctx context.Context, userID int, listItemID int) (bool, error) {
	var count int64

	if err := t.tr.Table("list_item").Joins(
		"INNER JOIN list on list_item.list_id = list.id",
	).Joins(
		"INNER JOIN board on list.board_id = board.id",
	).Joins(
		"INNER JOIN team on board.team_id = team.id",
	).Joins(
		"INNER JOIN team_has_member on team_has_member.team_id = team.id",
	).Where("list_item.id = ? and team_has_member.user_id = ?", listItemID, userID).Count(&count).Error; err != nil {
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

//ListItemUpdatePointerValue Update Pointer Value (Ignore Null)
func (t *GormSuite) ListItemUpdatePointerValue(ctx context.Context, id int, next *int, prev *int) (string, error) {
	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}

	if next != nil {
		data["next"] = next
	}
	if prev != nil {
		data["prev"] = prev
	}

	if err := t.tr.Table("list_item").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

func (t *GormSuite) ListItemUpdatePointerNextNull(ctx context.Context, id int) (string, error) {
	return t.ListItemUpdatePointerNull(ctx, id, true, false)
}

func (t *GormSuite) ListItemUpdatePointerPrevNull(ctx context.Context, id int) (string, error) {
	return t.ListItemUpdatePointerNull(ctx, id, false, true)
}

func (t *GormSuite) ListItemUpdatePointerNull(ctx context.Context, id int, next bool, prev bool) (string, error) {
	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
	}

	if next {
		data["next"] = nil
	}
	if prev {
		data["prev"] = nil
	}

	if err := t.tr.Table("list_item").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}

//ListItemUpdatePointer Update Pointer
func (t *GormSuite) ListItemUpdatePointerAndListID(ctx context.Context, id int, listID int, next *int, prev *int) (string, error) {
	data := map[string]interface{}{
		"updated_at": time.Now().UTC(),
		"list_id":    listID,
		"next":       next,
		"prev":       prev,
	}

	if err := t.tr.Table("list_item").Where("id = ?", id).Updates(data).Error; err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return "Success", nil
}
