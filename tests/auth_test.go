package tests

import (
	"context"
	"myapp/service"
	"testing"
)

func TestTokenValidate(t *testing.T) {
	type args struct {
		ctx context.Context
		t   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{}

	var tokenCreate string

	tokenCreate, _ = service.JwtTokenCreate(context.Background(), 1)
	tests = append(tests, struct {
		name    string
		args    args
		wantErr bool
	}{
		name:    "Test Case 1",
		args:    args{context.Background(), tokenCreate},
		wantErr: false,
	})
	tokenCreate, _ = service.JwtTokenCreate(context.Background(), 2)
	tests = append(tests, struct {
		name    string
		args    args
		wantErr bool
	}{
		name:    "Test Case 2",
		args:    args{context.Background(), tokenCreate},
		wantErr: false,
	})
	tokenCreate, _ = service.JwtTokenCreate(context.Background(), 3)
	tests = append(tests, struct {
		name    string
		args    args
		wantErr bool
	}{
		name:    "Test Case 3",
		args:    args{context.Background(), tokenCreate},
		wantErr: false,
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.TokenValidate(tt.args.ctx, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("TokenValidate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
