package cli

import (
	"fmt"
	"io"
	"os"
)

func cliWriter(w io.Writer) io.Writer {
	if w == nil {
		return os.Stdout
	}
	return w
}

func cliPrintf(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(cliWriter(w), format, args...)
}

func cliPrintln(w io.Writer, args ...any) {
	_, _ = fmt.Fprintln(cliWriter(w), args...)
}
