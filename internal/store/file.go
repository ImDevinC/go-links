package store

import (
	"cmp"
	"context"
	"encoding/json"
	"io"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

type file struct {
	path  string
	links map[string]Link
	mu    sync.Mutex
}

var _ Store = (*file)(nil)

func NewFileStore(path string, createFile bool) (*file, error) {
	flags := os.O_RDWR
	if createFile {
		flags = os.O_RDWR | os.O_CREATE | os.O_APPEND
	}
	f, err := os.OpenFile(path, flags, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	existing, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	links := map[string]Link{}
	if len(existing) > 0 {
		err = json.Unmarshal(existing, &links)
		if err != nil {
			return nil, err
		}
	}

	return &file{
		path:  path,
		links: links,
		mu:    sync.Mutex{},
	}, nil
}

func (f *file) saveLinks() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	data, err := json.Marshal(f.links)
	if err != nil {
		return err
	}
	return os.WriteFile(f.path, data, os.ModePerm)
}

// CreateLink implements Store.
func (f *file) CreateLink(ctx context.Context, link Link) error {
	_, err := f.GetLinkByName(ctx, link.Name)
	if err != nil && err != ErrLinkNotFound {
		return err
	}
	link.Created = time.Now()
	link.Updated = link.Created
	f.links[link.Name] = link
	return f.saveLinks()
}

// DisableLink implements Store.
func (f *file) DisableLink(ctx context.Context, name string) error {
	delete(f.links, name)
	return f.saveLinks()
}

// GetLinkByName implements Store.
func (f *file) GetLinkByName(ctx context.Context, name string) (Link, error) {
	for _, link := range f.links {
		if !link.Disabled && strings.EqualFold(link.Name, name) {
			return link, nil
		}
	}
	return Link{}, ErrLinkNotFound
}

// GetLinkByURL implements Store.
func (f *file) GetLinkByURL(ctx context.Context, url string) (Link, error) {
	for _, link := range f.links {
		if !link.Disabled && strings.EqualFold(link.URL, url) {
			return link, nil
		}
	}
	return Link{}, ErrLinkNotFound
}

// GetOwnedLinks implements Store.
func (f *file) GetOwnedLinks(ctx context.Context, email string) ([]Link, error) {
	links := []Link{}
	for _, link := range f.links {
		if !link.Disabled && strings.EqualFold(link.CreatedBy, email) {
			links = append(links, link)
		}
	}
	return links, nil
}

// GetPopularLinks implements Store.
func (f *file) GetPopularLinks(ctx context.Context, size int) ([]Link, error) {
	links := []Link{}
	for _, link := range f.links {
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

// GetRecentLinks implements Store.
func (f *file) GetRecentLinks(ctx context.Context, size int) ([]Link, error) {
	links := []Link{}
	for _, link := range f.links {
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

// IncrementLinkViews implements Store.
func (f *file) IncrementLinkViews(ctx context.Context, name string) error {
	if _, ok := f.links[name]; !ok {
		return ErrLinkNotFound
	}
	link := f.links[name]
	link.Views++
	f.links[name] = link
	return f.saveLinks()
}

// QueryLinks implements Store.
func (f *file) QueryLinks(ctx context.Context, query string) ([]Link, error) {
	links := []Link{}
	for _, link := range f.links {
		if !link.Disabled && (strings.Contains(strings.ToLower(link.Name), strings.ToLower(query)) || strings.Contains(strings.ToLower(link.Description), strings.ToLower(query))) {
			links = append(links, link)
		}
	}
	return links, nil
}

// Close implements Store.
func (*file) Close(ctx context.Context) error {
	return nil
}
