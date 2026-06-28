package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Release struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

func LatestRelease(owner, repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", owner, repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 15 * time.Second}
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

func IsNewer(current, latest string) bool {
	return strings.TrimSpace(current) != "" && strings.TrimSpace(latest) != "" && current != latest
}
