package requests

import "time"

type CreateWishlistRequest struct {
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	EventDate time.Time `json:"eventDate"`
}

type DeleteWishlistRequest struct {
	Id int `json:"id"`
}

type UpdateWishlistRequest struct {
	Title     string    `json:"title"`
	Desc      string    `json:"desc"`
	EventDate time.Time `json:"eventDate"`
	Id        int       `json:"id"`
}

type GetWishlistRequest struct {
	Id int `json:"id"`
}
