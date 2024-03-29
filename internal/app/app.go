package app

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/crewjam/saml/samlsp"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/imdevinc/go-links/internal/config"
	"github.com/imdevinc/go-links/internal/store"
)

type App struct {
	Store  store.Store
	Logger *slog.Logger
	config *config.Config
	sp     *samlsp.Middleware
}

type GetLinksType string

const (
	Recent  GetLinksType = "recent"
	Popular GetLinksType = "popular"
	Owned   GetLinksType = "owned"
)

var staticRegexp = regexp.MustCompile(`^static(\/.*)?|api\/popular|api\/recent`)
var validLinkRegexp = regexp.MustCompile(`^[a-zA-Z0-9\/\-]*$`)

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

type QueryInput struct {
	Query string `json:"query"`
}

func (a *App) Start(ctx context.Context, cfg *config.Config) error {
	if a.Store == nil {
		return fmt.Errorf("store is nil")
	}

	if cfg.StaticPath == "" {
		cfg.StaticPath = "/"
	}
	a.config = cfg

	sp, err := a.configureSaml()
	if err != nil && a.config.SSO.Require {
		return fmt.Errorf("failed to configure saml: %w", err)
	}
	a.sp = sp
	authWrapper := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if sp == nil {
				next.ServeHTTP(w, r)
			} else {
				sp.RequireAccount(next).ServeHTTP(w, r)
			}
		})
	}

	r := mux.NewRouter()
	r.Use(corsHandler)
	r.Use(a.indexHandler)

	fs := http.FileServer(http.Dir(cfg.StaticPath))

	if sp != nil {
		r.PathPrefix("/saml/").Handler(sp)
	}
	r.PathPrefix("/api").Handler(authWrapper(http.HandlerFunc(a.handleApi)))
	r.Path("/").Handler(authWrapper(fs))
	r.PathPrefix("/static").Methods(http.MethodGet).Handler(authWrapper(fs))
	assets := []string{
		"/asset-manifest.json",
		"/favicon.ico",
		"/logo192.png",
		"/logo512.png",
		"/manifest.json",
		"/robots.txt",
	}
	for _, a := range assets {
		r.Path(a).Handler(authWrapper(fs))
	}
	r.PathPrefix("/static").Handler(authWrapper(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sendError(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "Protected route"})
	})))
	r.Path("/{link:.*}").Handler(authWrapper(http.HandlerFunc(a.handleLink)))
	a.Logger.With("port", cfg.Port).Info("starting server")
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), r)
}

func (a *App) indexHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == a.config.FQDN {
			next.ServeHTTP(w, r)
		} else {
			w.Header().Set("Location", fmt.Sprintf("//%s%s", a.config.FQDN, r.URL.Path))
			w.WriteHeader(http.StatusMovedPermanently)
		}
	})
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
		link, err := cleanLink(link)
		if err != nil {
			sendError(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
			return
		}
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
	link, err := cleanLink(mux.Vars(r)["link"])
	if err != nil {
		sendError(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	result, err := a.Store.GetLinkByName(r.Context(), link)
	if err != nil {
		if !errors.Is(err, store.ErrLinkNotFound) {
			a.Logger.Error(err.Error())
		}
		w.Header().Set("Location", fmt.Sprintf("//%s", a.config.FQDN))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}
	err = a.Store.IncrementLinkViews(r.Context(), result.Name)
	if err != nil {
		a.Logger.Error(err.Error())
	}
	http.Redirect(w, r, result.URL, http.StatusFound)
}

func (a *App) handleCreateLink(w http.ResponseWriter, r *http.Request) {
	email, err := a.getEmailFromRequest(r)
	if err != nil {
		a.Logger.Error(err.Error())
		sendError(w, http.StatusUnauthorized, ErrorResponse{Error: "missing authentication token"})
		return
	}
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
	clean, err := cleanLink(mux.Vars(r)["link"])
	if err != nil {
		sendError(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
	}
	link.Name = clean
	link.CreatedBy = email
	err = a.Store.CreateLink(r.Context(), link)
	if err == store.ErrIDExists {
		sendError(w, http.StatusConflict, ErrorResponse{Error: "link already exists"})
		return
	}
	if err != nil {
		a.Logger.Error(err.Error())
		sendError(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *App) handleDeleteLink(w http.ResponseWriter, r *http.Request) {
	name, err := cleanLink(mux.Vars(r)["link"])
	if err != nil {
		sendError(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}
	err = a.Store.DisableLink(r.Context(), name)
	if err != nil {
		a.Logger.Error(err.Error())
		sendError(w, http.StatusInternalServerError, ErrorResponse{Error: "internal server error"})
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

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
		case Owned:
			email, err := a.getEmailFromRequest(r)
			if err != nil {
				a.Logger.Error(err.Error())
				sendError(w, http.StatusUnauthorized, ErrorResponse{Error: "missing authentication token"})
				return
			}
			links, err = a.Store.GetOwnedLinks(r.Context(), email)
			if err != nil {
				a.Logger.Error(err.Error())
				sendError(w, http.StatusUnauthorized, ErrorResponse{Error: "internal server error"})
				return
			}
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

func cleanLink(link string) (string, error) {
	if len(link) == 0 {
		return "", fmt.Errorf("link is empty")
	}
	if link[0] == '/' {
		link = link[1:]
	}
	if link[len(link)-1] == '/' {
		link = link[:len(link)-1]
	}
	if link[0] == '-' {
		return "", fmt.Errorf("name input is invalid")
	}
	if link[len(link)-1] == '-' {
		return "", fmt.Errorf("name input is invalid")
	}
	if !validLinkRegexp.Match([]byte(link)) {
		return "", fmt.Errorf("name input is invalid")
	}
	return strings.ToLower(link), nil
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
	case "/api/owned":
		a.handleGetLinkList(Owned)(w, r)
	case "/api/query":
		a.handleQueryLinks(w, r)
	default:
		sendError(w, http.StatusNotFound, ErrorResponse{Error: "not found"})
	}
}

func (a *App) configureSaml() (*samlsp.Middleware, error) {
	keypair, err := tls.X509KeyPair(a.config.SSO.SamlCert, a.config.SSO.SamlKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load keypair: %w", err)
	}
	keypair.Leaf, err = x509.ParseCertificate(keypair.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}
	metadataContent, err := os.ReadFile(a.config.SSO.MetadataFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata file: %w", err)
	}
	idpMetadata, err := samlsp.ParseMetadata(metadataContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	rootUrl, err := url.Parse(a.config.SSO.CallbackURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse root url: %w", err)
	}
	sp, err := samlsp.New(samlsp.Options{
		URL:               *rootUrl,
		Key:               keypair.PrivateKey.(*rsa.PrivateKey),
		Certificate:       keypair.Leaf,
		IDPMetadata:       idpMetadata,
		EntityID:          a.config.SSO.EntityID,
		AllowIDPInitiated: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create samlsp: %w", err)
	}
	return sp, nil
}

func (a *App) getEmailFromRequest(r *http.Request) (string, error) {
	if a.sp == nil {
		return "untracked", nil
	}
	c, err := r.Cookie("token")
	if err != nil {
		return "", fmt.Errorf("failed to get cookie: %w", err)
	}
	token, _ := jwt.Parse(c.Value, func(t *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to parse claims")
	}
	return claims.GetSubject()
}

func (a *App) handleQueryLinks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, ErrorResponse{})
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.Logger.Error(err.Error())
		sendError(w, http.StatusBadRequest, ErrorResponse{Error: "bad request"})
		return
	}
	input := QueryInput{}
	err = json.Unmarshal(body, &input)
	if err != nil {
		a.Logger.Error(err.Error())
		sendError(w, http.StatusBadRequest, ErrorResponse{Error: "bad request"})
		return
	}

	links, err := a.Store.QueryLinks(r.Context(), input.Query)
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
