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

// CreateGift
// @Summary      Add gift to wishlist
// @Description  Creates a new gift in specified wishlist
// @Tags         gifts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body requests.CreateGiftRequest true "Gift data"
// @Success      201 {object} responses.SuccessResponse "Gift created"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Router       /gifts [post]
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

// DeleteGift
// @Summary      Delete gift
// @Description  Removes a gift from wishlist
// @Tags         gifts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body requests.DeleteGiftRequest true "Gift ID"
// @Success      200 {object} responses.SuccessResponse "Gift deleted"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      401 {object} map[string]string "Unauthorized"
// @Router       /gifts [delete]
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

// ReserveGift
// @Summary      Reserve a gift
// @Description      Reserves a gift by share token (public endpoint)
// @Tags         gifts
// @Accept       json
// @Produce      json
// @Param        request body requests.ReserveGiftRequest true "Reservation data"
// @Success      200 {object} responses.SuccessResponse "Gift reserved"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      409 {object} map[string]string "Gift already reserved"
// @Router       /gifts/reserve [post]
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

// GetGiftsByWishlist
// @Summary      Get gifts by wishlist ID
// @Description  Returns all gifts for a specific wishlist (public)
// @Tags         gifts
// @Produce      json
// @Param        wishlist_id query int true "Wishlist ID"
// @Success      200 {array} domain.Gift "List of gifts"
// @Failure      400 {object} map[string]string "Invalid wishlist_id"
// @Failure      404 {object} map[string]string "No gifts found"
// @Router       /gifts [get]
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
