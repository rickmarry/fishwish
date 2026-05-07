package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"fishwish/services/user-service/internal/config"
	"fishwish/services/user-service/internal/handler"
	"fishwish/services/user-service/internal/repository"
	"fishwish/services/user-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := repository.NewDB(cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	userService := service.NewUserService(db)
	userHandler := handler.NewUserHandler(userService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Post("/auth/register", userHandler.Register)
	r.Post("/auth/login", userHandler.Login)
	r.Post("/auth/refresh", userHandler.Refresh)

	r.Group(func(r chi.Router) {
		r.Get("/users/me", userHandler.GetProfile)
		r.Put("/users/me", userHandler.UpdateProfile)
		r.Get("/users/me/preferences", userHandler.GetPreferences)
		r.Put("/users/me/preferences", userHandler.UpdatePreferences)
		r.Get("/users/me/saved-spots", userHandler.GetSavedSpots)
		r.Post("/users/me/saved-spots/{spotID}", userHandler.SaveSpot)
		r.Delete("/users/me/saved-spots/{spotID}", userHandler.UnsaveSpot)
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	go func() {
		log.Printf("user-service starting on %s", cfg.Port)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("user-service stopped")
}
