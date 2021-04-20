package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"myapp/dataloader"
	"myapp/graph/generated"
	"myapp/graph/model"
	"myapp/service"

	"github.com/99designs/gqlgen/graphql"
)

func (r *userResolver) Teams(ctx context.Context, obj *model.User) ([]*model.Team, error) {
	return dataloader.For(ctx).TeamBatchByUserIds.Load(obj.ID)
}

func (r *userOpsResolver) EditName(ctx context.Context, obj *model.UserOps, name string) (string, error) {
	return service.UserUpdateName(ctx, name)
}

func (r *userOpsResolver) EditAvatar(ctx context.Context, obj *model.UserOps, image *graphql.Upload) (*string, error) {
	return service.UserUpdateAvatar(ctx, image)
}

func (r *userOpsResolver) EditPassword(ctx context.Context, obj *model.UserOps, newPassword string) (string, error) {
	return service.UserEditPassword(ctx, newPassword)
}

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

// UserOps returns generated.UserOpsResolver implementation.
func (r *Resolver) UserOps() generated.UserOpsResolver { return &userOpsResolver{r} }

type userResolver struct{ *Resolver }
type userOpsResolver struct{ *Resolver }
