package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"myapp/graph/generated"
	"myapp/graph/model"
)

func (r *listItemOpsResolver) Create(ctx context.Context, obj *model.ListItemOps, input model.NewListItem) (*model.ListItem, error) {
	panic(fmt.Errorf("not implemented"))
}

// ListItemOps returns generated.ListItemOpsResolver implementation.
func (r *Resolver) ListItemOps() generated.ListItemOpsResolver { return &listItemOpsResolver{r} }

type listItemOpsResolver struct{ *Resolver }
