package requests

type CreateGiftRequest struct {
	WishlistID int    `json:"wishlist_id"`
	Title      string `json:"title"`
	Desc       string `json:"desc"`
	Link       string `json:"link"`
	Priority   int    `json:"priority"`
}

type DeleteGiftRequest struct {
	ID int `json:"id"`
}

type ReserveGiftRequest struct {
	GiftID     int    `json:"gift_id"`
	ReservedBy string `json:"reserved_by"`
}
