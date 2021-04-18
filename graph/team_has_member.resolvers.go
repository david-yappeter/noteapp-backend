package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"myapp/graph/generated"
	"myapp/graph/model"
)

func (r *teamHasMemberOpsResolver) Create(ctx context.Context, obj *model.TeamHasMemberOps, input model.NewTeamHasMember) (*model.TeamHasMember, error) {
	panic(fmt.Errorf("not implemented"))
}

// TeamHasMemberOps returns generated.TeamHasMemberOpsResolver implementation.
func (r *Resolver) TeamHasMemberOps() generated.TeamHasMemberOpsResolver {
	return &teamHasMemberOpsResolver{r}
}

type teamHasMemberOpsResolver struct{ *Resolver }
