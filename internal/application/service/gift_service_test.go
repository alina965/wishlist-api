package service

import (
	"database/sql"
	"testing"
	"wishlists_project/internal/domain"
)

type MockGiftRepository struct {
	gifts  map[int]domain.Gift
	nextID int
}

func NewMockGiftRepository() *MockGiftRepository {
	return &MockGiftRepository{
		gifts:  make(map[int]domain.Gift),
		nextID: 1,
	}
}

func (m *MockGiftRepository) CreateGift(gift *domain.Gift) error {
	gift.ID = m.nextID
	m.gifts[m.nextID] = *gift
	m.nextID++
	return nil
}

func (m *MockGiftRepository) FindGiftByID(id int) (*domain.Gift, error) {
	if gift, ok := m.gifts[id]; ok {
		return &gift, nil
	}
	return nil, sql.ErrNoRows
}

func (m *MockGiftRepository) FindGiftsByWishlistID(wishlistID int) ([]domain.Gift, error) {
	var result []domain.Gift
	for _, gift := range m.gifts {
		if gift.WishlistID == wishlistID {
			result = append(result, gift)
		}
	}
	if len(result) == 0 {
		return nil, sql.ErrNoRows
	}
	return result, nil
}

func (m *MockGiftRepository) UpdateGift(gift *domain.Gift) {
	m.gifts[gift.ID] = *gift
}

func (m *MockGiftRepository) DeleteGift(id int) error {
	delete(m.gifts, id)
	return nil
}

func TestGiftService_CreateGift_Success(t *testing.T) {
	mockRepo := NewMockGiftRepository()
	service := NewGiftService(mockRepo)

	err := service.CreateGift("iPhone 15", "New smartphone", "https://apple.com/iphone15", 1, 5)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	gifts, _ := service.GetWishlistGifts(1)
	if len(gifts) != 1 {
		t.Errorf("Expected 1 gift, got %d", len(gifts))
	}
	if gifts[0].Title != "iPhone 15" {
		t.Errorf("Expected title 'iPhone 15', got %v", gifts[0].Title)
	}
	if gifts[0].Priority != 5 {
		t.Errorf("Expected priority 5, got %d", gifts[0].Priority)
	}
	if gifts[0].Link != "https://apple.com/iphone15" {
		t.Errorf("Expected link 'https://apple.com/iphone15', got %v", gifts[0].Link)
	}
	if gifts[0].IsReserved {
		t.Error("Expected gift not to be reserved initially")
	}
}

func TestGiftService_CreateGift_MultipleGifts(t *testing.T) {
	mockRepo := NewMockGiftRepository()
	service := NewGiftService(mockRepo)

	service.CreateGift("Gift 1", "Description 1", "link1", 1, 3)
	service.CreateGift("Gift 2", "Description 2", "link2", 1, 4)
	service.CreateGift("Gift 3", "Description 3", "link3", 2, 5)

	gifts, err := service.GetWishlistGifts(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(gifts) != 2 {
		t.Errorf("Expected 2 gifts for wishlist 1, got %d", len(gifts))
	}

	gifts, err = service.GetWishlistGifts(2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(gifts) != 1 {
		t.Errorf("Expected 1 gift for wishlist 2, got %d", len(gifts))
	}
}

func TestGiftService_DeleteGift_Success(t *testing.T) {
	mockRepo := NewMockGiftRepository()
	service := NewGiftService(mockRepo)

	err := service.CreateGift("iPhone", "New phone", "", 1, 5)
	if err != nil {
		t.Fatalf("Create gift failed: %v", err)
	}

	gifts, _ := service.GetWishlistGifts(1)
	giftID := gifts[0].ID

	err = service.DeleteGift(giftID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	_, err = service.GetWishlistGifts(1)
	if err == nil {
		t.Error("Expected error after delete, got nil")
	}
}

func TestGiftService_DeleteGift_NotFound(t *testing.T) {
	mockRepo := NewMockGiftRepository()
	service := NewGiftService(mockRepo)

	err := service.DeleteGift(999)

	if err == nil {
		t.Error("Expected error for non-existent gift, got nil")
	}
	if err.Error() != "gift not found" {
		t.Errorf("Expected 'gift not found', got %v", err.Error())
	}
}

func TestGiftService_GetWishlistGifts_Success(t *testing.T) {
	mockRepo := NewMockGiftRepository()
	service := NewGiftService(mockRepo)

	service.CreateGift("Gift 1", "Desc 1", "link1", 1, 3)
	service.CreateGift("Gift 2", "Desc 2", "link2", 1, 4)

	gifts, err := service.GetWishlistGifts(1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(gifts) != 2 {
		t.Errorf("Expected 2 gifts, got %d", len(gifts))
	}
	if gifts[0].Title != "Gift 1" && gifts[1].Title != "Gift 1" {
		t.Error("Expected to find created gifts")
	}
}

func TestGiftService_GetWishlistGifts_Empty(t *testing.T) {
	mockRepo := NewMockGiftRepository()
	service := NewGiftService(mockRepo)

	_, err := service.GetWishlistGifts(999)

	if err == nil {
		t.Error("Expected error for empty wishlist, got nil")
	}
	if err.Error() != "gifts not found" {
		t.Errorf("Expected 'gifts not found', got %v", err.Error())
	}
}

func TestGiftService_ReserveGift_Success(t *testing.T) {
	mockRepo := NewMockGiftRepository()
	service := NewGiftService(mockRepo)

	err := service.CreateGift("iPhone", "New phone", "", 1, 5)
	if err != nil {
		t.Fatalf("Create gift failed: %v", err)
	}

	gifts, _ := service.GetWishlistGifts(1)
	giftID := gifts[0].ID

	err = service.ReserveGift(giftID, "Alice")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	gift, _ := mockRepo.FindGiftByID(giftID)
	if !gift.IsReserved {
		t.Error("Expected gift to be reserved")
	}
	if *gift.ReservedBy != "Alice" {
		t.Errorf("Expected reserved by 'Alice', got %v", *gift.ReservedBy)
	}
}

func TestGiftService_ReserveGift_AlreadyReserved(t *testing.T) {
	mockRepo := NewMockGiftRepository()
	service := NewGiftService(mockRepo)

	err := service.CreateGift("iPhone", "New phone", "", 1, 5)
	if err != nil {
		t.Fatalf("Create gift failed: %v", err)
	}

	gifts, _ := service.GetWishlistGifts(1)
	giftID := gifts[0].ID

	err = service.ReserveGift(giftID, "Alice")
	if err != nil {
		t.Fatalf("First reservation failed: %v", err)
	}

	err = service.ReserveGift(giftID, "Bob")

	if err == nil {
		t.Error("Expected error for already reserved gift, got nil")
	}
	if err.Error() != "gift is already reserved" {
		t.Errorf("Expected 'gift is already reserved', got %v", err.Error())
	}
}

func TestGiftService_ReserveGift_SamePerson_Unreserve(t *testing.T) {
	mockRepo := NewMockGiftRepository()
	service := NewGiftService(mockRepo)

	err := service.CreateGift("iPhone", "New phone", "", 1, 5)
	if err != nil {
		t.Fatalf("Create gift failed: %v", err)
	}

	gifts, _ := service.GetWishlistGifts(1)
	giftID := gifts[0].ID

	err = service.ReserveGift(giftID, "Alice")
	if err != nil {
		t.Fatalf("Reservation failed: %v", err)
	}

	err = service.ReserveGift(giftID, "Alice")
	if err != nil {
		t.Errorf("Expected no error for unreserve, got %v", err)
	}

	gift, _ := mockRepo.FindGiftByID(giftID)
	if gift.IsReserved {
		t.Error("Expected gift to be unreserved")
	}
	if gift.ReservedBy != nil {
		t.Errorf("Expected reserved_by to be nil, got %v", *gift.ReservedBy)
	}
}

func TestGiftService_ReserveGift_GiftNotFound(t *testing.T) {
	mockRepo := NewMockGiftRepository()
	service := NewGiftService(mockRepo)

	err := service.ReserveGift(999, "Alice")

	if err == nil {
		t.Error("Expected error for non-existent gift, got nil")
	}
	if err.Error() != "gift not found" {
		t.Errorf("Expected 'gift not found', got %v", err.Error())
	}
}

func TestGiftService_CreateGift_InvalidPriority(t *testing.T) {
	mockRepo := NewMockGiftRepository()
	service := NewGiftService(mockRepo)

	err := service.CreateGift("High priority", "Desc", "", 1, 5)
	if err != nil {
		t.Errorf("Priority 5 should be valid, got error: %v", err)
	}

	err = service.CreateGift("Low priority", "Desc", "", 1, 1)
	if err != nil {
		t.Errorf("Priority 1 should be valid, got error: %v", err)
	}
}
