package store

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"sync"
	"time"
)

type memory struct {
	links sync.Map
}

var _ Store = (*memory)(nil)

func NewMemoryStore() *memory {
	return &memory{
		links: sync.Map{},
	}
}

// DeleteLink implements Store.
func (m *memory) DisableLink(ctx context.Context, name string) error {
	m.links.Delete(name)
	return nil
}

// CreateLink implements Store.
func (m *memory) CreateLink(ctx context.Context, link Link) error {
	if _, ok := m.links.Load(link.Name); ok {
		return ErrIDExists
	}
	link.Created = time.Now()
	link.Updated = link.Created
	m.links.Store(link.Name, link)
	return nil
}

// GetLinkByName implements Store.
func (m *memory) GetLinkByName(ctx context.Context, name string) (Link, error) {
	var link Link
	m.links.Range(func(key, value any) bool {
		l := value.(Link)
		if !l.Disabled && strings.EqualFold(l.Name, name) {
			link = l
			return false
		}
		return true
	})
	if link.Name == "" {
		return Link{}, ErrLinkNotFound
	}
	return link, nil
}

// GetLinkByURL implements Store.
func (m *memory) GetLinkByURL(ctx context.Context, url string) (Link, error) {
	var link Link
	m.links.Range(func(key, value any) bool {
		l := value.(Link)
		if !l.Disabled && strings.EqualFold(l.URL, url) {
			link = l
			return false
		}
		return true
	})
	if link.Name == "" {
		return Link{}, ErrLinkNotFound
	}
	return link, nil
}

func (m *memory) GetPopularLinks(ctx context.Context, size int) ([]Link, error) {
	links := []Link{}
	m.links.Range(func(key, value any) bool {
		l := value.(Link)
		if !l.Disabled {
			links = append(links, l)
		}
		return true
	})
	slices.SortFunc(links, func(a Link, b Link) int {
		return cmp.Compare(b.Views, a.Views)
	})
	if size < len(links) {
		links = links[:size]
	}
	return links, nil
}

func (m *memory) GetRecentLinks(ctx context.Context, size int) ([]Link, error) {
	links := []Link{}
	m.links.Range(func(key, value any) bool {
		l := value.(Link)
		if !l.Disabled {
			links = append(links, l)
		}
		return true
	})
	slices.SortFunc(links, func(a Link, b Link) int {
		return b.Updated.Compare(a.Updated)
	})
	if size < len(links) {
		links = links[:size]
	}
	return links, nil
}

func (m *memory) GetOwnedLinks(ctx context.Context, email string) ([]Link, error) {
	links := []Link{}
	m.links.Range(func(key, value any) bool {
		l := value.(Link)
		if !l.Disabled {
			links = append(links, l)
		}
		return true
	})
	return links, nil
}

func (m *memory) IncrementLinkViews(ctx context.Context, name string) error {
	l, ok := m.links.Load(name)
	if !ok {
		return ErrLinkNotFound
	}
	link := l.(Link)
	link.Views++
	m.links.Store(name, link)
	return nil
}

func (m *memory) QueryLinks(ctx context.Context, query string) ([]Link, error) {
	links := []Link{}
	m.links.Range(func(key, value any) bool {
		link := value.(Link)
		if !link.Disabled && (strings.Contains(strings.ToLower(link.Name), strings.ToLower(query)) || strings.Contains(strings.ToLower(link.Description), strings.ToLower(query))) {
			links = append(links, link)
		}
		return true
	})
	return links, nil
}

// Close implements Store.
func (*memory) Close(ctx context.Context) error {
	return nil
}
