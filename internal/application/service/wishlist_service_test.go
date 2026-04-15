package service

import (
	"database/sql"
	"testing"
	"time"
	"wishlists_project/internal/domain"
)

type MockWishlistRepository struct {
	wishlists map[int]domain.Wishlist
	nextID    int
}

func NewMockWishlistRepository() *MockWishlistRepository {
	return &MockWishlistRepository{
		wishlists: make(map[int]domain.Wishlist),
		nextID:    1,
	}
}

func (m *MockWishlistRepository) CreateWishlist(wishlist *domain.Wishlist) error {
	wishlist.ID = m.nextID
	m.wishlists[m.nextID] = *wishlist
	m.nextID++
	return nil
}

func (m *MockWishlistRepository) FindWishlistByID(id int) (*domain.Wishlist, error) {
	if wishlist, ok := m.wishlists[id]; ok {
		return &wishlist, nil
	}
	return nil, sql.ErrNoRows
}

func (m *MockWishlistRepository) FindWishlistsByUserID(userID int) ([]domain.Wishlist, error) {
	var result []domain.Wishlist
	for _, w := range m.wishlists {
		if w.UserID == userID {
			result = append(result, w)
		}
	}
	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}
	return result, nil
}

func (m *MockWishlistRepository) UpdateWishlist(wishlist *domain.Wishlist) error {
	m.wishlists[wishlist.ID] = *wishlist
	return nil
}

func (m *MockWishlistRepository) DeleteWishlist(id int) error {
	delete(m.wishlists, id)
	return nil
}

func (m *MockWishlistRepository) FindWishlistByToken(token string) (*domain.Wishlist, error) {
	for _, w := range m.wishlists {
		if w.ShareToken == token {
			return &w, nil
		}
	}
	return nil, sql.ErrNoRows
}

func TestWishlistService_CreateWishlist_Success(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	err := service.CreateWishlist("Birthday Party", "My birthday wishes", time.Now(), 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	wishlists, _ := service.GetWishlistsByUserID(1)
	if len(wishlists) != 1 {
		t.Errorf("Expected 1 wishlist, got %d", len(wishlists))
	}
	if wishlists[0].Title != "Birthday Party" {
		t.Errorf("Expected title 'Birthday Party', got %v", wishlists[0].Title)
	}
	if wishlists[0].Desc != "My birthday wishes" {
		t.Errorf("Expected description 'My birthday wishes', got %v", wishlists[0].Desc)
	}
	if wishlists[0].UserID != 1 {
		t.Errorf("Expected userID 1, got %d", wishlists[0].UserID)
	}
	if wishlists[0].ShareToken == "" {
		t.Error("Expected share token to be generated")
	}
}

func TestWishlistService_CreateWishlist_MultipleWishlists(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	err := service.CreateWishlist("Wishlist 1", "Desc 1", time.Now(), 1)
	if err != nil {
		t.Fatalf("Create wishlist 1 failed: %v", err)
	}
	err = service.CreateWishlist("Wishlist 2", "Desc 2", time.Now(), 1)
	if err != nil {
		t.Fatalf("Create wishlist 2 failed: %v", err)
	}

	err = service.CreateWishlist("Wishlist 3", "Desc 3", time.Now(), 2)
	if err != nil {
		t.Fatalf("Create wishlist 3 failed: %v", err)
	}

	wishlists, _ := service.GetWishlistsByUserID(1)
	if len(wishlists) != 2 {
		t.Errorf("Expected 2 wishlists for user 1, got %d", len(wishlists))
	}

	wishlists, _ = service.GetWishlistsByUserID(2)
	if len(wishlists) != 1 {
		t.Errorf("Expected 1 wishlist for user 2, got %d", len(wishlists))
	}
}

func TestWishlistService_DeleteWishlist_Success(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	err := service.CreateWishlist("Test Wishlist", "Description", time.Now(), 1)
	if err != nil {
		t.Fatalf("Create wishlist failed: %v", err)
	}

	wishlists, _ := service.GetWishlistsByUserID(1)
	wishlistID := wishlists[0].ID

	err = service.DeleteWishlist(wishlistID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, err = service.GetWishlistsByUserID(1)
	if err == nil {
		t.Error("Expected error after delete, got nil")
	}
}

func TestWishlistService_DeleteWishlist_NotFound(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	err := service.DeleteWishlist(999)

	if err == nil {
		t.Error("Expected error for non-existent wishlist, got nil")
	}
}

func TestWishlistService_UpdateWishlist_Success(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	err := service.CreateWishlist("Original Title", "Original Desc", time.Now(), 1)
	if err != nil {
		t.Fatalf("Create wishlist failed: %v", err)
	}

	wishlists, _ := service.GetWishlistsByUserID(1)
	wishlistID := wishlists[0].ID

	newEventDate := time.Now().AddDate(1, 0, 0)
	err = service.UpdateWishlist("Updated Title", "Updated Desc", newEventDate, wishlistID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	updatedWishlist, _ := service.GetWishlistByID(wishlistID)
	if updatedWishlist.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %v", updatedWishlist.Title)
	}
	if updatedWishlist.Desc != "Updated Desc" {
		t.Errorf("Expected description 'Updated Desc', got %v", updatedWishlist.Desc)
	}
}

func TestWishlistService_UpdateWishlist_PartialUpdate(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	err := service.CreateWishlist("Original Title", "Original Desc", time.Now(), 1)
	if err != nil {
		t.Fatalf("Create wishlist failed: %v", err)
	}

	wishlists, _ := service.GetWishlistsByUserID(1)
	wishlistID := wishlists[0].ID

	err = service.UpdateWishlist("Updated Title", "", time.Time{}, wishlistID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	updatedWishlist, _ := service.GetWishlistByID(wishlistID)
	if updatedWishlist.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %v", updatedWishlist.Title)
	}
	if updatedWishlist.Desc != "Original Desc" {
		t.Errorf("Expected description to remain 'Original Desc', got %v", updatedWishlist.Desc)
	}
}

func TestWishlistService_UpdateWishlist_NotFound(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	err := service.UpdateWishlist("New Title", "New Desc", time.Now(), 999)

	if err == nil {
		t.Error("Expected error for non-existent wishlist, got nil")
	}
}

func TestWishlistService_GetWishlistsByUserID_Success(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	service.CreateWishlist("Wishlist A", "Desc A", time.Now(), 1)
	service.CreateWishlist("Wishlist B", "Desc B", time.Now(), 1)

	wishlists, err := service.GetWishlistsByUserID(1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(wishlists) != 2 {
		t.Errorf("Expected 2 wishlists, got %d", len(wishlists))
	}
}

func TestWishlistService_GetWishlistsByUserID_Empty(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	_, err := service.GetWishlistsByUserID(999)

	if err == nil {
		t.Error("Expected error for user with no wishlists, got nil")
	}
}

func TestWishlistService_GetWishlistByID_Success(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	err := service.CreateWishlist("Test Wishlist", "Test Desc", time.Now(), 1)
	if err != nil {
		t.Fatalf("Create wishlist failed: %v", err)
	}

	wishlists, _ := service.GetWishlistsByUserID(1)
	wishlistID := wishlists[0].ID

	wishlist, err := service.GetWishlistByID(wishlistID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if wishlist.Title != "Test Wishlist" {
		t.Errorf("Expected title 'Test Wishlist', got %v", wishlist.Title)
	}
	if wishlist.UserID != 1 {
		t.Errorf("Expected userID 1, got %d", wishlist.UserID)
	}
}

func TestWishlistService_GetWishlistByID_NotFound(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	_, err := service.GetWishlistByID(999)

	if err == nil {
		t.Error("Expected error for non-existent wishlist, got nil")
	}
}

func TestWishlistService_GetWishlistByToken_Success(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	err := service.CreateWishlist("Test Wishlist", "Test Desc", time.Now(), 1)
	if err != nil {
		t.Fatalf("Create wishlist failed: %v", err)
	}

	wishlists, _ := service.GetWishlistsByUserID(1)
	token := wishlists[0].ShareToken

	wishlist, err := service.GetWishlistByToken(token)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if wishlist.Title != "Test Wishlist" {
		t.Errorf("Expected title 'Test Wishlist', got %v", wishlist.Title)
	}
}

func TestWishlistService_GetWishlistByToken_NotFound(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	_, err := service.GetWishlistByToken("invalid_token_12345")

	if err == nil {
		t.Error("Expected error for non-existent token, got nil")
	}
	if err.Error() != "wishlist not found" {
		t.Errorf("Expected 'wishlist not found', got %v", err.Error())
	}
}

func TestWishlistService_GenerateShareToken_Uniqueness(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	token1, err1 := service.generateShareToken()
	token2, err2 := service.generateShareToken()
	token3, err3 := service.generateShareToken()

	if err1 != nil || err2 != nil || err3 != nil {
		t.Error("Expected no errors generating tokens")
	}
	if token1 == "" || token2 == "" || token3 == "" {
		t.Error("Expected non-empty tokens")
	}
	if token1 == token2 || token2 == token3 || token1 == token3 {
		t.Error("Expected all tokens to be unique")
	}
	if len(token1) != 32 {
		t.Errorf("Expected token length 32, got %d", len(token1))
	}
}

func TestWishlistService_ShareToken_Format(t *testing.T) {
	mockRepo := NewMockWishlistRepository()
	service := NewWishlistService(mockRepo)

	token, err := service.generateShareToken()
	if err != nil {
		t.Fatalf("Generate token failed: %v", err)
	}

	for _, c := range token {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("Token contains invalid character: %c", c)
		}
	}
}
