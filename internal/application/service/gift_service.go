package service

import (
	"errors"
	"wishlists_project/internal/domain"
	"wishlists_project/internal/domain/repository"
)

type GiftService struct {
	giftRepo repository.GiftRepositoryInterface
}

func NewGiftService(giftRepo repository.GiftRepositoryInterface) *GiftService {
	return &GiftService{
		giftRepo: giftRepo,
	}
}

func (service *GiftService) CreateGift(title, description, link string, wishlistID, priority int) error {
	gift := &domain.Gift{Title: title, Desc: description, Link: link, IsReserved: false, WishlistID: wishlistID, Priority: priority}
	err := service.giftRepo.CreateGift(gift)
	if err != nil {
		return err
	}
	return nil
}

func (service *GiftService) DeleteGift(id int) error {
	_, err := service.giftRepo.FindGiftByID(id)
	if err != nil {
		return errors.New("gift not found")
	}

	err = service.giftRepo.DeleteGift(id)
	if err != nil {
		return err
	}

	return nil
}

func (service *GiftService) GetWishlistGifts(id int) ([]domain.Gift, error) {
	gifts, err := service.giftRepo.FindGiftsByWishlistID(id)
	if err != nil {
		return nil, errors.New("gifts not found")
	}
	if len(gifts) == 0 {
		return nil, errors.New("there are no gifts with the given wishlist id")
	}

	return gifts, nil
}

func (service *GiftService) ReserveGift(id int, reservedBy string) error {
	gift, err := service.giftRepo.FindGiftByID(id)
	if err != nil {
		return errors.New("gift not found")
	}
	if gift.IsReserved == true {
		if *gift.ReservedBy == reservedBy {
			gift.IsReserved = false
			gift.ReservedBy = nil
		} else {
			return errors.New("gift is already reserved")
		}
	} else {
		gift.IsReserved = true
		gift.ReservedBy = &reservedBy
	}

	service.giftRepo.UpdateGift(gift)

	return nil
}
