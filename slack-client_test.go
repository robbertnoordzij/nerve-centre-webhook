package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendSlack(t *testing.T) {
	type args struct {
		payload *SlackPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Empty",
			args: args{
				payload: nil,
			},
			wantErr: true,
		},
		{
			name: "Payload",
			args: args{
				payload: &SlackPayload{
					Text: "Hello",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				body, _ := ioutil.ReadAll(r.Body)

				if len(body) == 0 || string(body) == "null" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				w.WriteHeader(http.StatusOK)
			}))
			defer ts.Close()

			if err := SendSlack(ts.URL, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("SendSlack() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
