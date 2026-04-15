package domain

import "time"

type Wishlist struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Desc       string    `json:"desc"`
	EventDate  time.Time `json:"event_date"`
	ShareToken string    `json:"share_token"`
	Gifts      []Gift    `json:"gifts"`
	UserID     int       `json:"user_id"`
}
