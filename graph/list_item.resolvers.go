package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"myapp/dataloader"
	"myapp/graph/generated"
	"myapp/graph/model"
	"myapp/service"
	"strconv"
)

func (r *listItemOpsResolver) Create(ctx context.Context, obj *model.ListItemOps, input model.NewListItem) (*model.ListItem, error) {
	return service.ListItemCreateNext(ctx, input)
}

func (r *listItemOpsResolver) Move(ctx context.Context, obj *model.ListItemOps, input model.MoveListItem) (map[string]interface{}, error) {
	if _, err := service.ListItemMovePlace(ctx, input); err != nil {
		fmt.Println(err)
		return nil, err
	}

	mappedObject := map[string]interface{}{}
	var err error
	mappedObject[strconv.Itoa(input.MoveBeforeListID)], err = dataloader.For(ctx).ListItemBatchByListIds.Load(input.MoveBeforeListID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if input.MoveAfterListID != input.MoveBeforeListID {
		mappedObject[strconv.Itoa(input.MoveAfterListID)], err = dataloader.For(ctx).ListItemBatchByListIds.Load(input.MoveAfterListID)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}

	return mappedObject, nil
}

// ListItemOps returns generated.ListItemOpsResolver implementation.
func (r *Resolver) ListItemOps() generated.ListItemOpsResolver { return &listItemOpsResolver{r} }

type listItemOpsResolver struct{ *Resolver }
