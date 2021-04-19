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

func (r *teamResolver) Members(ctx context.Context, obj *model.Team) ([]*model.User, error) {
	return dataloader.For(ctx).UserBatchByTeamIds.Load(obj.ID)
}

func (r *teamResolver) Boards(ctx context.Context, obj *model.Team) ([]*model.Board, error) {
	return dataloader.For(ctx).BoardBatchByTeamIds.Load(obj.ID)
}

func (r *teamOpsResolver) Create(ctx context.Context, obj *model.TeamOps, name string) (*model.Team, error) {
	return service.TeamCreate(ctx, name)
}

func (r *teamOpsResolver) AddMember(ctx context.Context, obj *model.TeamOps, input model.NewTeamHasMember) (*model.TeamHasMember, error) {
	return service.TeamAddMember(ctx, input)
}

func (r *teamOpsResolver) RemoveMember(ctx context.Context, obj *model.TeamOps, input model.NewTeamHasMember) (string, error) {
	return service.TeamRemoveMember(ctx, input)
}

// Team returns generated.TeamResolver implementation.
func (r *Resolver) Team() generated.TeamResolver { return &teamResolver{r} }

// TeamOps returns generated.TeamOpsResolver implementation.
func (r *Resolver) TeamOps() generated.TeamOpsResolver { return &teamOpsResolver{r} }

type teamResolver struct{ *Resolver }
type teamOpsResolver struct{ *Resolver }
