package update

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsNewer(t *testing.T) {
	if !IsNewer("0.1.0", "0.2.0") {
		t.Fatal("expected different versions to be newer")
	}
	if IsNewer("", "0.2.0") {
		t.Fatal("empty current should not compare as newer")
	}
	if IsNewer("1.0.0", "1.0.0") {
		t.Fatal("same version is not newer")
	}
	if IsNewer("2.0.0", "") {
		t.Fatal("empty latest should not compare as newer")
	}
}

func TestFetchRelease(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method %s", r.Method)
		}
		_ = json.NewEncoder(w).Encode(Release{
			TagName: " v1.2.3 ",
			HTMLURL: "https://example.com/releases/v1.2.3",
		})
	}))
	defer srv.Close()

	release, err := fetchRelease(srv.Client(), srv.URL)
	if err != nil {
		t.Fatalf("fetchRelease: %v", err)
	}
	if release.TagName != "v1.2.3" {
		t.Fatalf("tag=%q", release.TagName)
	}
	if release.HTMLURL != "https://example.com/releases/v1.2.3" {
		t.Fatalf("url=%q", release.HTMLURL)
	}
}

func TestFetchReleaseNonOK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	if _, err := fetchRelease(srv.Client(), srv.URL); err == nil {
		t.Fatal("expected error for 404")
	}
}
