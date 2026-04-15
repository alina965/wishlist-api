package service

import (
	"database/sql"
	"testing"
	"wishlists_project/internal/domain"
)

type MockUserRepository struct {
	users map[string]domain.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]domain.User),
	}
}

func (m *MockUserRepository) CreateUser(user *domain.User) error {
	user.ID = len(m.users) + 1
	m.users[user.Email] = *user
	return nil
}

func (m *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	if user, ok := m.users[email]; ok {
		return &user, nil
	}
	return nil, sql.ErrNoRows
}

func (m *MockUserRepository) GetUserById(id int) (*domain.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, sql.ErrNoRows
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewAuthService(mockRepo, "test-secret")

	user, err := service.Register("test@example.com", "Test User", "password123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %v", user.Email)
	}
	if user.Name != "Test User" {
		t.Errorf("Expected name Test User, got %v", user.Name)
	}
	if user.Token == "" {
		t.Error("Expected token to be generated")
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewAuthService(mockRepo, "test-secret")

	_, err := service.Register("test@example.com", "Test User", "password123")
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	_, err = service.Register("test@example.com", "Another User", "password456")

	if err == nil {
		t.Error("Expected error for duplicate email, got nil")
	}
	if err.Error() != "this email is already taken" {
		t.Errorf("Expected 'this email is already taken', got %v", err.Error())
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewAuthService(mockRepo, "test-secret")

	_, err := service.Register("test@example.com", "Test User", "password123")
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	user, err := service.Login("test@example.com", "password123")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %v", user.Email)
	}
	if user.Token == "" {
		t.Error("Expected token to be generated")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewAuthService(mockRepo, "test-secret")

	_, err := service.Register("test@example.com", "Test User", "password123")
	if err != nil {
		t.Fatalf("Registration failed: %v", err)
	}

	_, err = service.Login("test@example.com", "wrongpassword")

	if err == nil {
		t.Error("Expected error for wrong password, got nil")
	}
	if err.Error() != "invalid password" {
		t.Errorf("Expected 'invalid password', got %v", err.Error())
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	service := NewAuthService(mockRepo, "test-secret")

	_, err := service.Login("nonexistent@example.com", "password123")

	if err == nil {
		t.Error("Expected error for nonexistent user, got nil")
	}
	if err.Error() != "cannot find user by email" {
		t.Errorf("Expected 'cannot find user by email', got %v", err.Error())
	}
}
