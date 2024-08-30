package service

import (
	"context"
	"crypto/sha1"
	"fmt"
	"github.com/ZnNr/notes-keeper.git/intenal/errors"
	"github.com/ZnNr/notes-keeper.git/intenal/midleware/repository"
	"github.com/ZnNr/notes-keeper.git/intenal/users/usermodel"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"time"
)

type AuthService struct {
	repo     repository.User
	logger   *slog.Logger
	SignKey  string
	TokenTTL time.Duration
	Salt     string
}

type AuthDependencies struct {
	userRepo repository.User
	logger   *slog.Logger
	signKey  string
	tokenTTL time.Duration
	salt     string
}

func NewAuthService(deps AuthDependencies) *AuthService {
	fmt.Println(deps)
	return &AuthService{
		repo:     deps.userRepo,
		logger:   deps.logger,
		SignKey:  deps.signKey,
		TokenTTL: deps.tokenTTL,
		Salt:     deps.salt,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	const op = "service.Auth.Login"
	s.logger.With("op", op)

	if username == "" {
		s.logger.Error("username is required")
		return "", errors.ErrUsernameRequired
	}
	if password == "" {
		s.logger.Error("password is required")
		return "", errors.ErrPasswordRequired
	}
	user, err := s.repo.GetUser(ctx, username, s.generatePasswordHash(password))
	if err != nil {
		s.logger.Error("cannot get user", slog.String("username", username))
		return "", errors.ErrCannotGetUser
	}
	token, err := s.generateToken(user)
	if err != nil {
		s.logger.Error("cannot generate token", slog.String("username", username))
		return "", errors.ErrCannotGenerateToken
	}
	return token, nil
}

func (s *AuthService) Register(ctx context.Context, username, password string) error {
	const op = "service.Auth.Register"
	s.logger.With("op", op)

	if username == "" {
		s.logger.Error("username is required")
		return errors.ErrUsernameRequired
	}
	if password == "" {
		s.logger.Error("password is required")
		return errors.ErrPasswordRequired
	}
	_, err := s.repo.CreateUser(ctx, username, s.generatePasswordHash(password))
	if err != nil {
		if err == errors.ErrUserAlreadyExists {
			s.logger.Error("user already exists", slog.String("username", username))
			return errors.ErrUserAlreadyExists
		}
		s.logger.Error("cannot create user", slog.String("username", username))
		return errors.ErrCannotCreateUser
	}
	return nil
}

func (s *AuthService) generateToken(user usermodel.User) (string, error) {
	const op = "service.Auth.generateToken"
	s.logger.With("op", op)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expires_at": time.Now().Add(s.TokenTTL).Unix(),
		"issued_at":  time.Now().Unix(),
		"user_id":    user.Id,
	})
	tokenString, err := token.SignedString([]byte(s.SignKey))
	if err != nil {
		s.logger.Error("cannot sign token", slog.String("username", user.Username))
		return "", errors.ErrCannotSignToken
	}
	return tokenString, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	const op = "service.Auth.ParseToken"
	s.logger.With("op", op)

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Error(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.SignKey), nil
	})
	if err != nil {
		s.logger.Error("cannot parse token", slog.String("accessToken", accessToken))
		return 0, errors.ErrCannotParseToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error("cannot parse token claims", slog.String("accessToken", accessToken))
		return 0, errors.ErrCannotParseToken
	}
	userId, ok := claims["user_id"].(float64)
	if !ok {
		s.logger.Error("cannot get user_id from token claims", slog.String("accessToken", accessToken))
		return 0, errors.ErrCannotParseToken
	}
	return int(userId), nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(s.Salt)))
}
