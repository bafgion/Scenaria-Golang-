package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/mod/semver"
)

type Release struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func LatestRelease(owner, repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	return fetchRelease(&http.Client{Timeout: 15 * time.Second}, url)
}

func fetchRelease(client *http.Client, url string) (*Release, error) {
	if client == nil {
		client = &http.Client{Timeout: 15 * time.Second}
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch latest release: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github release API returned %s", resp.Status)
	}
	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("decode release metadata: %w", err)
	}
	release.TagName = strings.TrimSpace(release.TagName)
	return &release, nil
}

// canonicalVersionTag normalizes app and GitHub release versions for comparison.
func canonicalVersionTag(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	if !strings.HasPrefix(strings.ToLower(v), "v") {
		v = "v" + v
	}
	return v
}

// DisplayVersion strips a leading v/V for UI messages.
func DisplayVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	v = strings.TrimPrefix(v, "V")
	return v
}

func IsNewer(current, latest string) bool {
	cur := canonicalVersionTag(current)
	lat := canonicalVersionTag(latest)
	if cur == "" || lat == "" {
		return false
	}
	if semver.IsValid(cur) && semver.IsValid(lat) {
		return semver.Compare(cur, lat) < 0
	}
	return DisplayVersion(current) != DisplayVersion(latest)
}
