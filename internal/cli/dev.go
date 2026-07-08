package cli

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
)

func RunDev(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: scenaria dev <goroutines|pprof> [flags]")
	}
	switch args[0] {
	case "goroutines":
		return runDevGoroutines()
	case "pprof":
		return runDevPprof(args[1:])
	default:
		return fmt.Errorf("unknown dev subcommand: %s", args[0])
	}
}

func runDevGoroutines() error {
	buf := make([]byte, 1<<20)
	for {
		n := runtime.Stack(buf, true)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	fmt.Println(string(buf))
	fmt.Printf("\n-- total: %d goroutines\n", runtime.NumGoroutine())
	return nil
}

func runDevPprof(args []string) error {
	outPath := "goroutine.pprof"
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--out":
			if i+1 >= len(args) {
				return fmt.Errorf("pprof requires --out <path>")
			}
			outPath = args[i+1]
			i++
		default:
			return fmt.Errorf("unknown pprof flag: %s", args[i])
		}
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("create pprof output: %w", err)
	}
	defer f.Close()
	if err := pprof.Lookup("goroutine").WriteTo(f, 0); err != nil {
		return fmt.Errorf("write goroutine profile: %w", err)
	}
	fmt.Printf("Wrote goroutine profile to %s (%d goroutines)\n", outPath, runtime.NumGoroutine())
	return nil
}
