package rules

//
//import (
//	"circa/message"
//	"errors"
//	"reflect"
//	"testing"
//	"time"
//
//	"github.com/stretchr/testify/mock"
//)
//
//type mockStorage struct {
//	mock.Mock
//}
//
//func (s *mockStorage) String() string {
//	return "mock"
//}
//
//func (s *mockStorage) Set(key string, value *message.Response, ttl time.Duration) (bool, error) {
//	args := s.Called(key, value, ttl)
//	return args.Bool(0), args.Error(1)
//}
//
//func (s *mockStorage) Del(key string) (bool, error) {
//	args := s.Called(key)
//	return args.Bool(0), args.Error(1)
//}
//
//func (s *mockStorage) Get(key string) (*message.Response, error) {
//	args := s.Called(key)
//	return &message.Response{Status: 200, Body: []byte(`mock`)}, args.Error(0)
//}
//
//func TestCacheRule_Process(t *testing.T) {
//	tests := []struct {
//		name    string
//		req     *message.Request
//		cache   string
//		resp    *message.Response
//		wantErr bool
//	}{
//		{"simple", &message.Request{Method: "GET"}, "cached", &message.Response{Body: []byte(`mock`), Status: 200, CachedKey: "key"}, false},
//		{"no cache", &message.Request{Method: "GET"}, "", &message.Response{Body: []byte(`data`), Status: 200}, false},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			storage := &mockStorage{}
//			var cacheError error
//			if tt.cache == "" {
//				cacheError = errors.New("empty cache")
//				storage.On("Set", "key", mock.Anything, time.Second).Return(true, nil)
//			}
//			storage.On("Get", "key").Return(cacheError)
//			r := &CacheRule{
//				TTL: time.Second,
//			}
//			got, _, err := r.Process(tt.req, "key", storage, func(request *message.Request) (*message.Response, error) {
//				return &message.Response{Body: []byte(`data`), Status: 200}, nil
//			})
//			storage.AssertExpectations(t)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.resp) {
//				t.Errorf("Process() got = %v, want %v", got, tt.resp)
//			}
//		})
//	}
//}
