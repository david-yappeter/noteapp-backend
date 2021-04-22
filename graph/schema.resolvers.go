package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"myapp/graph/generated"
	"myapp/graph/model"
	"myapp/service"
)

func (r *mutationResolver) Auth(ctx context.Context) (*model.AuthOps, error) {
	return &model.AuthOps{}, nil
}

func (r *mutationResolver) User(ctx context.Context) (*model.UserOps, error) {
	return &model.UserOps{}, nil
}

func (r *mutationResolver) Team(ctx context.Context) (*model.TeamOps, error) {
	return &model.TeamOps{}, nil
}

func (r *mutationResolver) Board(ctx context.Context) (*model.BoardOps, error) {
	return &model.BoardOps{}, nil
}

func (r *mutationResolver) List(ctx context.Context) (*model.ListOps, error) {
	return &model.ListOps{}, nil
}

func (r *mutationResolver) ListItem(ctx context.Context) (*model.ListItemOps, error) {
	return &model.ListItemOps{}, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.User, error) {
	return service.UserGetByToken(ctx)
}

func (r *queryResolver) Team(ctx context.Context, id int) (*model.Team, error) {
	return service.TeamGetByIDAuthorize(ctx, id)
}

func (r *queryResolver) Board(ctx context.Context, id int) (*model.Board, error) {
	return service.BoardGetByID(ctx, id)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
