package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/imdevinc/go-links/internal/config"
	"github.com/imdevinc/go-links/internal/store"
)

type App struct {
	Store  store.Store
	Logger *slog.Logger
	config *config.Config
}

func (a *App) Start(ctx context.Context, cfg *config.Config) error {
	if a.Store == nil {
		return fmt.Errorf("store is nil")
	}

	if cfg.StaticPath == "" {
		cfg.StaticPath = "/"
	}
	a.config = cfg

	r := mux.NewRouter()
	r.Use(corsHandler)

	r.HandleFunc("/api/popular", a.getPopular)
	fs := http.FileServer(http.Dir(cfg.StaticPath))
	r.Path("/").Handler(fs)
	r.PathPrefix("/static/").Handler(fs)
	r.HandleFunc("/{link:.*}", a.handleLink)

	a.Logger.Info("starting server on port 8080")
	return http.ListenAndServe(fmt.Sprintf(":%d", 8080), r)
}

func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")
		next.ServeHTTP(w, r)
	})
}

func (a *App) handleLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case http.MethodGet:
		a.handleGetLink(w, r)
	case http.MethodPost:
		a.handleCreateLink(w, r)
	case http.MethodDelete:
		a.handleDeleteLink(w, r)
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *App) handleGetLink(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	if _, ok := v["link"]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "missing link"}`))
		return
	}

	link, err := a.Store.GetLinkByName(r.Context(), v["link"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "link not found"}`))
		return
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}

func (a *App) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		a.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}
	link, err := store.CreateLinkFromPayload(body)
	if err != nil {
		a.Logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid payload"}`))
		return
	}
	link.Name = mux.Vars(r)["link"]

	err = a.Store.CreateLink(r.Context(), link)
	if err != nil {
		a.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *App) handleDeleteLink(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["link"]
	err := a.Store.DisableLink(r.Context(), name)
	if err != nil {
		a.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (a *App) getPopular(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	links, err := a.Store.GetPopularLinks(r.Context(), 10)
	if err != nil {
		a.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}
	err = json.NewEncoder(w).Encode(links)
	if err != nil {
		a.Logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}
}
