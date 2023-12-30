package users

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
	ErrMissingToken = errors.New("missing token")
)

type Authenticator interface {
	Authenticate(ctx context.Context, token string) (*User, error)
	ObtainToken(user *User) (string, error)
	Login(c echo.Context, user *User) error
	GetUser(c echo.Context) UserI
}

type JwtToken struct {
	UserId    primitive.ObjectID
	ExpiredAt time.Time
}

type jwtAuthentication struct {
	secretKey []byte
	storage   Storage
}

func NewJWTAuthentication(secretKey string, storage Storage) *jwtAuthentication {
	return &jwtAuthentication{
		secretKey: []byte(secretKey),
		storage:   storage,
	}

}

func (auth *jwtAuthentication) ObtainToken(user *User) (string, error) {
	jwtToken := JwtToken{
		UserId:    user.Id,
		ExpiredAt: time.Now().Add(time.Hour * 6),
	}

	return auth.Encode(jwtToken)
}

func (auth *jwtAuthentication) Encode(jwtToken JwtToken) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    jwtToken.UserId.Hex(),
		"expired_at": jwtToken.ExpiredAt.Format(time.RFC822Z),
	})

	return token.SignedString(auth.secretKey)
}

func (auth *jwtAuthentication) Authenticate(ctx context.Context, tokenStr string) (*User, error) {
	if tokenStr == "" {
		return nil, ErrInvalidToken
	}

	token, err := auth.DecodeToken(tokenStr)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if err := IsExpired(token.ExpiredAt); err != nil {
		return nil, ErrExpiredToken
	}

	user, err := auth.storage.GetUserById(ctx, token.UserId)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return user, nil
}

func (auth *jwtAuthentication) DecodeToken(tokenStr string) (JwtToken, error) {
	claims, err := auth.decodeToken(tokenStr)
	if err != nil {
		return JwtToken{}, err
	}

	userIdStr, ok := claims["user_id"]
	if !ok {
		return JwtToken{}, fmt.Errorf("invalid token: user_id does not exist")
	}
	userId, err := UserIdFromString(userIdStr.(string))
	if err != nil {
		return JwtToken{}, fmt.Errorf("invalid token: user_id is invalid ObjectId")
	}

	expiredAtStr, ok := claims["expired_at"]
	if !ok {
		return JwtToken{}, fmt.Errorf("invalid token: expired_at does not exist")
	}
	expiredAt, err := time.Parse(time.RFC822Z, expiredAtStr.(string))
	if err != nil {
		return JwtToken{}, fmt.Errorf("invalid token: expired_at is invalid time.Time")
	}

	return JwtToken{
		UserId:    userId,
		ExpiredAt: expiredAt,
	}, nil
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

func IsExpired(t time.Time) error {
	if !t.After(time.Now()) {
		return fmt.Errorf("expired token: %v", t)
	}
	return nil
}

func (auth *jwtAuthentication) Login(c echo.Context, user *User) error {
	token, err := auth.ObtainToken(user)
	if err != nil {
		return err
	}
	c.SetCookie(&http.Cookie{Name: "sessionID", Value: token, Path: "/"})
	return nil
}

func (auth *jwtAuthentication) GetUser(c echo.Context) UserI {
	cookie, err := c.Cookie("sessionID")
	if err != nil {
		return NewAnonymousUser()
	}

	user, err := auth.Authenticate(c.Request().Context(), cookie.Value)
	if err != nil {
		return NewAnonymousUser()
	}

	return user
}
