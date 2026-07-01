package player

import (
	"errors"
)

// ExecutionFailure is returned when a run stops early but partial scenario results are available.
type ExecutionFailure struct {
	Err    error
	Result ExecutionResult
}

func (e *ExecutionFailure) Error() string {
	if e == nil || e.Err == nil {
		return "execution failed"
	}
	return e.Err.Error()
}

func (e *ExecutionFailure) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

// NewExecutionFailure wraps err with partial execution results.
func NewExecutionFailure(err error, result ExecutionResult) error {
	if err == nil {
		return nil
	}
	return &ExecutionFailure{Err: err, Result: result}
}

// PartialExecutionResult extracts results from an ExecutionFailure.
func PartialExecutionResult(err error) (ExecutionResult, bool) {
	var fail *ExecutionFailure
	if errors.As(err, &fail) {
		return fail.Result, true
	}
	return ExecutionResult{}, false
}

func executionFailure(err error, result ExecutionResult) error {
	if len(result.ScenarioResults) == 0 {
		return err
	}
	return NewExecutionFailure(err, result)
}
