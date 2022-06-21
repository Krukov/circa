package key_template

import "testing"

func Test_FormatTemplate(t *testing.T) {
	type args struct {
		template string
		params   map[string]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple", args{"key", map[string]string{}}, "key"},
		{"param", args{"key:{id}:test", map[string]string{"id": "1"}}, "key:1:test"},
		{"params", args{"key:{id}:{name}", map[string]string{"id": "1", "name": "test"}}, "key:1:test"},
		{"param-func", args{"key:{name|lower}", map[string]string{"name": "TEST"}}, "key:test"},
		{"param-strip", args{"key:{name|trim}", map[string]string{"name": " t "}}, "key:t"},
		{"param-strip2", args{"key:{name|trim:T}", map[string]string{"name": "TtT"}}, "key:t"},
		{"param-replace", args{"key:{name|replace:t}", map[string]string{"name": "TtT"}}, "key:TT"},
		{"param-miss", args{"key:{id}:{name}", map[string]string{"name": "test"}}, "key::test"},
		{"param-jwt", args{"key:{token|jwt:user_id}", map[string]string{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMiJ9.l-XkseImddERtMT8e3f-XxfPhhQhgplGvX3uPLYy1IQ"}}, "key:2"},
		{"param-hash", args{"key:{name|hash}", map[string]string{"name": "test"}}, "key:098f6bcd4621d373cade4e832627b4f6"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatTemplate(tt.args.template, tt.args.params); got != tt.want {
				t.Errorf("formatTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}
