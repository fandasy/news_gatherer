package req_controller

import (
	"sync"
	"testing"
	"time"
)

func testLimitOptions(m uint, t time.Duration, b time.Duration) LimitOptions {
	return LimitOptions{
		MaxNumReq: m,
		TimeSlice: t,
		BanTime:   b,
	}
}

func TestReqCounter_Checking_SingleUser(t *testing.T) {
	type args struct {
		reqNum    int
		sleepTime time.Duration
		username  string
		options   LimitOptions
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid requests",
			args: args{
				reqNum:    5,
				sleepTime: 500 * time.Millisecond,
				username:  "Maks",
				options:   testLimitOptions(4, 2*time.Second, 60*time.Second),
			},
			want: true,
		},
		{
			name: "invalid requests",
			args: args{
				reqNum:   2,
				username: "Igor",
				options:  testLimitOptions(1, 2*time.Second, 60*time.Second),
			},
			want: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			r := &ReqCounter{}

			var wg sync.WaitGroup

			result := make(chan bool, 1)

			for i := 0; i < tt.args.reqNum; i++ {
				time.Sleep(tt.args.sleepTime)

				wg.Add(1)
				go func() {
					defer wg.Done()

					ok := r.Checking(tt.args.username, tt.args.options)

					if !ok {
						select {
						case result <- false:
						default:
						}
					}

				}()
			}

			wg.Wait()

			ok := true

			select {
			case ok = <-result:
			default:
			}

			if ok != tt.want {
				t.Errorf("Checking() = %v, want %v", ok, tt.want)
			}
		})
	}
}

func TestReqCounter_Checking_ManyUsers(t *testing.T) {
	type user struct {
		reqNum    int
		username  string
		sleepTime time.Duration
	}
	type args struct {
		users   []user
		options LimitOptions
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "valid requests",
			args: args{
				users: []user{
					{
						reqNum:    4,
						sleepTime: 1 * time.Millisecond,
						username:  "Maks",
					},
					{
						reqNum:    5,
						sleepTime: 500 * time.Millisecond,
						username:  "Vlad",
					},
					{
						reqNum:   2,
						username: "Ivan",
					},
				},
				options: testLimitOptions(4, 2*time.Second, 60*time.Second),
			},
			want: true,
		},
		{
			name: "invalid requests",
			args: args{
				users: []user{
					{
						reqNum:   2,
						username: "Igor",
					},
					{
						reqNum:   1,
						username: "Sonya",
					},
				},
				options: testLimitOptions(1, 2*time.Second, 60*time.Second),
			},
			want: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			r := &ReqCounter{}

			var wg sync.WaitGroup

			result := make(chan bool, 1)
			for _, user := range tt.args.users {
				user := user

				for i := 0; i < user.reqNum; i++ {
					time.Sleep(user.sleepTime)

					wg.Add(1)
					go func() {
						defer wg.Done()

						ok := r.Checking(user.username, tt.args.options)

						if !ok {
							select {
							case result <- false:
							default:
							}
						}

					}()
				}
			}

			wg.Wait()

			ok := true

			select {
			case ok = <-result:
			default:
			}

			if ok != tt.want {
				t.Errorf("Checking() = %v, want %v", ok, tt.want)
			}
		})
	}
}
