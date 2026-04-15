package repository

import (
	"database/sql"
	"wishlists_project/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (userRepo *UserRepository) CreateUser(user *domain.User) error {
	query := `INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id`
	return userRepo.db.QueryRow(query, user.Email, user.Password, user.Name).Scan(&user.ID)
}

func (userRepo *UserRepository) GetUserById(id int) (*domain.User, error) {
	query := `SELECT id, email, name, password FROM users WHERE id = $1`
	user := &domain.User{}

	err := userRepo.db.QueryRow(query, id).Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return &domain.User{}, err
	}

	return user, nil
}

func (userRepo *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	query := `SELECT id, email, name, password FROM users WHERE email = $1`
	user := &domain.User{}

	err := userRepo.db.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		return &domain.User{}, err
	}
	return user, nil
}
