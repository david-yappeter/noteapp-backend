package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"myapp/graph/generated"
	"myapp/graph/model"
	"myapp/service"
)

func (r *teamOpsResolver) Create(ctx context.Context, obj *model.TeamOps, name string) (*model.Team, error) {
	return service.TeamCreate(ctx, name)
}

// TeamOps returns generated.TeamOpsResolver implementation.
func (r *Resolver) TeamOps() generated.TeamOpsResolver { return &teamOpsResolver{r} }

type teamOpsResolver struct{ *Resolver }
