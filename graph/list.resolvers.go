package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"myapp/graph/generated"
	"myapp/graph/model"
)

func (r *listOpsResolver) Create(ctx context.Context, obj *model.ListOps, name string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// ListOps returns generated.ListOpsResolver implementation.
func (r *Resolver) ListOps() generated.ListOpsResolver { return &listOpsResolver{r} }

type listOpsResolver struct{ *Resolver }
