package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"myapp/graph/generated"
	"myapp/graph/model"

	"github.com/99designs/gqlgen/graphql"
)

func (r *userOpsResolver) EditName(ctx context.Context, obj *model.UserOps, name string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userOpsResolver) EditAvatar(ctx context.Context, obj *model.UserOps, image *graphql.Upload) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

// UserOps returns generated.UserOpsResolver implementation.
func (r *Resolver) UserOps() generated.UserOpsResolver { return &userOpsResolver{r} }

type userOpsResolver struct{ *Resolver }
