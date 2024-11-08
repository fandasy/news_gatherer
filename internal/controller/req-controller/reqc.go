package req_controller

import (
	"sync"
	"time"
)

type ReqCounter struct {
	m  map[string]*UserControl
	rw sync.RWMutex
}

type UserControl struct {
	msgCounter  uint
	lastMsgTime time.Time
	bannedUntil time.Time
}

type LimitOptions struct {
	MaxNumReq uint
	TimeSlice time.Duration
	BanTime   time.Duration
}

func NewLimitOptions(muxNumReq uint, timeSlice time.Duration, banTime time.Duration) LimitOptions {
	return LimitOptions{
		MaxNumReq: muxNumReq,
		TimeSlice: timeSlice,
		BanTime:   banTime,
	}
}

func New() *ReqCounter {
	return &ReqCounter{
		m: make(map[string]*UserControl),
	}
}

func (r *ReqCounter) Checking(username string, options LimitOptions) bool {

	user, ok := r.Get(username)
	if !ok {
		r.Set(username)

	} else {

		if !user.bannedUntil.IsZero() && user.bannedUntil.After(time.Now()) {
			return false
		}

		if time.Since(user.lastMsgTime) < options.TimeSlice*time.Second {
			if user.msgCounter >= options.MaxNumReq {
				r.Ban(username, time.Now().Add(options.BanTime*time.Second))

				return false
			}
		} else {

			r.Reset(username)
		}
	}

	return true
}

func (r *ReqCounter) Add(key string, number uint) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.m[key].msgCounter += number
}

func (r *ReqCounter) Get(key string) (*UserControl, bool) {
	r.rw.RLock()
	defer r.rw.RUnlock()
	val, ok := r.m[key]
	return val, ok
}

func (r *ReqCounter) Set(key string) {
	r.rw.Lock()
	defer r.rw.Unlock()

	user := &UserControl{
		msgCounter:  0,
		lastMsgTime: time.Now(),
	}

	r.m[key] = user
}

func (r *ReqCounter) Reset(key string) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.m[key].msgCounter = 0
	r.m[key].lastMsgTime = time.Now()
}

func (r *ReqCounter) Ban(key string, bannedUntil time.Time) {
	r.rw.Lock()
	defer r.rw.Unlock()
	r.m[key].msgCounter = 0
	r.m[key].bannedUntil = bannedUntil
}
