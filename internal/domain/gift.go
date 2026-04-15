package domain

type Gift struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	Desc       string  `json:"description"`
	Link       string  `json:"link"`
	IsReserved bool    `json:"is_reserved"`
	ReservedBy *string `json:"reserved_by"`
	WishlistID int     `json:"wishlist_id"`
	Priority   int     `json:"priority"`
}
