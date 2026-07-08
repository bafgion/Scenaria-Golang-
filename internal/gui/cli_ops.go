package gui

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/cli"
)

// CLIOps centralizes GUI calls into CLI command handlers with captured output.
type CLIOps struct {
	runInitFn       func(args []string, out io.Writer) error
	runValidateFn   func(ctx context.Context, args []string, out io.Writer) error
	runUpdateFn     func(args []string, out io.Writer) error
	runExportFn     func(args []string, out io.Writer) error
	runImportJSONFn func(args []string, out io.Writer) error
	runRecordFn     func(args []string, out io.Writer) error
	runVAFn         func(args []string, out io.Writer) error
	runRunFn        func(ctx context.Context, args []string, out io.Writer) error
}

func NewCLIOps() *CLIOps {
	return &CLIOps{
		runInitFn:       cli.RunInitWithOutput,
		runValidateFn:   cli.RunValidateContextWithOutput,
		runUpdateFn:     cli.RunUpdateWithOutput,
		runExportFn:     cli.RunExportWithOutput,
		runImportJSONFn: cli.RunImportJSONWithOutput,
		runRecordFn:     cli.RunRecordWithOutput,
		runVAFn:         cli.RunVAWithOutput,
		runRunFn:        cli.RunRunContextWithOutput,
	}
}

func runCLIWithOutput(fn func(io.Writer) error) (out string, runErr error) {
	var buf strings.Builder
	defer func() {
		if recovered := recover(); recovered != nil {
			runErr = fmt.Errorf("cli panic: %v", recovered)
		}
	}()
	runErr = fn(&buf)
	return buf.String(), runErr
}

func (c *CLIOps) InitProject(path string) (string, error) {
	return runCLIWithOutput(func(out io.Writer) error { return c.runInitFn([]string{path}, out) })
}

func (c *CLIOps) Validate(ctx context.Context, args []string) (string, error) {
	return runCLIWithOutput(func(out io.Writer) error { return c.runValidateFn(ctx, args, out) })
}

func (c *CLIOps) CheckUpdates() (string, error) {
	return runCLIWithOutput(func(out io.Writer) error { return c.runUpdateFn([]string{"--check"}, out) })
}

func (c *CLIOps) Export(args []string) (string, error) {
	return runCLIWithOutput(func(out io.Writer) error { return c.runExportFn(args, out) })
}

func (c *CLIOps) ImportJSON(args []string) (string, error) {
	return runCLIWithOutput(func(out io.Writer) error { return c.runImportJSONFn(args, out) })
}

func (c *CLIOps) Record(args []string) (string, error) {
	return runCLIWithOutput(func(out io.Writer) error { return c.runRecordFn(args, out) })
}

func (c *CLIOps) VA(args []string) (string, error) {
	return runCLIWithOutput(func(out io.Writer) error { return c.runVAFn(args, out) })
}

func (c *CLIOps) Run(ctx context.Context, args []string) (string, error) {
	return runCLIWithOutput(func(out io.Writer) error { return c.runRunFn(ctx, args, out) })
}
