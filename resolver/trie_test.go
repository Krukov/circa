package resolver

import (
	"reflect"
	"testing"
)

func createRouter() *node {
	t := newTrie("")
	t.addRule("/")
	t.addRule("user/me")
	t.addRule("user/{id}")
	t.addRule("user/{id}/top")
	t.addRule("user/{id}/top/{item}")
	t.addRule("/user/")

	t.addRule("posts/last/{id}")
	t.addRule("posts/*")
	t.addRule("posts/{id}/*")
	t.addRule("posts/{id}/*")
	t.addRule("posts/{id}")
	t.addRule("posts/{id}/info")
	t.addRule("posts/info/{id}")
	return t
}

func Test_node_getRoute(t1 *testing.T) {
	t := createRouter()

	tests := []struct {
		name       string
		gotPath    string
		wantNames  []string
		wantParams map[string]string
	}{
		{"user list", "/user", []string{"/user/"}, map[string]string{}},
		{"user list", "/user/", []string{"/user/"}, map[string]string{}},
		{"user list", "user/", []string{"/user/"}, map[string]string{}},

		{"user detail", "/user/1", []string{"user/{id}"}, map[string]string{"id": "1"}},
		{"user detail", "/user/test/", []string{"user/{id}"}, map[string]string{"id": "test"}},

		{"user me", "/user/me/", []string{"user/me"}, map[string]string{}},
		{"user me", "/user/me", []string{"user/me"}, map[string]string{}},

		{"user top", "/user/1/top", []string{"user/{id}/top"}, map[string]string{"id": "1"}},
		{"user top", "/user/me/top", []string{"user/{id}/top"}, map[string]string{"id": "me"}},
		{"user top list", "/user/1/top/1", []string{"user/{id}/top/{item}"}, map[string]string{"id": "1", "item": "1"}},

		{"posts info", "/posts/1/info", []string{"posts/*", "posts/{id}/*", "posts/{id}/info"}, map[string]string{"id": "1"}},
		{"posts info", "/posts/info/1", []string{"posts/*", "posts/info/{id}"}, map[string]string{"id": "1"}},
		{"posts last", "/posts/last/1", []string{"posts/*", "posts/last/{id}"}, map[string]string{"id": "1"}},
		{"posts detail", "/posts/1", []string{"posts/*", "posts/{id}"}, map[string]string{"id": "1"}},

		{"posts pass", "/posts/1/inf/some/here", []string{"posts/*", "posts/{id}/*"}, map[string]string{}},
		{"posts pass 2", "/posts/1/inf", []string{"posts/*", "posts/{id}/*"}, map[string]string{}},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {

			gotNames, gotParams, err := t.resolve(tt.gotPath)
			if err != nil {
				t1.Errorf("resolve() error = %v", err)
				return
			}
			if len(tt.wantNames) > 0 && !reflect.DeepEqual(gotNames, tt.wantNames) {
				t1.Errorf("resolve() gotNames = %v, want %v", gotNames, tt.wantNames)
			}
			if len(tt.wantParams) > 0 && !reflect.DeepEqual(gotParams, tt.wantParams) {
				t1.Errorf("resolve() gotParams = %v, want %v", gotParams, tt.wantParams)
			}
		})
	}
}

func Test_node_getRoute_Error(t1 *testing.T) {
	t := createRouter()

	tests := []struct {
		name       string
		gotPath    string
		wantName   string
		wantParams map[string]string
	}{

		{"not find", "/user/top/me", "", map[string]string{}},
		{"not find", "/posts/", "", map[string]string{}},
		{"not find", "/use/me", "", map[string]string{}},
		{"not find", "/user/1/1", "", map[string]string{}},
		{"not find", "/user/1/2/top/", "", map[string]string{}},
		{"not find", "?", "", map[string]string{}},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			names, _, err := t.resolve(tt.gotPath)
			if err == nil {
				t1.Errorf("resolve() error = %v names=%v, wantErr", err, names)
			}
		})
	}
}

func BenchmarkRoot(b *testing.B) {

	t := createRouter()
	for n := 0; n < b.N; n++ {
		t.resolve("/")
	}
}

func BenchmarkNotFound(b *testing.B) {
	t := createRouter()
	for n := 0; n < b.N; n++ {
		t.resolve("user/me/this/path/do/not/exists")
	}
}

func BenchmarkLastNode(b *testing.B) {
	b.ReportAllocs()
	t := createRouter()
	for n := 0; n < b.N; n++ {
		t.resolve("/user/1/top/1")
	}
}

func BenchmarkMultiRuleNode(b *testing.B) {
	b.ReportAllocs()
	t := createRouter()
	for n := 0; n < b.N; n++ {
		t.resolve("/posts/1/info")
	}
}
