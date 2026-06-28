package cli

import (
	"fmt"

	"github.com/bafgion/scenaria-golang/internal/update"
	"github.com/bafgion/scenaria-golang/internal/version"
)

func RunUpdate(args []string) error {
	checkOnly := false
	for _, arg := range args {
		if arg == "--check" {
			checkOnly = true
		}
	}

	release, err := update.LatestRelease("bafgion", "Scenaria-Golang-")
	if err != nil {
		return err
	}
	if !update.IsNewer(version.Version, release.TagName) {
		fmt.Printf("You are on the latest version (%s).\n", version.Version)
		return nil
	}
	fmt.Printf("Update available: %s -> %s\n", version.Version, release.TagName)
	fmt.Printf("Release: %s\n", release.HTMLURL)
	if checkOnly {
		return nil
	}
	fmt.Println("Download the portable ZIP from the release page or run: make build-portable")
	return nil
}
