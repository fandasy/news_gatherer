package rss

import (
	"context"
	"testing"
)

const (
	validUrl = "https://ria.ru/export/rss2/archive/index.xml"
)

func TestParsing(t *testing.T) {
	type args struct {
		ctx     context.Context
		feedURL string
	}
	tests := []struct {
		name    string
		args    args
		valid   bool
		wantErr bool
	}{
		{
			name: "valid URL",
			args: args{
				ctx:     context.Background(),
				feedURL: validUrl,
			},
			valid:   true,
			wantErr: false,
		},
		{
			name: "invalid URL",
			args: args{
				ctx:     context.Background(),
				feedURL: "https://ria.ru/example/index.xml",
			},
			valid:   false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if ok := ValidateFeedURL(tt.args.ctx, tt.args.feedURL); ok != tt.valid {
				t.Errorf("ValidateFeedURL() = %v, want %v", ok, tt.valid)
			}

			if tt.valid {
				news, err := Parsing(tt.args.ctx, tt.args.feedURL)

				if (err != nil) != tt.wantErr {
					t.Errorf("Parsing() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !tt.wantErr {
					for _, v := range news {
						t.Log(v)
					}
				}
			}
		})
	}
}
