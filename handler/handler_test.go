package handler

import (
	"circa/message"
	"circa/rules"
	"circa/storages"
	"reflect"
	"testing"
)

func fakeCall (r *message.Request) (*message.Response, error) {
	return &message.Response{Status: 100}, nil
}


func TestRunner_Handle(t *testing.T) {

	r := NewRunner(fakeCall)
	h := NewHandler(&rules.ProxyRule{}, &storages.Memory{}, "key", &message.Request{}, []string{"get"})
	r.AddHandler("/", h)

	tests := []struct {
		name     string
		request  *message.Request
		wantResp *message.Response
		wantErr  bool
	}{
		{"simple", &message.Request{Path: "/"}, &message.Response{Status: 100}, false},
		{"404", &message.Request{Path: "/404"}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := r.Handle(tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("Handle() got = %v, want %v", got, tt.wantResp)
			}
		})
	}
}
