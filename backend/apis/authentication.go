package apis

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sina-am/chess/database"
	"github.com/sina-am/chess/types"
)

type Authenticator interface {
	Authenticate(ctx context.Context, req *http.Request) (*types.User, error)
	ObtainToken(user *types.User) (string, error)
}

type jwtAuthentication struct {
	secretKey []byte
	database  database.Database
}

func NewJWTAuthentication(secretKey string, database database.Database) *jwtAuthentication {
	return &jwtAuthentication{
		secretKey: []byte(secretKey),
		database:  database,
	}
}

func (auth *jwtAuthentication) ObtainToken(user *types.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    user.Id,
		"expired_at": time.Now().Add(time.Hour * 24).Format(time.RFC822),
	})

	return token.SignedString(auth.secretKey)
}

func (auth *jwtAuthentication) Authenticate(ctx context.Context, req *http.Request) (*types.User, error) {
	tokenStr := req.Header.Get("Authorization")
	if tokenStr == "" {
		return nil, fmt.Errorf("authorization header is not set")
	}

	userId, err := auth.GetUserIdFromToken(tokenStr)
	if err != nil {
		return nil, err
	}

	return auth.database.GetUserById(ctx, userId)
}

func (auth *jwtAuthentication) GetUserIdFromToken(tokenStr string) (string, error) {
	claims, err := auth.decodeToken(tokenStr)
	if err != nil {
		return "", err
	}

	return claims["user_id"].(string), nil
}

func (auth *jwtAuthentication) decodeToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return auth.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func IsExpired(t string) error {
	expiredAt, err := time.Parse(time.RFC822, t)
	if err != nil {
		return err
	}
	if expiredAt.After(time.Now()) {
		return fmt.Errorf("expired token: %v", expiredAt)
	}
	return nil
}
