package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fishwish/services/weather-service/internal/service"
)

type WeatherHandler struct {
	service *service.WeatherService
}

func NewWeatherHandler(svc *service.WeatherService) *WeatherHandler {
	return &WeatherHandler{service: svc}
}

func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	lat, errLat := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lng, errLng := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)

	if errLat != nil || errLng != nil || lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		http.Error(w, `{"error":"invalid lat/lng parameters","message":"lat (-90 to 90) and lng (-180 to 180) are required"}`, http.StatusBadRequest)
		return
	}

	resp, err := h.service.GetWeather(r.Context(), lat, lng)
	if err != nil {
		http.Error(w, `{"error":"weather API unavailable","message":"`+err.Error()+`"}`, http.StatusBadGateway)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *WeatherHandler) GetForecast(w http.ResponseWriter, r *http.Request) {
	lat, errLat := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lng, errLng := strconv.ParseFloat(r.URL.Query().Get("lng"), 64)

	if errLat != nil || errLng != nil || lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		http.Error(w, `{"error":"invalid lat/lng parameters"}`, http.StatusBadRequest)
		return
	}

	forecast, err := h.service.GetForecast(r.Context(), lat, lng)
	if err != nil {
		http.Error(w, `{"error":"weather API unavailable"}`, http.StatusBadGateway)
		return
	}

	writeJSON(w, http.StatusOK, forecast)
}

func (h *WeatherHandler) GetTides(w http.ResponseWriter, r *http.Request) {
	lat, _ := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	lon, _ := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)

	if lat == 0 && lon == 0 {
		http.Error(w, "lat and lon are required", http.StatusBadRequest)
		return
	}

	tides, err := h.service.GetTides(r.Context(), lat, lon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, tides)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
