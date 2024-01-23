package store

import (
	"context"
	"strings"
)

type memory struct {
	links map[string]Link
}

var _ Store = (*memory)(nil)

func NewMemoryStore() *memory {
	return &memory{
		links: make(map[string]Link),
	}
}

// DeleteLink implements Store.
func (m *memory) DisableLink(ctx context.Context, name string) error {
	if _, ok := m.links[name]; !ok {
		return ErrLinkNotFound
	}
	link := m.links[name]
	link.Disabled = true
	m.links[name] = link
	return nil
}

// CreateLink implements Store.
func (m *memory) CreateLink(ctx context.Context, link Link) error {
	if _, ok := m.links[link.Name]; ok {
		return ErrIDExists
	}
	m.links[link.Name] = link
	return nil
}

// GetLinkByName implements Store.
func (m *memory) GetLinkByName(ctx context.Context, name string) (Link, error) {
	for _, link := range m.links {
		if strings.EqualFold(link.Name, name) && !link.Disabled {
			return link, nil
		}
	}
	return Link{}, ErrLinkNotFound
}

// GetLinkByURL implements Store.
func (m *memory) GetLinkByURL(ctx context.Context, url string) (Link, error) {
	for _, link := range m.links {
		if strings.EqualFold(link.URL, url) {
			return link, nil
		}
	}
	return Link{}, ErrLinkNotFound
}
