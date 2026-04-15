package main

import (
	"fmt"
	"net/http"
	"os"
	"wishlists_project/internal/application/service"
	"wishlists_project/internal/infrastructure/api"
	"wishlists_project/internal/infrastructure/repository"
	"wishlists_project/internal/infrastructure/storage"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/wishlists?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	store, err := storage.NewStorage(dbURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer store.Close()

	userRepo := repository.NewUserRepository(store.GetDB())
	wishlistRepo := repository.NewWishlistRepository(store.GetDB())
	giftRepo := repository.NewGiftRepository(store.GetDB())

	authService := service.NewAuthService(userRepo, jwtSecret)
	wishlistService := service.NewWishlistService(wishlistRepo)
	giftService := service.NewGiftService(giftRepo)

	authHandler := api.NewAuthHandler(authService)
	wishlistHandler := api.NewWishlistsHandler(wishlistService)
	giftHandler := api.NewGiftsHandler(giftService)

	authMiddleware := api.NewAuthMiddleware(authService)

	http.HandleFunc("POST /api/register", authHandler.Register)
	http.HandleFunc("POST /api/login", authHandler.Login)
	http.HandleFunc("GET /api/wishlists/public", wishlistHandler.GetWishlistByToken)
	http.HandleFunc("POST /api/gifts/reserve", giftHandler.ReserveGift)
	http.HandleFunc("GET /api/gifts", giftHandler.GetGiftsByWishlist)

	http.HandleFunc("POST /api/wishlists", authMiddleware.Authenticate(wishlistHandler.CreateWishlist))
	http.HandleFunc("GET /api/wishlists", authMiddleware.Authenticate(wishlistHandler.GetWishlist))
	http.HandleFunc("GET /api/wishlists/get", authMiddleware.Authenticate(wishlistHandler.GetWishlistByID))
	http.HandleFunc("PUT /api/wishlists/update", authMiddleware.Authenticate(wishlistHandler.UpdateWishlist))
	http.HandleFunc("DELETE /api/wishlists/delete", authMiddleware.Authenticate(wishlistHandler.DeleteWishlist))

	http.HandleFunc("POST /api/gifts", authMiddleware.Authenticate(giftHandler.CreateGift))
	http.HandleFunc("DELETE /api/gifts", authMiddleware.Authenticate(giftHandler.DeleteGift))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println(err)
		return
	}
}
