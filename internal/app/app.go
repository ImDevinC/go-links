package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/imdevinc/go-links/internal/config"
	"github.com/imdevinc/go-links/internal/store"
)

type App struct {
	Store  store.Store
	Logger *slog.Logger
	config *config.Config
}

type GetLinksType string

const (
	Recent  GetLinksType = "recent"
	Popular GetLinksType = "popular"
)

var staticRegexp = regexp.MustCompile(`^static(\/.*)?|api\/popular|api\/recent`)

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
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

	r.PathPrefix("/api").Methods(http.MethodGet).HandlerFunc(a.handleApi)
	r.PathPrefix("/api").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendError(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "Protected route"})
	})
	fs := http.FileServer(http.Dir(cfg.StaticPath))
	r.Path("/").Handler(fs)
	r.PathPrefix("/static").Methods(http.MethodGet).Handler(fs)
	r.PathPrefix("/static").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendError(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "Protected route"})
	})
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
	// static/ is protected due to web display resources
	v := mux.Vars(r)
	if link, ok := v["link"]; ok {
		link = cleanLink(link)
		if staticRegexp.Match([]byte(link)) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}
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
	result, err := a.Store.GetLinkByName(r.Context(), cleanLink(mux.Vars(r)["link"]))
	if err != nil {
		sendError(w, http.StatusNotFound, ErrorResponse{Error: "link not found"})
		return
	}

	http.Redirect(w, r, result.URL, http.StatusFound)
}

func (a *App) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		a.Logger.Error(err.Error())
		sendError(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	link, err := store.CreateLinkFromPayload(body)
	if err != nil {
		a.Logger.Error(err.Error())
		sendError(w, http.StatusBadRequest, ErrorResponse{Error: "invalid payload"})
		return
	}
	link.Name = cleanLink(mux.Vars(r)["link"])
	err = a.Store.CreateLink(r.Context(), link)
	if err != nil {
		a.Logger.Error(err.Error())
		sendError(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *App) handleDeleteLink(w http.ResponseWriter, r *http.Request) {
	name := cleanLink(mux.Vars(r)["link"])
	err := a.Store.DisableLink(r.Context(), name)
	if err != nil {
		a.Logger.Error(err.Error())
		sendError(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

// func (a *App) handleGetLinkList(w http.ResponseWriter, r *http.Request) {
func (a *App) handleGetLinkList(t GetLinksType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet {
			sendError(w, http.StatusMethodNotAllowed, ErrorResponse{})
			return
		}
		var links []store.Link
		var err error
		switch t {
		case Popular:
			links, err = a.Store.GetPopularLinks(r.Context(), 10)
		case Recent:
			links, err = a.Store.GetRecentLinks(r.Context(), 10)
		default:
			sendError(w, http.StatusBadRequest, ErrorResponse{Error: "unknown link list type"})
			return
		}

		if err != nil {
			a.Logger.Error(err.Error())
			sendError(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
			return
		}
		err = json.NewEncoder(w).Encode(links)
		if err != nil {
			a.Logger.Error(err.Error())
			sendError(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
			return
		}
	}
}

func cleanLink(link string) string {
	if len(link) == 0 {
		return link
	}
	if link[0] == '/' {
		link = link[1:]
	}
	if link[len(link)-1] == '/' {
		link = link[:len(link)-1]
	}
	return link
}

func sendError(w http.ResponseWriter, code int, message ErrorResponse) {
	w.WriteHeader(code)
	if message.Error != "" {
		json.NewEncoder(w).Encode(message)
	}
}

func (a *App) handleApi(w http.ResponseWriter, r *http.Request) {
	switch strings.ToLower(r.URL.Path) {
	case "/api/popular":
		a.handleGetLinkList(Popular)(w, r)
	case "/api/recent":
		a.handleGetLinkList(Recent)(w, r)
	default:
		sendError(w, http.StatusNotFound, ErrorResponse{Error: "not found"})
	}
}
