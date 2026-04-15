package repository

import (
	"database/sql"
	"wishlists_project/internal/domain"
)

type WishlistRepository struct {
	db *sql.DB
}

func NewWishlistRepository(db *sql.DB) *WishlistRepository {
	return &WishlistRepository{db: db}
}

func (wishlistRepo *WishlistRepository) CreateWishlist(wishlist *domain.Wishlist) error {
	query := `INSERT INTO wishlists (title, description, event_date, share_token, user_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return wishlistRepo.db.QueryRow(query, wishlist.Title, wishlist.Desc, wishlist.EventDate, wishlist.ShareToken, wishlist.UserID).Scan(&wishlist.ID)
}

func (wishlistRepo *WishlistRepository) FindWishlistByID(id int) (*domain.Wishlist, error) {
	query := `SELECT title, description, event_date, share_token, user_id FROM wishlists WHERE id = $1`
	wishlist := &domain.Wishlist{}

	err := wishlistRepo.db.QueryRow(query, id).Scan(&wishlist.Title, &wishlist.Desc, &wishlist.EventDate, &wishlist.ShareToken, &wishlist.UserID)
	if err != nil {
		return nil, err
	}

	return wishlist, nil
}

func (wishlistRepo *WishlistRepository) FindWishlistsByUserID(id int) ([]domain.Wishlist, error) {
	query := `SELECT title, description, event_date, share_token FROM wishlists WHERE user_id = $1`
	wishlists := []domain.Wishlist{}

	rows, err := wishlistRepo.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		wishlist := domain.Wishlist{}
		err = rows.Scan(&wishlist.Title, &wishlist.Desc, &wishlist.EventDate, &wishlist.ShareToken)
		if err != nil {
			return nil, err
		}
		wishlists = append(wishlists, wishlist)
	}

	return wishlists, nil
}

func (wishlistRepo *WishlistRepository) UpdateWishlist(wishlist *domain.Wishlist) {
	query := `UPDATE wishlists
	SET title = $1, description = $2, event_date = $3, share_token = $4, user_id = $5
	WHERE id = $6`

	wishlistRepo.db.QueryRow(query, wishlist.Title, wishlist.Desc, wishlist.EventDate, wishlist.ShareToken, wishlist.UserID, wishlist.ID)
}

func (wishlistRepo *WishlistRepository) DeleteWishlist(id int) error {
	query := `DELETE FROM wishlists WHERE id = $1`
	_, err := wishlistRepo.db.Exec(query, id)
	return err
}

func (wishlistRepo *WishlistRepository) FindWishlistByToken(token string) (*domain.Wishlist, error) {
	query := `SELECT id, title, description, event_date, share_token, user_id FROM wishlists WHERE share_token = $1`
	wishlist := &domain.Wishlist{}

	err := wishlistRepo.db.QueryRow(query, token).Scan(&wishlist.ID, &wishlist.Title, &wishlist.Desc, &wishlist.EventDate, &wishlist.ShareToken, &wishlist.UserID)
	if err != nil {
		return nil, err
	}

	return wishlist, nil
}
