package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"fishwish/services/spot-service/internal/model"
	"fishwish/services/spot-service/internal/service"
)

type SpotHandler struct {
	service *service.SpotService
}

func NewSpotHandler(svc *service.SpotService) *SpotHandler {
	return &SpotHandler{service: svc}
}

func (h *SpotHandler) ListSpots(w http.ResponseWriter, r *http.Request) {
	params := model.ListSpotsParams{
		Type:      r.URL.Query().Get("type"),
		Species:   r.URL.Query().Get("species"),
		Difficulty: r.URL.Query().Get("difficulty"),
		Page:      1,
		Limit:     20,
	}

	if p := r.URL.Query().Get("page"); p != "" {
		params.Page, _ = strconv.Atoi(p)
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		params.Limit, _ = strconv.Atoi(l)
	}

	spots, err := h.service.ListSpots(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, spots)
}

func (h *SpotHandler) GetSpot(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	spot, err := h.service.GetSpot(r.Context(), id)
	if err != nil {
		http.Error(w, "spot not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, spot)
}

func (h *SpotHandler) GetSpotDetails(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	details, err := h.service.GetSpotDetails(r.Context(), id)
	if err != nil {
		http.Error(w, "spot not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, details)
}

func (h *SpotHandler) NearbySpots(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
	radiusMi := 25.0
	if r := r.URL.Query().Get("radius"); r != "" {
		radiusMi, _ = strconv.ParseFloat(r, 64)
	}

	params := model.NearbyParams{
		Lat:      lat,
		Lon:      lon,
		RadiusMi: radiusMi,
		Limit:    20,
	}

	spots, err := h.service.NearbySpots(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, spots)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
