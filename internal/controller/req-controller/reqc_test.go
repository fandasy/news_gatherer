package req_controller

import (
	"sync"
	"testing"
	"time"
)

func testLimitOptions(m uint32, t time.Duration, b time.Duration) RateLimit {
	return RateLimit{
		limit:    m,
		interval: t,
		banTime:  b,
	}
}

func TestReqCounter_Checking_SingleUser(t *testing.T) {
	type args struct {
		reqNum    int
		sleepTime time.Duration
		username  string
		options   RateLimit
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
			r := &ReqCounter{
				rl: tt.args.options,
			}

			var wg sync.WaitGroup

			result := make(chan bool, 1)

			for i := 0; i < tt.args.reqNum; i++ {
				time.Sleep(tt.args.sleepTime)

				wg.Add(1)
				go func() {
					defer wg.Done()

					ok := r.Checking(tt.args.username)

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
		options RateLimit
	}
	tests := []struct {
		name              string
		args              args
		want              bool
		number_of_blocked int
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
			want:              false,
			number_of_blocked: 1,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			r := &ReqCounter{
				rl: tt.args.options,
			}

			var (
				wg  sync.WaitGroup
				buf int
			)

			for _, user := range tt.args.users {
				buf += user.reqNum
			}

			result := make(chan bool, buf)

			for _, user := range tt.args.users {
				user := user

				for i := 0; i < user.reqNum; i++ {
					time.Sleep(user.sleepTime)

					wg.Add(1)
					go func() {
						defer wg.Done()

						ok := r.Checking(user.username)

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

			blocked := make([]bool, 0, buf)

			for i := 0; i < buf; i++ {
				select {
				case <-result:
					blocked = append(blocked, false)
				default:
				}
			}

			var ok bool

			if len(blocked) == 0 {
				ok = true
			}

			if ok != tt.want {
				t.Errorf("Checking() = %v, want %v", ok, tt.want)

				return
			}

			if ok == false {
				if len(blocked) != tt.number_of_blocked {
					t.Errorf("Checking() = %v, want %v", len(blocked), tt.number_of_blocked)
				}
			}
		})
	}
}
