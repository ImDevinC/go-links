package store

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

type Link struct {
	Name        string    `json:"name" bson:"_id"`
	Description string    `json:"description" bson:"description"`
	URL         string    `json:"url" bson:"url"`
	Views       int       `json:"views" bson:"views"`
	Created     time.Time `json:"created_at" bson:"created_at"`
	Updated     time.Time `json:"updated_at" bson:"updated_at"`
	CreatedBy   string    `json:"created_by" bson:"created_by"`
	Disabled    bool      `json:"disabled" bson:"disabled"`
}

func CreateLinkFromPayload(payload []byte) (Link, error) {
	var link Link
	err := json.Unmarshal(payload, &link)
	if err != nil {
		return Link{}, err
	}
	return link, nil
}

var ErrIDExists = errors.New("id exists")
var ErrLinkNotFound = errors.New("link not found")

type Store interface {
	CreateLink(ctx context.Context, link Link) error
	GetLinkByName(ctx context.Context, name string) (Link, error)
	GetLinkByURL(ctx context.Context, url string) (Link, error)
	DisableLink(ctx context.Context, name string) error
	GetPopularLinks(ctx context.Context, size int) ([]Link, error)
	GetRecentLinks(ctx context.Context, size int) ([]Link, error)
	GetOwnedLinks(ctx context.Context, email string) ([]Link, error)
	IncrementLinkViews(ctx context.Context, name string) error
	QueryLinks(ctx context.Context, query string) ([]Link, error)
	Close(ctx context.Context) error
}
