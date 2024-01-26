package store

import (
	"cmp"
	"context"
	"slices"
	"strings"
	"time"
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
	delete(m.links, name)
	// if _, ok := m.links[name]; !ok {
	// 	return ErrLinkNotFound
	// }
	// link := m.links[name]
	// link.Updated = time.Now()
	// // link.Disabled = true

	// m.links[name] = link
	return nil
}

// CreateLink implements Store.
func (m *memory) CreateLink(ctx context.Context, link Link) error {
	if _, ok := m.links[link.Name]; ok {
		return ErrIDExists
	}
	link.Created = time.Now()
	link.Updated = link.Created
	m.links[link.Name] = link
	return nil
}

// GetLinkByName implements Store.
func (m *memory) GetLinkByName(ctx context.Context, name string) (Link, error) {
	for _, link := range m.links {
		if !link.Disabled && strings.EqualFold(link.Name, name) {
			return link, nil
		}
	}
	return Link{}, ErrLinkNotFound
}

// GetLinkByURL implements Store.
func (m *memory) GetLinkByURL(ctx context.Context, url string) (Link, error) {
	for _, link := range m.links {
		if !link.Disabled && strings.EqualFold(link.URL, url) {
			return link, nil
		}
	}
	return Link{}, ErrLinkNotFound
}

func (m *memory) GetPopularLinks(ctx context.Context, size int) ([]Link, error) {
	links := []Link{}
	for _, link := range m.links {
		if !link.Disabled {
			links = append(links, link)
		}
	}
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
	for _, link := range m.links {
		if !link.Disabled {
			links = append(links, link)
		}
	}
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
	for _, link := range m.links {
		if !link.Disabled && strings.EqualFold(link.CreatedBy, email) {
			links = append(links, link)
		}
	}
	return links, nil
}

func (m *memory) IncrementLinkViews(ctx context.Context, name string) error {
	if _, ok := m.links[name]; !ok {
		return ErrLinkNotFound
	}
	link := m.links[name]
	link.Views++
	m.links[name] = link
	return nil
}
