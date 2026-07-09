package gui

import (
	"context"
	"fmt"
	"strings"
)

func (c *CLIOps) ValidateProject(ctx context.Context, projectPath string, req ValidateRequest) (string, error) {
	path := strings.TrimSpace(projectPath)
	if path == "" {
		return "", fmt.Errorf("open a project folder first")
	}
	args := []string{}
	if len(req.Targets) > 0 {
		args = append(args, req.Targets...)
	} else {
		args = append(args, path)
	}
	if req.SkipBrowser {
		args = append(args, "--no-browser")
	} else if req.Browser != "" {
		args = append(args, "--browser", req.Browser)
	}
	return c.Validate(ctx, args)
}
