package api

import (
	"encoding/json"
	"net/http"
	"wishlists_project/internal/application/service"
	"wishlists_project/internal/domain/requests"
	"wishlists_project/internal/domain/responses"
)

type WishlistsHandler struct {
	wishlistService *service.WishlistService
}

func NewWishlistsHandler(wishlistService *service.WishlistService) *WishlistsHandler {
	return &WishlistsHandler{
		wishlistService: wishlistService,
	}
}

func (handler *WishlistsHandler) CreateWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.CreateWishlistRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	err = handler.wishlistService.CreateWishlist(req.Title, req.Desc, req.EventDate, userID)
	if err != nil {
		http.Error(w, "Error creating wishlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responses.SuccessResponse{Message: "Wishlist created successfully"})
}

func (handler *WishlistsHandler) UpdateWishlist(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.UpdateWishlistRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
	}

	err = handler.wishlistService.UpdateWishlist(req.Title, req.Desc, req.EventDate, req.Id)
	if err != nil {
		http.Error(w, "Error updating wishlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.SuccessResponse{Message: "Wishlist updated successfully"})
}

func (handler *WishlistsHandler) DeleteWishlist(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.DeleteWishlistRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = handler.wishlistService.DeleteWishlist(req.Id)
	if err != nil {
		http.Error(w, "Error deleting wishlist", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responses.SuccessResponse{Message: "Wishlist deleted successfully"})
}

func (handler *WishlistsHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wishlists, err := handler.wishlistService.GetWishlistsByUserID(userID)
	if err != nil {
		http.Error(w, "Error getting wishlists", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wishlists)
}

func (handler *WishlistsHandler) GetWishlistByID(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req requests.GetWishlistRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	wishlist, err := handler.wishlistService.GetWishlistByID(req.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wishlist)
}

func (handler *WishlistsHandler) GetWishlistByToken(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "token parameter is required", http.StatusBadRequest)
		return
	}

	wishlist, err := handler.wishlistService.GetWishlistByToken(token)
	if err != nil {
		http.Error(w, "Wishlist not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wishlist)
}
