package handler

import (
	"reflect"
	"testing"
)

func createRouter() *node {
	t := newTrie()
	t.addRule("/", "ROOT", false)
	t.addRule("user/me", "USER_ME",false)
	t.addRule("user/{id}", "USER_DETAIL", false)
	t.addRule("posts/{id}/info", "POSTS_DETAIL", false)
	t.addRule("posts/{id}/last", "POSTS_LAST", true)
	t.addRule("posts/*", "POSTS_PASS", false)
	t.addRule("user/{id}/top", "USER_TOP_LIST", false)
	t.addRule("user/{id}/top/{item}", "USER_TOP_ITEM", false)
	t.addRule("/user/", "USER_LIST", false)
	return t
}

func Test_node_getRoute(t1 *testing.T) {
	t := createRouter()

	tests := []struct {
		name       string
		gotPath    string
		wantName   string
		wantParams map[string]string
		wantErr    bool
	}{
		{"user list", "/user", "USER_LIST", map[string]string{}, false},
		{"user list", "/user/", "USER_LIST", map[string]string{}, false},
		{"user list", "user/", "USER_LIST", map[string]string{}, false},

		{"user detail", "/user/1", "USER_DETAIL", map[string]string{"id": "1"}, false},
		{"user detail", "/user/test/", "USER_DETAIL", map[string]string{"id": "test"}, false},

		{"user me", "/user/me/", "USER_ME", map[string]string{}, false},
		{"user me", "/user/me", "USER_ME", map[string]string{}, false},

		{"user top", "/user/1/top", "USER_TOP_LIST", map[string]string{"id": "1"}, false},
		{"user top", "/user/me/top", "USER_TOP_LIST", map[string]string{"id": "me"}, false},
		{"user top list", "/user/1/top/1", "USER_TOP_ITEM", map[string]string{"id": "1", "item": "1"}, false},

		{"posts info", "/posts/1/info", "POSTS_DETAIL", map[string]string{"id": "1"}, false},
		{"posts info", "/posts/1/inf", "POSTS_PASS", map[string]string{}, false},

		{"not find", "/user/top/me", "", map[string]string{}, true},
		{"not find", "/use/me", "", map[string]string{}, true},
		{"not find", "/user/1/1", "", map[string]string{}, true},
		{"not find", "/user/1/2/top/", "", map[string]string{}, true},
		{"not find", "/posts", "", map[string]string{}, true},
		{"not find", "?", "", map[string]string{}, true},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			gotName, gotParams, err := t.getRoute(tt.gotPath)
			if (err != nil) != tt.wantErr {
				t1.Errorf("getRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(gotName) != tt.wantName {
				t1.Errorf("getRoute() gotName = %v, want %v", gotName, tt.wantName)
			}
			if len(tt.wantParams) > 0 && !reflect.DeepEqual(gotParams, tt.wantParams) {
				t1.Errorf("getRoute() gotParams = %v, want %v", gotParams, tt.wantParams)
			}
		})
	}
}

func Test_node_remove(t1 *testing.T) {
	t := createRouter()

	t.removeRule("USER_TOP_ITEM")

	path, _, err := t.getRoute("user/1/top/test")
	if err == nil {
		t1.Errorf("the route steal exists %v", path)
	}
}

func BenchmarkRoot(b *testing.B) {

	t := createRouter()
	for n := 0; n < b.N; n++ {
		t.getRoute("/")
	}
}

func BenchmarkNotFound(b *testing.B) {
	t := createRouter()
	for n := 0; n < b.N; n++ {
		t.getRoute("user/me/this/path/do/not/exists")
	}
}

func BenchmarkLastNode(b *testing.B) {
	b.ReportAllocs()
	t := createRouter()
	for n := 0; n < b.N; n++ {
		t.getRoute("/user/1/top/1")
	}
}