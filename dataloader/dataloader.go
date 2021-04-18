package dataloader

import (
	"context"
	"database/sql"
	"fmt"
	"myapp/graph/model"
	"myapp/service"
	"net/http"
	"time"
)

const loadersKey = "dataloaders"

type Loaders struct {
	TeamBatchByUserIds TeamBatchLoaderByUserIds
}

func Middleware(conn *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			TeamBatchByUserIds: TeamBatchLoaderByUserIds{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int) ([][]*model.Team, []error) {
					resp, err := service.TeamBatchMapByUserIds(context.Background(), ids)

					if err != nil {
						fmt.Println(err)
						return nil, []error{err}
					}

					itemById := map[int][]*model.Team{}
					for key, val := range resp {
						itemById[key] = val
					}

					items := make([][]*model.Team, len(ids))
					for i, id := range ids {
						items[i] = itemById[id]
					}

					return items, nil
				},
			},
		})
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}
