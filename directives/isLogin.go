package directives

import (
	"context"
	"myapp/service"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func IsLogin(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	tokenUser := service.ForContext(ctx)
	if tokenUser == nil {
		return nil, &gqlerror.Error{
			Message: "access denied (Not Logged In)",
			Extensions: map[string]interface{}{
				"code": "ACCESS_DENIED",
			},
		}
	}
	return next(ctx)
}
