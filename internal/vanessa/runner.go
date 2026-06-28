package vanessa

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type CaseResult struct {
	Path    string
	Name    string
	Success bool
	Message string
}

type BatchResult struct {
	Success  bool
	Cases    []CaseResult
	RunDir   string
	ExitCode int
	Error    string
}

func Run(req RunRequest) (BatchResult, error) {
	projectRoot := req.ProjectRoot
	if projectRoot == "" && len(req.Paths) > 0 {
		projectRoot = filepath.Dir(req.Paths[0])
	}
	cfg, err := LoadSettings(projectRoot)
	if err != nil {
		return BatchResult{}, err
	}
	if issues := ValidateSettings(cfg); len(issues) > 0 {
		return BatchResult{Success: false, Error: strings.Join(issues, "; ")}, fmt.Errorf("%s", strings.Join(issues, "; "))
	}
	if req.DryRun || cfg.DryRunOnly {
		files, ferr := resolveFeatureFiles(req)
		if ferr != nil {
			return BatchResult{}, ferr
		}
		cases := make([]CaseResult, 0, len(files))
		for _, file := range files {
			cases = append(cases, CaseResult{Path: file, Name: filepath.Base(file), Success: true, Message: "dry-run"})
		}
		return BatchResult{Success: true, Cases: cases}, nil
	}

	runDir := filepath.Join(ResolveRunsDir(cfg), fmt.Sprintf("run-%d", time.Now().UnixNano()))
	if err := os.MkdirAll(runDir, 0o755); err != nil {
		return BatchResult{}, err
	}
	_, vaPath, err := MergeVAParams(cfg, req, runDir)
	if err != nil {
		return BatchResult{}, err
	}

	cmd := buildPlatformCommand(cfg, vaPath)
	processLog := filepath.Join(runDir, "process.log")
	logFile, err := os.Create(processLog)
	if err != nil {
		return BatchResult{}, err
	}
	defer logFile.Close()
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		return BatchResult{RunDir: runDir, Error: err.Error()}, err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()

	timeout := time.Duration(cfg.ProcessTimeoutSec) * time.Second
	if timeout <= 0 {
		timeout = time.Hour
	}
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	var waitErr error
	select {
	case waitErr = <-done:
	case <-timer.C:
		_ = cmd.Process.Kill()
		waitErr = fmt.Errorf("process timeout after %s", timeout)
	}

	exitCode := 0
	if waitErr != nil {
		if exitErr, ok := waitErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return BatchResult{RunDir: runDir, ExitCode: -1, Error: waitErr.Error()}, waitErr
		}
	}

	junitDir := filepath.Join(runDir, "junit")
	cases := ParseJUnitDir(junitDir)
	if len(cases) == 0 {
		files, _ := resolveFeatureFiles(req)
		success := exitCode == 0
		for _, file := range files {
			cases = append(cases, CaseResult{
				Path:    file,
				Name:    strings.TrimSuffix(filepath.Base(file), ".feature"),
				Success: success,
				Message: fmt.Sprintf("exit code %d", exitCode),
			})
		}
	}
	success := exitCode == 0
	for _, c := range cases {
		if !c.Success {
			success = false
			break
		}
	}
	return BatchResult{
		Success:  success,
		Cases:    cases,
		RunDir:   runDir,
		ExitCode: exitCode,
	}, nil
}

func buildPlatformCommand(cfg Settings, vaParamsPath string) *exec.Cmd {
	args := []string{cfg.PlatformMode}
	if conn := strings.TrimSpace(cfg.IBConnection); conn != "" {
		args = append(args, conn)
	}
	if user := strings.TrimSpace(cfg.User); user != "" {
		args = append(args, "/N"+user)
	}
	if pass := cfg.Password; pass != "" {
		args = append(args, "/P"+pass)
	}
	args = append(args, `/Execute`+cfg.EPFPath)
	args = append(args, `/C`+vaParamsPath)
	args = append(args, cfg.PlatformExtraArgs...)
	return exec.Command(cfg.PlatformExecutable, args...)
}

func ParseJUnitDir(dir string) []CaseResult {
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		return nil
	}
	out := make([]CaseResult, 0)
	_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".xml" {
			return nil
		}
		out = append(out, parseJUnitFile(path)...)
		return nil
	})
	return out
}

func parseJUnitFile(path string) []CaseResult {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()
	payload, err := io.ReadAll(file)
	if err != nil {
		return nil
	}
	var root junitNode
	if err := xml.Unmarshal(payload, &root); err != nil {
		return nil
	}
	suites := root.TestSuites
	if root.Tests > 0 || len(root.TestCases) > 0 {
		suites = []junitNode{root}
	}
	out := make([]CaseResult, 0)
	for _, suite := range suites {
		for _, tc := range suite.TestCases {
			success := tc.Failure == nil && tc.Error == nil
			message := "ok"
			if tc.Failure != nil {
				message = strings.TrimSpace(tc.Failure.Message)
				if message == "" {
					message = strings.TrimSpace(tc.Failure.Body)
				}
			}
			if tc.Error != nil {
				success = false
				message = strings.TrimSpace(tc.Error.Message)
			}
			out = append(out, CaseResult{
				Path:    path,
				Name:    tc.Name,
				Success: success,
				Message: message,
			})
		}
	}
	return out
}

type junitNode struct {
	XMLName    xml.Name    `xml:"testsuite"`
	Tests      int         `xml:"tests,attr"`
	TestCases  []junitCase `xml:"testcase"`
	TestSuites []junitNode `xml:"testsuite"`
}

type junitCase struct {
	Name    string       `xml:"name,attr"`
	Failure *junitFault `xml:"failure"`
	Error   *junitFault `xml:"error"`
}

type junitFault struct {
	Message string `xml:"message,attr"`
	Body    string `xml:",chardata"`
}
