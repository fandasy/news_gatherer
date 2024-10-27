package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"telegramBot/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("postgres", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, p *storage.Page) error {

	q := `INSERT INTO pages (url, user_name, assembler) VALUES($1, $2, $3)`

	if _, err := s.db.ExecContext(ctx, q, p.URL, p.UserName, p.Assembler); err != nil {
		return fmt.Errorf("can't save page: %w", err)
	}

	return nil
}

func (s *Storage) Remove(ctx context.Context, page *storage.Page) error {

	q := `DELETE FROM pages WHERE url = $1 AND user_name = $2`

	if _, err := s.db.ExecContext(ctx, q, page.URL, page.UserName); err != nil {
		return fmt.Errorf("can't pick remove page: %w", err)
	}

	return nil
}

func (s *Storage) PickPageList(ctx context.Context, username string) (*storage.PageList, int, error) {

	qCount := `SELECT COUNT(*) FROM pages WHERE user_name = $1`

	var count int

	if err := s.db.QueryRowContext(ctx, qCount, username).Scan(&count); err != nil {
		return nil, 0, fmt.Errorf("unable to select page existence check: %w", err)
	}

	urls := make([]string, 0, count)

	q := `SELECT url FROM pages WHERE user_name = $1 ORDER BY assembler DESC`

	rows, err := s.db.QueryContext(ctx, q, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, 0, fmt.Errorf("can't pick page list: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, 0, fmt.Errorf("can't pick rows page list: %w", err)
		}
		urls = append(urls, url)
	}

	return &storage.PageList{
		URLS:     urls,
		UserName: username,
	}, count, nil
}

func (s *Storage) GetAllNews(ctx context.Context, username string) (*storage.NewsList, error) {

	qCount := `SELECT COUNT(*) FROM pages WHERE user_name = $1`

	var count int

	if err := s.db.QueryRowContext(ctx, qCount, username).Scan(&count); err != nil {
		return nil, fmt.Errorf("unable to select page existence check: %w", err)
	}

	if count == 0 {
		return nil, storage.ErrNoSavedPages
	}

	newsArr := make([]storage.News, 0, count)

	q := `SELECT url, assembler FROM pages WHERE user_name = $1 ORDER BY assembler ASC`

	rows, err := s.db.QueryContext(ctx, q, username)
	if err != nil {
		return nil, fmt.Errorf("can't get news list: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var (
			url       string
			assembler string
		)

		if err := rows.Scan(&url, &assembler); err != nil {
			return nil, fmt.Errorf("can't get rows news list: %w", err)
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

	qCount := `SELECT COUNT(*) FROM pages WHERE user_name = $1`

	var count int

	if err := s.db.QueryRowContext(ctx, qCount, page.UserName).Scan(&count); err != nil {
		return nil, fmt.Errorf("unable to select page existence check: %w", err)
	}

	if count == 0 {
		return nil, storage.ErrNoSavedPages
	}

	newsArr := make([]storage.News, 0, count)

	q := `SELECT url FROM pages WHERE assembler = $1 AND user_name = $2`

	rows, err := s.db.QueryContext(ctx, q, page.Assembler, page.UserName)
	if err != nil {
		return nil, fmt.Errorf("can't pick news list: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var url string

		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("can't pick rows news list: %w", err)
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
	var count int

	if page.URL == "" {
		q := `SELECT COUNT(*) FROM pages WHERE assembler = $1 AND user_name = $2`

		if err := s.db.QueryRowContext(ctx, q, page.Assembler, page.UserName).Scan(&count); err != nil {
			return false, fmt.Errorf("can't pick check if assembler exists: %w", err)
		}
	} else {
		q := `SELECT COUNT(*) FROM pages WHERE url = $1 AND user_name = $2`

		if err := s.db.QueryRowContext(ctx, q, page.URL, page.UserName).Scan(&count); err != nil {
			return false, fmt.Errorf("can't pick check if page exists: %w", err)
		}
	}

	return count > 0, nil
}

func (s *Storage) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (user_name TEXT, url TEXT, assembler TEXT)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}

	return nil
}
