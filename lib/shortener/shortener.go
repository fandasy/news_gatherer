package shortener

import (
	"math/rand"
	"sync"
	"time"
)

// Я не использую sync.Map потому что в мапу чаще всего будут записывать чем читать

type UrlsMap struct {
	m  map[string]string
	rw sync.RWMutex
}

// Если у вас нагрузка (кол-во пар) на кэш-мапу известна, то лучше будет буферизировать мапу, ну или использовать redis для этой задачи

func NewUrlsMap() *UrlsMap {
	return &UrlsMap{
		m: make(map[string]string),
	}
}

func (c *UrlsMap) Get(key string) (string, bool) {
	c.rw.RLock()
	defer c.rw.RUnlock()
	val, ok := c.m[key]
	return val, ok
}

func (c *UrlsMap) Set(key, value string) {
	c.rw.Lock()
	defer c.rw.Unlock()
	c.m[key] = value
}

const (
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength = 10
)

func (c *UrlsMap) GenerateShortKey(originalUrls string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	c.rw.Lock()
	defer c.rw.Unlock()

	for {
		shortKey := make([]byte, keyLength)
		for i := range shortKey {
			shortKey[i] = charset[r.Intn(len(charset))]
		}

		if _, ok := c.m[string(shortKey)]; !ok {
			c.m[string(shortKey)] = originalUrls

			return string(shortKey)
		}
	}
}
