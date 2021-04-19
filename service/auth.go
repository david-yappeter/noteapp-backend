package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

var jwtSecret = []byte(os.Getenv("SECRET"))

//JwtTokenCreate Create
func JwtTokenCreate(ctx context.Context, userID int) (string, error) {
	var signingMethod = jwt.SigningMethodHS256
	var expiredTime = time.Now().UTC().Unix()

	customClaims := UserClaims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiredTime,
		},
	}

	token := jwt.NewWithClaims(signingMethod, customClaims)

	signedToken, err := token.SignedString(jwtSecret)

	if err != nil {
		fmt.Println(err)
		return "Failed", err
	}

	return signedToken, nil
}

//TokenValidate Validate
func TokenValidate(ctx context.Context, t string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(t, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println(ok)
			return nil, fmt.Errorf("there was an error")
		}

		return jwtSecret, nil
	})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return token, nil
}
