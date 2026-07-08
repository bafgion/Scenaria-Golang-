package cli

import (
	"io"

	"github.com/bafgion/scenaria-golang/internal/update"
	"github.com/bafgion/scenaria-golang/internal/version"
)

func RunUpdate(args []string) error {
	return RunUpdateWithOutput(args, nil)
}

func RunUpdateWithOutput(args []string, out io.Writer) error {
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
		cliPrintf(out, "You are on the latest version (%s).\n", version.Version)
		return nil
	}
	cliPrintf(out, "Update available: %s -> %s\n", version.Version, release.TagName)
	cliPrintf(out, "Release: %s\n", release.HTMLURL)
	if checkOnly {
		return nil
	}
	cliPrintln(out, "Download the portable ZIP from the release page or run: make build-portable")
	return nil
}
