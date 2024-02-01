package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type postgres struct {
	conn *pgx.Conn
}

var _ (Store) = (*postgres)(nil)

func NewPostgresStore(ctx context.Context, user string, password string, host string, databaseName string) (*postgres, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s/%s", user, password, host, databaseName)
	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	p := &postgres{
		conn: conn,
	}
	err = p.ensureTable(ctx)
	if err != nil {
		p.Close(ctx)
		return nil, err
	}
	return p, nil
}

// Close implements Store.
func (p *postgres) Close(ctx context.Context) error {
	return p.conn.Close(ctx)
}

// CreateLink implements Store.
func (p *postgres) CreateLink(ctx context.Context, link Link) error {
	_, err := p.conn.Exec(ctx,
		`insert into links(name, description, url, created_at, updated_at, created_by) values ($1, $2, $3, $4, $5, $6)`,
		link.Name, link.Description, link.URL, link.Created, link.Created, link.CreatedBy,
	)
	if err != nil {
		return err
	}
	return nil
}

// DisableLink implements Store.
func (p *postgres) DisableLink(ctx context.Context, name string) error {
	resp, err := p.conn.Exec(ctx, `delete from links where name = $1`, name)
	if err != nil {
		return err
	}
	if resp.RowsAffected() == 0 {
		return ErrLinkNotFound
	}
	return nil
}

// GetLinkByName implements Store.
func (p *postgres) GetLinkByName(ctx context.Context, name string) (Link, error) {
	return p.getSingleResult(ctx, `select * from links where name = $1`, name)
}

// GetLinkByURL implements Store.
func (p *postgres) GetLinkByURL(ctx context.Context, url string) (Link, error) {
	return p.getSingleResult(ctx, `select * from links where url = $1`, url)
}

// GetOwnedLinks implements Store.
func (p *postgres) GetOwnedLinks(ctx context.Context, email string) ([]Link, error) {
	return p.getMultipleResults(ctx, `select * from links where created_by = $1`, email)
}

// GetPopularLinks implements Store.
func (p *postgres) GetPopularLinks(ctx context.Context, size int) ([]Link, error) {
	return p.getMultipleResults(ctx, fmt.Sprintf("select * from links order by views desc limit %d", size))
}

// GetRecentLinks implements Store.
func (p *postgres) GetRecentLinks(ctx context.Context, size int) ([]Link, error) {
	return p.getMultipleResults(ctx, fmt.Sprintf("select * from links order by updated_at desc limit %d", size))
}

// IncrementLinkViews implements Store.
func (p *postgres) IncrementLinkViews(ctx context.Context, name string) error {
	link, err := p.GetLinkByName(ctx, name)
	if err != nil {
		return err
	}
	views := link.Views + 1
	resp, err := p.conn.Exec(ctx, `update links set views=$1 where name=$2`, views, name)
	if err != nil {
		return err
	}
	if resp.RowsAffected() == 0 {
		return ErrLinkNotFound
	}
	return nil
}

// QueryLinks implements Store.
func (p *postgres) QueryLinks(ctx context.Context, query string) ([]Link, error) {
	return p.getMultipleResults(ctx, `select * from links where name ilike '%' || $1 || '%' or description ilike '%' || $1 || '%' order by views desc`, query)
}

func (p *postgres) getSingleResult(ctx context.Context, query string, args ...any) (Link, error) {
	row := p.conn.QueryRow(ctx, query, args...)
	link := Link{}
	err := row.Scan(&link.Name, &link.Description, &link.URL, &link.Views, &link.Created, &link.Updated, &link.CreatedBy, &link.Disabled)
	if errors.Is(err, pgx.ErrNoRows) {
		return link, ErrLinkNotFound
	}
	if err != nil {
		return link, err
	}
	return link, nil
}

func (p *postgres) getMultipleResults(ctx context.Context, query string, args ...any) ([]Link, error) {
	links := []Link{}
	rows, err := p.conn.Query(ctx, query, args...)
	if err != nil {
		return links, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var link Link
		err = rows.Scan(&link.Name, &link.Description, &link.URL, &link.Views, &link.Created, &link.Updated, &link.CreatedBy, &link.Disabled)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return links, fmt.Errorf("failed while scanning: %w", err)
		}
		links = append(links, link)
	}
	return links, nil
}

func (p *postgres) ensureTable(ctx context.Context) error {
	_, err := p.conn.Exec(ctx, `create table if not exists links (
		name text not null primary key,
		description text not null,
		url text not null,
		views integer default 0,
		created_at timestamptz not null,
		updated_at timestamptz not null,
		created_by text not null,
		disabled bool default false
	)`)
	if err != nil {
		return err
	}
	return nil
}
