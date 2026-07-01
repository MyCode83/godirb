package urlutil

import "testing"

func TestJoinPath(t *testing.T) {
	tests := []struct {
		name    string
		fullURL string
		newPath string
		want    string
	}{
		{
			name:    "host without slash",
			fullURL: "https://target.com",
			newPath: "admin",
			want:    "https://target.com/admin",
		},
		{
			name:    "host with slash and new path with slash",
			fullURL: "https://target.com/",
			newPath: "/admin",
			want:    "https://target.com/admin",
		},
		{
			name:    "recursive path does not double slash",
			fullURL: "https://target.com/DIR/",
			newPath: "/FILE",
			want:    "https://target.com/DIR/FILE",
		},
		{
			name:    "path without trailing slash",
			fullURL: "https://target.com/DIR",
			newPath: "FILE",
			want:    "https://target.com/DIR/FILE",
		},
		{
			name:    "preserves encoded dot dot",
			fullURL: "https://target.com/%2e%2e/",
			newPath: "admin",
			want:    "https://target.com/%2e%2e/admin",
		},
		{
			name:    "preserves invalid percent escape",
			fullURL: "https://target.com/%invalid/",
			newPath: "admin",
			want:    "https://target.com/%invalid/admin",
		},
		{
			name:    "drops query and fragment from base",
			fullURL: "https://target.com/admin?x=1#frag",
			newPath: "login/",
			want:    "https://target.com/admin/login/",
		},
		{
			name:    "drops query from base",
			fullURL: "https://target.com/admin?x=1",
			newPath: "login",
			want:    "https://target.com/admin/login",
		},
		{
			name:    "drops fragment from base",
			fullURL: "https://target.com/admin#frag",
			newPath: "login",
			want:    "https://target.com/admin/login",
		},
		{
			name:    "drops query from final result",
			fullURL: "https://target.com/admin",
			newPath: "login?x=1#frag",
			want:    "https://target.com/admin/login",
		},
		{
			name:    "empty new path returns stripped base",
			fullURL: "https://target.com/admin?x=1#frag",
			newPath: "",
			want:    "https://target.com/admin",
		},
		{
			name:    "empty full url and empty path",
			fullURL: "",
			newPath: "",
			want:    "/",
		},
		{
			name:    "empty full url trims leading slashes",
			fullURL: "",
			newPath: "///admin",
			want:    "/admin",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JoinPath(tt.fullURL, tt.newPath)
			if err != nil {
				t.Fatalf("JoinPath() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("JoinPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAddExtension(t *testing.T) {
	tests := []struct {
		name    string
		fullURL string
		ext     string
		want    string
	}{
		{
			name:    "adds extension without dot",
			fullURL: "https://target.com/admin",
			ext:     "php",
			want:    "https://target.com/admin.php",
		},
		{
			name:    "adds extension with dot",
			fullURL: "https://target.com/admin",
			ext:     ".php",
			want:    "https://target.com/admin.php",
		},
		{
			name:    "trims extension spaces",
			fullURL: "https://target.com/admin",
			ext:     " php ",
			want:    "https://target.com/admin.php",
		},
		{
			name:    "removes trailing slash before extension",
			fullURL: "https://target.com/admin/",
			ext:     "php",
			want:    "https://target.com/admin.php",
		},
		{
			name:    "drops query and fragment",
			fullURL: "https://target.com/admin?x=1#frag",
			ext:     "php",
			want:    "https://target.com/admin.php",
		},
		{
			name:    "drops query",
			fullURL: "https://target.com/admin?x=1",
			ext:     "php",
			want:    "https://target.com/admin.php",
		},
		{
			name:    "drops fragment",
			fullURL: "https://target.com/admin#frag",
			ext:     "php",
			want:    "https://target.com/admin.php",
		},
		{
			name:    "drops query from final result",
			fullURL: "https://target.com/admin",
			ext:     "php?x=1#frag",
			want:    "https://target.com/admin.php",
		},
		{
			name:    "empty extension returns stripped base",
			fullURL: "https://target.com/admin?x=1#frag",
			ext:     " ",
			want:    "https://target.com/admin",
		},
		{
			name:    "preserves encoded dot dot",
			fullURL: "https://target.com/%2e%2e/admin",
			ext:     "php",
			want:    "https://target.com/%2e%2e/admin.php",
		},
		{
			name:    "preserves invalid percent escape",
			fullURL: "https://target.com/%invalid/admin",
			ext:     "php",
			want:    "https://target.com/%invalid/admin.php",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddExtension(tt.fullURL, tt.ext)
			if err != nil {
				t.Fatalf("AddExtension() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("AddExtension() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDropQueryAndFragment(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{
			name: "no query or fragment",
			raw:  "https://target.com/%2e%2e/admin",
			want: "https://target.com/%2e%2e/admin",
		},
		{
			name: "drops query",
			raw:  "https://target.com/%invalid/admin?x=1",
			want: "https://target.com/%invalid/admin",
		},
		{
			name: "drops fragment",
			raw:  "https://target.com/admin#frag",
			want: "https://target.com/admin",
		},
		{
			name: "query before fragment",
			raw:  "https://target.com/admin?x=1#frag",
			want: "https://target.com/admin",
		},
		{
			name: "2 queries and fragment",
			raw:  "https://target.com/admin?x=1#frag&x=4/",
			want: "https://target.com/admin",
		},
		{
			name: "fragment before query",
			raw:  "https://target.com/admin#frag?x=1",
			want: "https://target.com/admin",
		},
		{
			name: "raw path stays raw",
			raw:  "https://target.com/%2e%2e/%invalid",
			want: "https://target.com/%2e%2e/%invalid",
		},
		{
			name: "empty string",
			raw:  "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DropQueryAndFragment(tt.raw); got != tt.want {
				t.Fatalf("DropQueryAndFragment() = %q, want %q", got, tt.want)
			}
		})
	}
}
