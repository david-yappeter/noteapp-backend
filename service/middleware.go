package service

import (
	"context"
	"net/http"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	user string
}

//Auth Middleware Token Check
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authorization := r.Header.Get("Authorization")

			if authorization == "" {
				next.ServeHTTP(w, r)
				return
			}
			tokenBeforeClaims, err := TokenValidate(context.Background(), authorization)

			if err != nil {
				http.Error(w, "Invalid Token", http.StatusForbidden)
				return
			}

			claims, ok := tokenBeforeClaims.Claims.(*UserClaims)
			if !ok && !tokenBeforeClaims.Valid {
				http.Error(w, "Invalid Token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, &UserClaims{
				ID: claims.ID,
			})

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		})
	}
}

func ForContext(ctx context.Context) *UserClaims {
	raw, _ := ctx.Value(userCtxKey).(*UserClaims)
	return raw
}
