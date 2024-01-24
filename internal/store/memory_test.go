package store_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/imdevinc/go-links/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestMemoryGetPopular(t *testing.T) {
	ctx := context.Background()
	m := store.NewMemoryStore()
	m.CreateLink(ctx, store.Link{
		Name:  "test",
		Views: 10,
	})
	m.CreateLink(ctx, store.Link{
		Name:  "test2",
		Views: 5,
	})
	m.CreateLink(ctx, store.Link{
		Name:  "test3",
		Views: 15,
	})
	for i := range [10]int{} {
		m.CreateLink(ctx, store.Link{
			Name:  fmt.Sprintf("test%d", i),
			Views: 1,
		})
	}
	links, _ := m.GetPopularLinks(ctx, 3)
	if !assert.Equal(t, 3, len(links)) {
		t.FailNow()
	}
	expected := []store.Link{
		{
			Name:  "test3",
			Views: 15,
		},
		{
			Name:  "test",
			Views: 10,
		},
		{
			Name:  "test2",
			Views: 5,
		},
	}
	diff := cmp.Diff(links, expected)
	if !assert.Equal(t, "", diff) {
		t.FailNow()
	}
}
