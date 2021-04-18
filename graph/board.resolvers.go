package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"myapp/graph/generated"
	"myapp/graph/model"
)

func (r *boardOpsResolver) Create(ctx context.Context, obj *model.BoardOps, name string) (*model.Board, error) {
	panic(fmt.Errorf("not implemented"))
}

// BoardOps returns generated.BoardOpsResolver implementation.
func (r *Resolver) BoardOps() generated.BoardOpsResolver { return &boardOpsResolver{r} }

type boardOpsResolver struct{ *Resolver }
