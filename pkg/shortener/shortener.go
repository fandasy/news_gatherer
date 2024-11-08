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

// Если у вас нагрузка (кол-во пар) на кэш-мапу известна, то лучше будет буферизировать мапу

func NewUrlsMap() *UrlsMap {
	return &UrlsMap{
		m: make(map[string]string),
	}
}

func (m *UrlsMap) Get(key string) (string, bool) {
	m.rw.RLock()
	defer m.rw.RUnlock()
	val, ok := m.m[key]
	return val, ok
}

func (m *UrlsMap) Set(key, value string) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.m[key] = value
}

func (m *UrlsMap) Delete(keySlice []string) {
	m.rw.Lock()
	defer m.rw.Unlock()
	for _, key := range keySlice {
		delete(m.m, key)
	}
}

const (
	charset   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	keyLength = 10
)

func (m *UrlsMap) GenerateShortKey(originalUrls string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	m.rw.Lock()
	defer m.rw.Unlock()

	for {
		shortKey := make([]byte, keyLength)
		for i := range shortKey {
			shortKey[i] = charset[r.Intn(len(charset))]
		}

		if _, ok := m.m[string(shortKey)]; !ok {
			m.m[string(shortKey)] = originalUrls

			return string(shortKey)
		}
	}
}
