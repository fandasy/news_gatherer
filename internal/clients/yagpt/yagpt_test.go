package yagpt

import (
	"context"
	"net/http"
	"testing"
)

const (
	host  = "300.ya.ru"
	token = ""

	validUrl = ""
)

var (
	yagptClient = &Client{
		host:   host,
		token:  newToken(token),
		client: &http.Client{},
	}
)

func TestClient_GetRetelling(t *testing.T) {
	type args struct {
		ctx     context.Context
		pageURL string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "validate",
			c:    yagptClient,
			args: args{
				ctx:     context.Background(),
				pageURL: validUrl,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.c
			gotRetelling, err := c.GetRetelling(tt.args.ctx, tt.args.pageURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRetelling() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				t.Log(gotRetelling)
			}
		})
	}
}
