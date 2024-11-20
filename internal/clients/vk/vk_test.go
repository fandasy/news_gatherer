package vk

import (
	"context"
	"net/http"
	"testing"
)

const (
	token      = ""
	host       = "api.vk.com"
	apiVersion = "5.131"

	validGroupID = ""
)

var (
	vkClient = &Client{
		host:    host,
		version: apiVersion,
		token:   token,
		client:  http.Client{},
	}
)

func TestClient_GetNews(t *testing.T) {
	type args struct {
		ctx     context.Context
		groupID string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantErr bool
	}{
		{
			name: "validate",
			c:    vkClient,
			args: args{
				ctx:     context.Background(),
				groupID: validGroupID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.c
			gotNews, err := c.GetNews(tt.args.ctx, tt.args.groupID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNews() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				for _, v := range gotNews {
					t.Log(v)
				}
			}
		})
	}
}

func TestClient_ValidateNewsGroup(t *testing.T) {
	type args struct {
		ctx     context.Context
		groupID string
	}
	tests := []struct {
		name    string
		c       *Client
		args    args
		wantVal bool
		wantErr bool
	}{
		{
			name: "validate",
			c:    vkClient,
			args: args{
				ctx:     context.Background(),
				groupID: validGroupID,
			},
			wantVal: true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.c
			gotVal, err := c.ValidateNewsGroup(tt.args.ctx, tt.args.groupID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNewsGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotVal != tt.wantVal {
				t.Errorf("ValidateNewsGroup() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
		})
	}
}
