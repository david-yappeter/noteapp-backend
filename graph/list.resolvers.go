package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"myapp/graph/generated"
	"myapp/graph/model"
	"myapp/service"
)

func (r *listOpsResolver) Create(ctx context.Context, obj *model.ListOps, input model.NewList) (*model.List, error) {
	return service.ListCreate(ctx, input)
}

// ListOps returns generated.ListOpsResolver implementation.
func (r *Resolver) ListOps() generated.ListOpsResolver { return &listOpsResolver{r} }

type listOpsResolver struct{ *Resolver }
