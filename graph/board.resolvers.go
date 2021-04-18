package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"myapp/graph/generated"
	"myapp/graph/model"
	"myapp/service"
)

func (r *boardOpsResolver) Create(ctx context.Context, obj *model.BoardOps, input model.NewBoard) (*model.Board, error) {
	return service.BoardCreate(ctx, input)
}

// BoardOps returns generated.BoardOpsResolver implementation.
func (r *Resolver) BoardOps() generated.BoardOpsResolver { return &boardOpsResolver{r} }

type boardOpsResolver struct{ *Resolver }
