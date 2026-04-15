package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"wishlists_project/internal/application/service"
	"wishlists_project/internal/domain/requests"
	"wishlists_project/internal/domain/responses"
)

type GiftsHandler struct {
	giftService *service.GiftService
}

func NewGiftsHandler(giftService *service.GiftService) *GiftsHandler {
	return &GiftsHandler{giftService: giftService}
}

func (handler *GiftsHandler) CreateGift(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.CreateGiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}
	if req.WishlistID == 0 {
		http.Error(w, "wishlist_id is required", http.StatusBadRequest)
		return
	}
	if req.Priority < 1 || req.Priority > 5 {
		req.Priority = 1
	}

	err := handler.giftService.CreateGift(req.Title, req.Desc, req.Link, req.WishlistID, req.Priority)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responses.SuccessResponse{Message: "Gift created successfully"})
}

func (handler *GiftsHandler) DeleteGift(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.DeleteGiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ID == 0 {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	err := handler.giftService.DeleteGift(req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.SuccessResponse{Message: "Gift deleted successfully"})
}

func (handler *GiftsHandler) ReserveGift(w http.ResponseWriter, r *http.Request) {
	var req requests.ReserveGiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.GiftID == 0 {
		http.Error(w, "gift_id is required", http.StatusBadRequest)
		return
	}
	if req.ReservedBy == "" {
		http.Error(w, "reserved_by is required", http.StatusBadRequest)
		return
	}

	err := handler.giftService.ReserveGift(req.GiftID, req.ReservedBy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.SuccessResponse{Message: "Gift reserved successfully"})
}

func (handler *GiftsHandler) GetGiftsByWishlist(w http.ResponseWriter, r *http.Request) {
	wishlistIDStr := r.URL.Query().Get("wishlist_id")
	if wishlistIDStr == "" {
		http.Error(w, "wishlist_id parameter is required", http.StatusBadRequest)
		return
	}

	wishlistID, err := strconv.Atoi(wishlistIDStr)
	if err != nil {
		http.Error(w, "Invalid wishlist_id parameter", http.StatusBadRequest)
		return
	}

	gifts, err := handler.giftService.GetWishlistGifts(wishlistID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gifts)
}
