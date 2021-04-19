package dataloader

import (
	"context"
	"myapp/graph/model"
	"myapp/service"
	"net/http"
	"time"
)

const loadersKey = "dataloaders"

type Loaders struct {
	TeamBatchByUserIds     TeamBatchLoaderByUserIds
	UserBatchByTeamIds     UserBatchLoaderByTeamIds
	BoardBatchByTeamIds    BoardBatchLoaderByTeamIds
	ListBatchByBoardIds    ListBatchLoaderByBoardIds
	ListItemBatchByListIds ListItemBatchLoaderByListIds
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), loadersKey, &Loaders{
			TeamBatchByUserIds: TeamBatchLoaderByUserIds{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int) ([][]*model.Team, []error) {
					return service.TeamDataLoaderBatchByUserIds(context.Background(), ids)
				},
			},
			UserBatchByTeamIds: UserBatchLoaderByTeamIds{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int) ([][]*model.User, []error) {
					return service.UserDataloaderBatchByTeamIds(context.Background(), ids)
				},
			},
			BoardBatchByTeamIds: BoardBatchLoaderByTeamIds{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int) ([][]*model.Board, []error) {
					return service.BoardDataloaderBatchByTeamIds(context.Background(), ids)
				},
			},
			ListBatchByBoardIds: ListBatchLoaderByBoardIds{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int) ([][]*model.List, []error) {
					return service.ListDataloaderBatchByBoardIds(context.Background(), ids)
				},
			},
			ListItemBatchByListIds: ListItemBatchLoaderByListIds{
				maxBatch: 100,
				wait:     1 * time.Millisecond,
				fetch: func(ids []int) ([][]*model.ListItem, []error) {
					return service.ListItemDataloaderBatchByListIds(context.Background(), ids)
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
