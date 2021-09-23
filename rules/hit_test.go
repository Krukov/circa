package rules

import (
	"circa/message"
	"circa/storages"
	"errors"

	"reflect"
	"testing"
	"time"
)

func TestHitRule_Process_Nocache(t *testing.T) {
	req := &message.Request{Method: "GET"}
	noCacheResp := message.NewResponse(200, []byte(`data`), map[string]string{})
	storage, _ := storages.StorageFormDSN("mem://")

	r := &HitRule{
		TTL:             time.Second,
		Hits:            3,
		UpdateAfterHits: 0,
	}
	got, _, _ := r.Process(req, "key", storage, func(request *message.Request) (*message.Response, error) {
		return noCacheResp, nil
	})

	if !reflect.DeepEqual(got, noCacheResp) {
		t.Errorf("Process() got = %v, want %v", got, noCacheResp)
	}

	cached, _ := storage.Get("key")
	if !reflect.DeepEqual(cached, noCacheResp) {
		t.Errorf("Process() got = %v, want %v", cached, noCacheResp)
	}
}

func TestHitRule_Process_Cache(t *testing.T) {
	req := &message.Request{Method: "GET"}
	noCacheResp := message.NewResponse(200, []byte(`data`), map[string]string{})
	cacheResp := message.NewResponse(200, []byte(`cache`), map[string]string{})
	storage, _ := storages.StorageFormDSN("mem://")

	storage.Set("key", cacheResp, time.Second)
	storage.Incr("key:hits") // 1

	r := &HitRule{
		TTL:             time.Second,
		Hits:            3,
		UpdateAfterHits: 0,
	}
	got, _, _ := r.Process(req, "key", storage, func(request *message.Request) (*message.Response, error) {
		return noCacheResp, nil
	})

	if !reflect.DeepEqual(got, cacheResp) {
		t.Errorf("Process() got = %v, want %v", got, cacheResp)
	}

	incr, _ := storage.Incr("key:hits")
	if incr != 3 {
		t.Errorf("Storage.incr got = %v, want 3", incr)
	}
}

func TestHitRule_Process_ProxyError(t *testing.T) {
	req := &message.Request{Method: "GET"}
	noCacheResp := &message.Response{Body: []byte(`data`), Status: 200}
	cacheResp := &message.Response{Body: []byte(`cache`), Status: 200}
	storage, _ := storages.StorageFormDSN("mem://")

	r := &HitRule{
		TTL:             time.Second,
		Hits:            3,
		UpdateAfterHits: 0,
	}
	got, _, err := r.Process(req, "key", storage, func(request *message.Request) (*message.Response, error) {
		return noCacheResp, errors.New("error")
	})

	if !reflect.DeepEqual(err, errors.New("error")) {
		t.Errorf("Process() got = %v, want %v", got, cacheResp)
	}

	stored, _ := storage.Get("key")
	if stored != nil {
		t.Errorf("do not store value at proxy errro")
	}
}

func TestHitRule_Process_CacheOverHit(t *testing.T) {
	req := &message.Request{Method: "GET"}
	noCacheResp := &message.Response{Body: []byte(`data`), Status: 200}
	cacheResp := &message.Response{Body: []byte(`cache`), Status: 200}
	storage, _ := storages.StorageFormDSN("mem://")

	storage.Set("key", cacheResp, time.Second)
	storage.Incr("key:hits") // 1
	storage.Incr("key:hits") // 2
	storage.Incr("key:hits") // 3

	r := &HitRule{
		TTL:  time.Second,
		Hits: 3,
	}
	got, _, _ := r.Process(req, "key", storage, func(request *message.Request) (*message.Response, error) {
		return noCacheResp, nil
	})

	if !reflect.DeepEqual(got, noCacheResp) {
		t.Errorf("Process() got = %v, want %v", got, noCacheResp)
	}

	incr, _ := storage.Incr("key:hits")
	if incr != 1 {
		t.Errorf("Storage.incr got = %v, want 1", incr)
	}
}

func TestHitRule_Process_CachePreHit(t *testing.T) {
	req := &message.Request{Method: "GET"}
	noCacheResp := &message.Response{Body: []byte(`data`), Status: 200}
	cacheResp := &message.Response{Body: []byte(`cache`), Status: 200}
	storage, _ := storages.StorageFormDSN("mem://")

	storage.Set("key", cacheResp, time.Second)
	storage.Incr("key:hits") // 1

	r := &HitRule{
		TTL:             time.Second,
		Hits:            3,
		UpdateAfterHits: 1,
	}
	got, _, _ := r.Process(req, "key", storage, func(request *message.Request) (*message.Response, error) {
		return noCacheResp, nil
	})

	if !reflect.DeepEqual(got, cacheResp) {
		t.Errorf("Process() got = %v, want %v", got, cacheResp)
	}

	time.Sleep(time.Millisecond)
	incr, _ := storage.Incr("key:hits")
	if incr != 1 {
		t.Errorf("Storage.incr got = %v, want 1", incr)
	}
}
