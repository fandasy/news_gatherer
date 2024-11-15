package req_controller

import (
	"sync"
	"sync/atomic"
	"time"

	"telegramBot/internal/config/j"
)

type ReqCounter struct {
	m sync.Map
}

type UserControl struct {
	msgCounter  int32
	lastMsgTime int64 // хранить время в наносекундах
	bannedUntil int64 // хранить время в наносекундах
}

type LimitOptions struct {
	MaxNumReq uint
	TimeSlice time.Duration
	BanTime   time.Duration
}

func NewLimitOptions(limit j.ReqLimit) LimitOptions {
	return LimitOptions{
		MaxNumReq: limit.MaxNumberReq,
		TimeSlice: limit.TimeSlice,
		BanTime:   limit.BanTime,
	}
}

func New() *ReqCounter {
	return &ReqCounter{}
}

func (r *ReqCounter) Checking(username string, options LimitOptions) bool {
	user, ok := r.GetOrSet(username)
	if !ok {
		return true
	}

	bannedUntil := atomic.LoadInt64(&user.bannedUntil)
	if bannedUntil > time.Now().UnixNano() {
		return false
	}

	lastMsgTime := atomic.LoadInt64(&user.lastMsgTime)
	if time.Since(time.Unix(0, lastMsgTime)) < options.TimeSlice {

		if atomic.LoadInt32(&user.msgCounter) >= int32(options.MaxNumReq) {

			user.Ban(time.Now().Add(options.BanTime))

			return false
		}

		user.Add(1)

	} else {
		user.Reset()
	}

	return true
}

func (r *ReqCounter) GetOrSet(key string) (*UserControl, bool) {
	user, loaded := r.m.LoadOrStore(key, &UserControl{
		msgCounter:  1,
		lastMsgTime: time.Now().UnixNano(),
	})
	if !loaded {
		return user.(*UserControl), false
	}

	return user.(*UserControl), true
}

func (u *UserControl) Add(number uint) {
	atomic.AddInt32(&u.msgCounter, int32(number))
}

func (u *UserControl) Reset() {
	atomic.StoreInt32(&u.msgCounter, 1)
	atomic.StoreInt64(&u.lastMsgTime, time.Now().UnixNano())
}

func (u *UserControl) Ban(bannedUntil time.Time) {
	atomic.StoreInt32(&u.msgCounter, 0)
	atomic.StoreInt64(&u.bannedUntil, bannedUntil.UnixNano())
}
