package telegram

import (
	"context"
	"net/http"
	"testing"
)

const (
	token = ""
	host  = "api.telegram.org"
)

var (
	tgClient = &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
)

func TestClient_Updates(t *testing.T) {
	type args struct {
		ctx    context.Context
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "validate",
			c:    tgClient,
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  100,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.c

			gotUpdates, err := c.Updates(tt.args.ctx, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Updates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for _, gotUpdate := range gotUpdates {
					if gotUpdate.CallbackQuery != nil {

						t.Log(gotUpdate.CallbackQuery)
					} else {

						t.Log(gotUpdate.Message)
					}
				}
			}
		})
	}
}

func TestClient_AnswerCallbackQuery(t *testing.T) {
	type args struct {
		ctx        context.Context
		callbackID string
		text       string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "validate",
			c:    tgClient,
			args: args{
				ctx:        context.Background(),
				callbackID: "",
				text:       "text",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.c

			if err := c.AnswerCallbackQuery(tt.args.ctx, tt.args.callbackID, tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("AnswerCallbackQuery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SendMessageText(t *testing.T) {
	type args struct {
		ctx    context.Context
		chatID int
		text   string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "validate",
			c:    tgClient,
			args: args{
				ctx:    context.Background(),
				chatID: 0,
				text:   "text",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.c

			if err := c.SendMessageText(tt.args.ctx, tt.args.chatID, tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("SendMessageText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_SendMessageTextAndButton(t *testing.T) {
	type args struct {
		ctx    context.Context
		chatID int
		text   string
		button InlineKeyboardMarkup
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "validate",
			c:    tgClient,
			args: args{
				ctx:    context.Background(),
				chatID: 0,
				text:   "text",
				button: InlineKeyboardMarkup{
					InlineKeyboard: [][]InlineKeyboardButton{
						{
							{
								Text:         "Button_1",
								CallbackData: "/Response",
							},
							{
								Text:         "Button_2",
								CallbackData: "/Response",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.c

			if err := c.SendMessageTextAndButton(tt.args.ctx, tt.args.chatID, tt.args.text, tt.args.button); (err != nil) != tt.wantErr {
				t.Errorf("SendMessageTextAndButton() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
