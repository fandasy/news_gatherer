package storage

import (
	"context"
	"errors"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickPageList(ctx context.Context, userName string) (*PageList, int, error)
	Remove(ctx context.Context, p *Page) error
	GetAllNews(ctx context.Context, username string) (*NewsList, error)
	PickNews(ctx context.Context, page *Page) (*NewsList, error)
	IsExists(ctx context.Context, p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL       string
	UserName  string
	Assembler string
}

type PageList struct {
	URLS     []string
	UserName string
}

type NewsList struct {
	News     []News
	UserName string
}

type News struct {
	URL       string
	Assembler string
}
