package main

import (
	"fmt"
	"net/http"
	"os"
	_ "wishlists_project/docs"
	"wishlists_project/internal/application/service"
	"wishlists_project/internal/infrastructure/api"
	"wishlists_project/internal/infrastructure/repository"
	"wishlists_project/internal/infrastructure/storage"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Wishlist API
// @version         1.0
// @description     Wishlist management service for creating and sharing wishlists
// @host           localhost:8080
// @BasePath       /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter JWT token in format: Bearer <token>

const (
	defaultDBURL     = "postgres://postgres:password@localhost:5432/wishlists?sslmode=disable"
	defaultJWTSecret = "secret-key"
	defaultPort      = "8080"
)

type config struct {
	dbURL     string
	jwtSecret string
	port      string
}

func loadConfig() *config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = defaultDBURL
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = defaultJWTSecret
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	return &config{
		dbURL:     dbURL,
		jwtSecret: jwtSecret,
		port:      port,
	}
}

type handlers struct {
	auth *api.AuthHandler
	wish *api.WishlistsHandler
	gift *api.GiftsHandler
	mw   *api.AuthMiddleware
}

func initHandlers(cfg *config) (*handlers, *storage.Storage, error) {
	store, err := storage.NewStorage(cfg.dbURL)
	if err != nil {
		return nil, nil, err
	}

	userRepo := repository.NewUserRepository(store.GetDB())
	wishlistRepo := repository.NewWishlistRepository(store.GetDB())
	giftRepo := repository.NewGiftRepository(store.GetDB())

	authService := service.NewAuthService(userRepo, cfg.jwtSecret)
	wishlistService := service.NewWishlistService(wishlistRepo)
	giftService := service.NewGiftService(giftRepo)

	return &handlers{
		auth: api.NewAuthHandler(authService),
		wish: api.NewWishlistsHandler(wishlistService),
		gift: api.NewGiftsHandler(giftService),
		mw:   api.NewAuthMiddleware(authService),
	}, store, nil
}

func registerPublicRoutes(h *handlers) {
	http.HandleFunc("POST /api/register", h.auth.Register)
	http.HandleFunc("POST /api/login", h.auth.Login)
	http.HandleFunc("GET /api/wishlists/public", h.wish.GetWishlistByToken)
	http.HandleFunc("POST /api/gifts/reserve", h.gift.ReserveGift)
	http.HandleFunc("GET /api/gifts", h.gift.GetGiftsByWishlist)
}

func registerProtectedRoutes(h *handlers) {
	http.HandleFunc("POST /api/wishlists", h.mw.Authenticate(h.wish.CreateWishlist))
	http.HandleFunc("GET /api/wishlists", h.mw.Authenticate(h.wish.GetWishlist))
	http.HandleFunc("GET /api/wishlists/get", h.mw.Authenticate(h.wish.GetWishlistByID))
	http.HandleFunc("PUT /api/wishlists/update", h.mw.Authenticate(h.wish.UpdateWishlist))
	http.HandleFunc("DELETE /api/wishlists/delete", h.mw.Authenticate(h.wish.DeleteWishlist))

	http.HandleFunc("POST /api/gifts", h.mw.Authenticate(h.gift.CreateGift))
	http.HandleFunc("DELETE /api/gifts", h.mw.Authenticate(h.gift.DeleteGift))
}

func registerSwaggerRoutes() {
	http.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)
	http.HandleFunc("GET /swagger/*", httpSwagger.WrapHandler)
}

func startServer(port string) error {
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Printf("Swagger UI: http://localhost:%s/swagger/index.html\n", port)
	return http.ListenAndServe(":"+port, nil)
}

func main() {
	cfg := loadConfig()

	h, store, err := initHandlers(cfg)
	if err != nil {
		fmt.Println("Failed to initialize handlers:", err)
		return
	}
	defer store.Close()

	registerPublicRoutes(h)
	registerProtectedRoutes(h)
	registerSwaggerRoutes()

	if err := startServer(cfg.port); err != nil {
		fmt.Println("Server error:", err)
		return
	}
}
