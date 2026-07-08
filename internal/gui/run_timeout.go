package gui

import "time"

// DefaultRunTimeout caps a single GUI Playwright run so a stuck site cannot block the IDE forever.
const DefaultRunTimeout = 20 * time.Minute

// DefaultValidateTimeout caps GUI validation jobs.
const DefaultValidateTimeout = 2 * time.Minute
