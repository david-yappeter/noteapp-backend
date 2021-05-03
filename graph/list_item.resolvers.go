package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"myapp/graph/generated"
	"myapp/graph/model"
	"myapp/service"
)

func (r *listItemOpsResolver) Create(ctx context.Context, obj *model.ListItemOps, input model.NewListItem) (*model.ListItem, error) {
	return service.ListItemCreateNext(ctx, input)
}

func (r *listItemOpsResolver) Move(ctx context.Context, obj *model.ListItemOps, input model.MoveListItem) (string, error) {
	return service.ListItemMovePlace(ctx, input)
}

func (r *listItemOpsResolver) UpdateName(ctx context.Context, obj *model.ListItemOps, id int, name string) (string, error) {
	return service.ListItemUpdateName(ctx, id, name)
}

func (r *listItemOpsResolver) Delete(ctx context.Context, obj *model.ListItemOps, id int) (string, error) {
	return service.ListItemDeleteByID(ctx, id)
}

// ListItemOps returns generated.ListItemOpsResolver implementation.
func (r *Resolver) ListItemOps() generated.ListItemOpsResolver { return &listItemOpsResolver{r} }

type listItemOpsResolver struct{ *Resolver }
