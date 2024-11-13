package psql

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log/slog"
	"telegramBot/internal/storage"
	"telegramBot/pkg/e"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

func New(path string, log *slog.Logger) (*Storage, error) {
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, e.Wrap("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.Wrap("can't ping database", err)
	}

	return &Storage{
		db:  db,
		log: log,
	}, nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {
	const op = "storage.psql.Save"

	q := `INSERT INTO pages (url, user_name, assembler) VALUES($1, $2, $3)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName, p.Assembler); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {
	const op = "storage.psql.Remove"

	q := `DELETE FROM pages WHERE url = $1 AND user_name = $2`

	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return e.Wrap(op, err)
	}

	return nil
}

func (s *Storage) PickPageList(ctx context.Context, username string) (*storage.PageList, int, error) {
	const op = "storage.psql.PickPageList"

	q := `SELECT url FROM pages WHERE user_name = $1 ORDER BY assembler ASC`

	rows, err := s.db.QueryContext(ctx, q, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, storage.ErrNoSavedPages
		}

		return nil, 0, e.Wrap(op, err)
	}

	defer rows.Close()

	var (
		urls  []string
		count int
	)

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, 0, e.Wrap(op, err)
		}

		count++
		urls = append(urls, url)
	}

	return &storage.PageList{
		URLS:     urls,
		UserName: username,
	}, count, nil
}

func (s *Storage) GetAllNews(ctx context.Context, username string) (*storage.NewsList, error) {
	const op = "storage.psql.GetAllNews"

	q := `SELECT url, assembler FROM pages WHERE user_name = $1 ORDER BY assembler ASC`

	rows, err := s.db.QueryContext(ctx, q, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoSavedPages
		}

		return nil, e.Wrap(op, err)
	}

	var newsArr []storage.News

	defer rows.Close()

	for rows.Next() {
		var (
			url       string
			assembler string
		)

		if err := rows.Scan(&url, &assembler); err != nil {
			return nil, e.Wrap(op, err)
		}

		news := storage.News{
			URL:       url,
			Assembler: assembler,
		}
		newsArr = append(newsArr, news)
	}

	return &storage.NewsList{
		News:     newsArr,
		UserName: username,
	}, nil
}

func (s *Storage) PickNews(ctx context.Context, page *storage.Page) (*storage.NewsList, error) {
	const op = "storage.psql.PickNews"

	q := `SELECT url FROM pages WHERE assembler = $1 AND user_name = $2`

	rows, err := s.db.QueryContext(ctx, q, page.Assembler, page.UserName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrNoSavedPages
		}

		return nil, e.Wrap(op, err)
	}

	defer rows.Close()

	var newsArr []storage.News

	for rows.Next() {
		var url string

		if err := rows.Scan(&url); err != nil {
			return nil, e.Wrap(op, err)
		}

		news := storage.News{
			URL:       url,
			Assembler: page.Assembler,
		}
		newsArr = append(newsArr, news)
	}

	return &storage.NewsList{
		News:     newsArr,
		UserName: page.UserName,
	}, nil
}

func (s *Storage) IsExists(ctx context.Context, page *storage.Page) (bool, error) {
	const op = "storage.psql.IsExists"

	var count int

	if page.URL == "" {
		q := `SELECT COUNT(*) FROM pages WHERE assembler = $1 AND user_name = $2`

		if err := s.db.QueryRowContext(ctx, q, page.Assembler, page.UserName).Scan(&count); err != nil {
			return false, e.Wrap(op, err)
		}
	} else {
		q := `SELECT COUNT(*) FROM pages WHERE url = $1 AND user_name = $2`

		if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
			return false, e.Wrap(op, err)
		}
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	const op = "storage.psql.Init"

	q := `CREATE TABLE IF NOT EXISTS pages (
    user_name TEXT,
    url TEXT,
    assembler TEXT);
	CREATE INDEX IF NOT EXISTS idx_user_name ON pages (user_name);`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return e.Wrap(op, err)
	}

	return nil
}
