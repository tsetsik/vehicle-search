package http

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/tsetsik/vehicle-search/internal/config"
	"github.com/tsetsik/vehicle-search/internal/core"
)

type Handler struct {
	seSvc    core.SearchEngine
	listings core.Listings
	ctx      context.Context
	cfg      config.Config
	logger   *slog.Logger
}

func NewHttpResolver(ctx context.Context, logger *slog.Logger, cfg config.Config, store core.Store) *Handler {
	return &Handler{
		seSvc:    core.NewSearchEngine(store),
		listings: core.NewListings(store),
		ctx:      ctx,
		cfg:      cfg,
		logger:   logger,
	}
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	search := r.Header.Get("q")

	filter := core.ItemsFilter{
		Q: search,
	}
	results, err := h.seSvc.Search(filter)
	if err != nil {
		http.Error(w, "Search error", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]any{"items": results})
	if err != nil {
		http.Error(w, "Search error", http.StatusInternalServerError)
	}
	defer r.Body.Close() //nolint:errcheck
	h.responde(w, http.StatusOK, response)
}

func (h *Handler) Listings(w http.ResponseWriter, r *http.Request) {
	raw := json.RawMessage(r.Header.Get("item"))
	item := core.Item{}
	if err := json.Unmarshal(raw, &item); err != nil {
		h.logger.Error("Invalid listing", slog.Any("error", err))
		h.responde(w, http.StatusBadRequest, []byte("Invalid listing"))
		return
	}

	if err := item.Validate(); err != nil {
		h.logger.Error("Invalid listing", slog.Any("error", err))
		h.responde(w, http.StatusBadRequest, []byte("Invalid listing"))
		return
	}

	if err := h.listings.AddItem(item); err != nil {
		h.responde(w, http.StatusBadRequest, []byte("Invalid listing"))
	}

	h.logger.Info("New listing added", slog.Any("items", item))
	h.responde(w, http.StatusOK, []byte("Request received"))
}

func (h *Handler) responde(w http.ResponseWriter, status int, body []byte) {
	w.WriteHeader(status)
	//nolint:errcheck
	w.Write(body)
}
