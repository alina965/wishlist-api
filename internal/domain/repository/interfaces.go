package repository

import (
	"wishlists_project/internal/domain"
)

type UserRepositoryInterface interface {
	CreateUser(user *domain.User) error
	GetUserById(id int) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
}

type WishlistRepositoryInterface interface {
	CreateWishlist(wishlist *domain.Wishlist) error
	FindWishlistByID(id int) (*domain.Wishlist, error)
	FindWishlistsByUserID(userID int) ([]domain.Wishlist, error)
	UpdateWishlist(wishlist *domain.Wishlist)
	DeleteWishlist(id int) error
	FindWishlistByToken(token string) (*domain.Wishlist, error)
}

type GiftRepositoryInterface interface {
	CreateGift(gift *domain.Gift) error
	FindGiftByID(id int) (*domain.Gift, error)
	FindGiftsByWishlistID(wishlistID int) ([]domain.Gift, error)
	UpdateGift(gift *domain.Gift)
	DeleteGift(id int) error
}
