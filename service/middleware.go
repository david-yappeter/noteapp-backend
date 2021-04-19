package service

import (
	"context"
	"net/http"
	"strings"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	user string
}

//Auth Middleware Token Check
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorization := r.Header.Get("Authorization")

		if authorization == "" {
			next.ServeHTTP(w, r)
			return
		}

		splitBearer := strings.Split(authorization, " ")
		if len(splitBearer) != 2 {
			http.Error(w, "Invalid Token", http.StatusForbidden)
			return
		}

		if !strings.EqualFold(splitBearer[0], "bearer") {
			http.Error(w, "Invalid Token", http.StatusForbidden)
			return
		}

		tokenBeforeClaims, err := TokenValidate(context.Background(), splitBearer[1])

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
	})
}

func ForContext(ctx context.Context) *UserClaims {
	raw, _ := ctx.Value(userCtxKey).(*UserClaims)
	return raw
}
