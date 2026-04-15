package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	"wishlists_project/internal/domain"
	"wishlists_project/internal/domain/repository"
)

type WishlistService struct {
	wishlistRepo repository.WishlistRepositoryInterface
}

func NewWishlistService(wishlistRepo repository.WishlistRepositoryInterface) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
	}
}

func (service *WishlistService) CreateWishlist(title string, desc string, eventDate time.Time, userId int) error {
	shareToken, err := service.generateShareToken()
	if err != nil {
		return errors.New("error generating share token" + err.Error())
	}

	wishlist := &domain.Wishlist{
		Title:      title,
		Desc:       desc,
		EventDate:  eventDate,
		ShareToken: shareToken,
		UserID:     userId,
	}

	err = service.wishlistRepo.CreateWishlist(wishlist)
	if err != nil {
		return errors.New("cannot create new wishlist" + err.Error())
	}

	return nil
}

func (service *WishlistService) DeleteWishlist(id int) error {
	_, err := service.wishlistRepo.FindWishlistByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("wishlist not found")
		}
		return errors.New("cannot delete wishlist: " + err.Error())
	}

	err = service.wishlistRepo.DeleteWishlist(id)
	if err != nil {
		return errors.New("cannot delete wishlist" + err.Error())
	}
	return nil
}

func (service *WishlistService) UpdateWishlist(title string, desc string, eventDate time.Time, id int) error {
	wishlist, err := service.wishlistRepo.FindWishlistByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("wishlist not found")
		}
		return errors.New("cannot update wishlist: " + err.Error())
	}

	if title != "" {
		wishlist.Title = title
	}

	if desc != "" {
		wishlist.Desc = desc
	}

	if !eventDate.IsZero() {
		wishlist.EventDate = eventDate
	}

	err = service.wishlistRepo.UpdateWishlist(wishlist)
	if err != nil {
		return errors.New("cannot update wishlist: " + err.Error())
	}

	return nil
}

func (service *WishlistService) generateShareToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (service *WishlistService) GetWishlistsByUserID(userID int) ([]domain.Wishlist, error) {
	wishlists, err := service.wishlistRepo.FindWishlistsByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wishlists: %w", err)
	}

	return wishlists, nil
}

func (service *WishlistService) GetWishlistByID(ID int) (domain.Wishlist, error) {
	wishlist, err := service.wishlistRepo.FindWishlistByID(ID)
	if err != nil {
		return domain.Wishlist{}, fmt.Errorf("failed to get wishlist by id: %w", err)
	}

	return *wishlist, nil
}

func (service *WishlistService) GetWishlistByToken(token string) (*domain.Wishlist, error) {
	wishlist, err := service.wishlistRepo.FindWishlistByToken(token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("wishlist not found")
		}
		return nil, fmt.Errorf("failed to get wishlist by token: %w", err)
	}

	return wishlist, nil
}
