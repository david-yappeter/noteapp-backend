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
)

func (r *listResolver) ListItems(ctx context.Context, obj *model.List) ([]*model.ListItem, error) {
	return dataloader.For(ctx).ListItemBatchByListIds.Load(obj.ID)
}

func (r *listOpsResolver) Create(ctx context.Context, obj *model.ListOps, input model.NewList) (*model.List, error) {
	return service.ListCreateNext(ctx, input)
}

func (r *listOpsResolver) Move(ctx context.Context, obj *model.ListOps, input model.MoveList) ([]*model.List, error) {
	boardID, err := service.ListMovePlace(ctx, input)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return dataloader.For(ctx).ListBatchByBoardIds.Load(boardID)
}

func (r *listOpsResolver) UpdateName(ctx context.Context, obj *model.ListOps, id int, name string) (string, error) {
	return service.ListUpdateName(ctx, id, name)
}

// List returns generated.ListResolver implementation.
func (r *Resolver) List() generated.ListResolver { return &listResolver{r} }

// ListOps returns generated.ListOpsResolver implementation.
func (r *Resolver) ListOps() generated.ListOpsResolver { return &listOpsResolver{r} }

type listResolver struct{ *Resolver }
type listOpsResolver struct{ *Resolver }
