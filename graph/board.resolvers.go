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

func (r *boardResolver) Lists(ctx context.Context, obj *model.Board) ([]*model.List, error) {
	return dataloader.For(ctx).ListBatchByBoardIds.Load(obj.ID)
}

func (r *boardOpsResolver) Create(ctx context.Context, obj *model.BoardOps, input model.NewBoard) (*model.Board, error) {
	return service.BoardCreate(ctx, input)
}

func (r *boardOpsResolver) UpdateName(ctx context.Context, obj *model.BoardOps, id int, name string) (string, error) {
	return service.BoardUpdateName(ctx, id, name)
}

// Board returns generated.BoardResolver implementation.
func (r *Resolver) Board() generated.BoardResolver { return &boardResolver{r} }

// BoardOps returns generated.BoardOpsResolver implementation.
func (r *Resolver) BoardOps() generated.BoardOpsResolver { return &boardOpsResolver{r} }

type boardResolver struct{ *Resolver }
type boardOpsResolver struct{ *Resolver }
