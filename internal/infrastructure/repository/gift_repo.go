package repository

import (
	"database/sql"
	"wishlists_project/internal/domain"
)

type GiftRepository struct {
	db *sql.DB
}

func NewGiftRepository(db *sql.DB) *GiftRepository {
	return &GiftRepository{db: db}
}

func (giftRepo *GiftRepository) CreateGift(gift *domain.Gift) error {
	query := `INSERT INTO gifts (wishlist_id, title, description, link, priority) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return giftRepo.db.QueryRow(query, gift.WishlistID, gift.Title, gift.Desc, gift.Link, gift.Priority).Scan(&gift.ID)
}

func (giftRepo *GiftRepository) FindGiftByID(id int) (*domain.Gift, error) {
	query := `SELECT id, wishlist_id, title, description, link, priority, is_reserved, reserved_by FROM gifts WHERE id = $1`
	gift := &domain.Gift{}

	err := giftRepo.db.QueryRow(query, id).Scan(&gift.ID, &gift.WishlistID, &gift.Title, &gift.Desc, &gift.Link, &gift.Priority, &gift.IsReserved, &gift.ReservedBy)
	if err != nil {
		return nil, err
	}

	return gift, nil
}

func (giftRepo *GiftRepository) FindGiftsByWishlistID(id int) ([]domain.Gift, error) {
	query := `SELECT id, wishlist_id, title, description, link, priority, is_reserved, reserved_by FROM gifts WHERE wishlist_id = $1`
	gifts := []domain.Gift{}

	rows, err := giftRepo.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		gift := domain.Gift{}
		err = rows.Scan(&gift.ID, &gift.WishlistID, &gift.Title, &gift.Desc, &gift.Link, &gift.Priority, &gift.IsReserved, &gift.ReservedBy)
		if err != nil {
			return nil, err
		}
		gifts = append(gifts, gift)
	}

	return gifts, nil
}

func (giftRepo *GiftRepository) UpdateGift(gift *domain.Gift) {
	query := `UPDATE gifts SET title = $1, description = $2, link = $3, priority = $4, is_reserved = $5, reserved_by = $6 WHERE id = $7`
	giftRepo.db.Exec(query, gift.Title, gift.Desc, gift.Link, gift.Priority, gift.IsReserved, gift.ReservedBy, gift.ID)
}

func (giftRepo *GiftRepository) DeleteGift(id int) error {
	query := `DELETE FROM gifts WHERE id = $1`
	_, err := giftRepo.db.Exec(query, id)
	return err
}
