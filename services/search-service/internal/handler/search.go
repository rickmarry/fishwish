package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fishwish/services/search-service/internal/model"
	"fishwish/services/search-service/internal/service"
)

type SearchHandler struct {
	service *service.SearchService
}

func NewSearchHandler(svc *service.SearchService) *SearchHandler {
	return &SearchHandler{service: svc}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	params := model.SearchParams{
		Query: r.URL.Query().Get("q"),
		Type:  r.URL.Query().Get("type"),
		Species: r.URL.Query().Get("species"),
		State: r.URL.Query().Get("state"),
		Limit: 20,
	}

	if l := r.URL.Query().Get("limit"); l != "" {
		params.Limit, _ = strconv.Atoi(l)
	}

	results, err := h.service.Search(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, results)
}

func (h *SearchHandler) Species(w http.ResponseWriter, r *http.Request) {
	species, err := h.service.GetSpecies(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, species)
}

func (h *SearchHandler) Suggestions(w http.ResponseWriter, r *http.Request) {
	prefix := r.URL.Query().Get("q")
	if prefix == "" {
		writeJSON(w, http.StatusOK, []string{})
		return
	}

	suggestions, err := h.service.GetSuggestions(r.Context(), prefix)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, suggestions)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
