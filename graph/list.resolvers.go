package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"myapp/dataloader"
	"myapp/graph/generated"
	"myapp/graph/model"
	"myapp/service"
)

func (r *listResolver) ListItems(ctx context.Context, obj *model.List) ([]*model.ListItem, error) {
	return dataloader.For(ctx).ListItemBatchByListIds.Load(obj.ID)
}

func (r *listOpsResolver) Create(ctx context.Context, obj *model.ListOps, input model.NewList) (*model.List, error) {
	return service.ListCreate(ctx, input)
}

// List returns generated.ListResolver implementation.
func (r *Resolver) List() generated.ListResolver { return &listResolver{r} }

// ListOps returns generated.ListOpsResolver implementation.
func (r *Resolver) ListOps() generated.ListOpsResolver { return &listOpsResolver{r} }

type listResolver struct{ *Resolver }
type listOpsResolver struct{ *Resolver }
