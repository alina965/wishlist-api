package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"wishlists_project/internal/domain"
	"wishlists_project/internal/domain/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo  repository.UserRepositoryInterface
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepositoryInterface, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (service *AuthService) Register(email string, name string, password string) (domain.User, error) {
	_, err := service.userRepo.GetUserByEmail(email)
	if err == nil {
		return domain.User{}, errors.New("this email is already taken")
	}

	if !errors.Is(err, sql.ErrNoRows) {
		// Произошла другая ошибка БД (не "пользователь не найден")
		return domain.User{}, fmt.Errorf("database error: %w", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, fmt.Errorf("error hashing password: %w", err)
	}

	user := domain.User{Email: email, Name: name, Password: string(hashedPassword)}
	err = service.userRepo.CreateUser(&user)
	if err != nil {
		return domain.User{}, errors.New("cannot create new user")
	}

	token, err := service.generateToken(user.ID)
	if err != nil {
		return domain.User{}, errors.New("error generating token")
	}

	user.Token = token

	return user, nil
}

func (service *AuthService) Login(email string, password string) (domain.User, error) {
	user, err := service.userRepo.GetUserByEmail(email)
	if err != nil {
		return domain.User{}, errors.New("cannot find user by email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return domain.User{}, errors.New("invalid password")
	}

	token, err := service.generateToken(user.ID)
	if err != nil {
		return domain.User{}, errors.New("error generating token")
	}

	user.Token = token
	user.Password = ""

	return *user, nil
}

func (service *AuthService) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(service.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (service *AuthService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(service.jwtSecret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token")
	}

	userID := int(claims["user_id"].(float64))
	return userID, nil
}
