package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"fishwish/services/social-service/internal/model"
	"fishwish/services/social-service/internal/service"
)

type SocialHandler struct {
	service *service.SocialService
}

func NewSocialHandler(svc *service.SocialService) *SocialHandler {
	return &SocialHandler{service: svc}
}

func (h *SocialHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	spotID := chi.URLParam(r, "spotID")

	reviews, err := h.service.ListReviews(r.Context(), spotID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, reviews)
}

func (h *SocialHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	spotID := chi.URLParam(r, "spotID")
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	review, err := h.service.CreateReview(r.Context(), userID, spotID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, review)
}

func (h *SocialHandler) LogCatch(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.CreateCatchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	catch, err := h.service.LogCatch(r.Context(), userID, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, catch)
}

func (h *SocialHandler) GetUserCatches(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	catches, err := h.service.GetUserCatches(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, catches)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
