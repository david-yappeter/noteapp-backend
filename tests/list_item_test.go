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
		ListID:    input.ListID,
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
				itemById[val.ListID] = append([]*model.ListItem{val}, itemById[val.ListID]...)
			} else {
				itemById[val.ListID] = append(itemById[val.ListID], val)
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
		mappedObject[val.ListID] = append(mappedObject[val.ListID], val)
	}

	return mappedObject, nil
}

//ListItemMovePlace Move PLace
func (t *GormSuite) ListItemMovePlace(ctx context.Context, input model.MoveListItem) (string, error) {
	getListItem, err := t.ListItemGetByID(ctx, input.ID)
	if err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	if getListItem.Prev == input.MoveBeforeID && getListItem.Next == input.MoveAfterID && getListItem.ListID == input.MoveBeforeListID && getListItem.ListID == input.MoveAfterListID {
		return "No Changes", nil
	}

	var updateMap = map[int][]*int{}
	if getListItem.Prev != nil {
		empty := 0
		updateMap[*getListItem.Prev] = []*int{&empty, getListItem.Next}
	}
	if getListItem.Next != nil {
		empty := 0
		updateMap[*getListItem.Next] = []*int{getListItem.Prev, &empty}
	}

	updateMap[getListItem.ID] = []*int{input.MoveBeforeID, input.MoveAfterID, &input.MoveAfterListID}

	if input.MoveBeforeID != nil {
		if len(updateMap[*input.MoveBeforeID]) == 2 {
			updateMap[*input.MoveBeforeID][1] = &getListItem.ID
		} else {
			empty := 0
			updateMap[*input.MoveBeforeID] = []*int{&empty, &getListItem.ID}
		}
	}
	if input.MoveAfterID != nil {
		if len(updateMap[*input.MoveAfterID]) == 2 {
			updateMap[*input.MoveAfterID][0] = &getListItem.ID
		} else {
			empty := 0
			updateMap[*input.MoveAfterID] = []*int{&getListItem.ID, &empty}
		}
	}

	if _, err := t.ListItemUpdateMove(ctx, updateMap); err != nil {
		fmt.Println(err)
		return "Failed", err
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
