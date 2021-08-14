package handler

import (
	"reflect"
	"testing"
)

func createRouter() *node {
	t := newTrie()
	t.addRule("/", "ROOT")
	t.addRule("user/me", "USER_ME")
	t.addRule("user/{id}", "USER_DETAIL")
	t.addRule("user/{id}/top", "USER_TOP_LIST")
	t.addRule("user/{id}/top/{item}", "USER_TOP_ITEM")
	t.addRule("/user/", "USER_LIST")

	t.addRule("posts/last/{id}", "POSTS_LAST")
	t.addRule("posts/*", "POSTS_PASS")
	t.addRule("posts/{id}/*", "POSTS_PASS2")
	t.addRule("posts/{id}/*", "POSTS_PASS2")
	t.addRule("posts/{id}", "POSTS_DETAIL")
	t.addRule("posts/{id}/info", "POSTS_DETAIL_INFO")
	t.addRule("posts/info/{id}", "POSTS_DETAIL_INFO2")
	return t
}

func Test_node_getRoute(t1 *testing.T) {
	t := createRouter()

	tests := []struct {
		name       string
		gotPath    string
		wantNames  []ruleName
		wantParams map[string]string
	}{
		{"user list", "/user", []ruleName{"USER_LIST"}, map[string]string{}},
		{"user list", "/user/", []ruleName{"USER_LIST"}, map[string]string{}},
		{"user list", "user/", []ruleName{"USER_LIST"}, map[string]string{}},

		{"user detail", "/user/1", []ruleName{"USER_DETAIL"}, map[string]string{"id": "1"}},
		{"user detail", "/user/test/", []ruleName{"USER_DETAIL"}, map[string]string{"id": "test"}},

		{"user me", "/user/me/", []ruleName{"USER_ME"}, map[string]string{}},
		{"user me", "/user/me", []ruleName{"USER_ME"}, map[string]string{}},

		{"user top", "/user/1/top", []ruleName{"USER_TOP_LIST"}, map[string]string{"id": "1"}},
		{"user top", "/user/me/top", []ruleName{"USER_TOP_LIST"}, map[string]string{"id": "me"}},
		{"user top list", "/user/1/top/1", []ruleName{"USER_TOP_ITEM"}, map[string]string{"id": "1", "item": "1"}},

		{"posts info", "/posts/1/info", []ruleName{"POSTS_DETAIL_INFO", "POSTS_PASS", "POSTS_PASS2"}, map[string]string{"id": "1"}},
		{"posts last", "/posts/last/1", []ruleName{"POSTS_LAST", "POSTS_PASS"}, map[string]string{"id": "1"}},
		{"posts detail", "/posts/1", []ruleName{"POSTS_DETAIL", "POSTS_PASS", "POSTS_PASS2"}, map[string]string{"id": "1"}},

		{"posts pass", "/posts/1/inf/some/here", []ruleName{"POSTS_PASS", "POSTS_PASS2"}, map[string]string{}},
		{"posts pass 2", "/posts/1/inf", []ruleName{"POSTS_PASS", "POSTS_PASS2"}, map[string]string{}},
		{"posts pass 3", "/posts/", []ruleName{"POSTS_PASS"}, map[string]string{}},
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
		{"not find", "/use/me", "", map[string]string{}},
		{"not find", "/user/1/1", "", map[string]string{}},
		{"not find", "/user/1/2/top/", "", map[string]string{}},
		{"not find", "?", "", map[string]string{}},
	}

	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			_, _, err := t.resolve(tt.gotPath)
			if err == nil {
				t1.Errorf("resolve() error = %v, wantErr", err)
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
