package player

import (
	"context"
	"sync"
)

// EmailCodePrompter asks the user for an OTP / email verification code.
type EmailCodePrompter func(email string) (string, error)

var (
	emailPromptMu sync.RWMutex
	emailPrompt   EmailCodePrompter

	otpCancelMu   sync.RWMutex
	otpCancelHook func()
)

// SetEmailCodePrompt registers the interactive OTP prompt (desktop GUI).
// Safe for concurrent reads during scenario execution.
func SetEmailCodePrompt(fn EmailCodePrompter) {
	emailPromptMu.Lock()
	emailPrompt = fn
	emailPromptMu.Unlock()
}

// SetOTPCancelHook registers a callback that aborts a blocking OTP prompt (GUI CancelOTP).
func SetOTPCancelHook(fn func()) {
	otpCancelMu.Lock()
	otpCancelHook = fn
	otpCancelMu.Unlock()
}

func emailCodePrompter() EmailCodePrompter {
	emailPromptMu.RLock()
	fn := emailPrompt
	emailPromptMu.RUnlock()
	return fn
}

func cancelPendingOTP() {
	otpCancelMu.RLock()
	fn := otpCancelHook
	otpCancelMu.RUnlock()
	if fn != nil {
		fn()
	}
}

// ctxAwareEmailPrompt wraps an OTP prompt so scenario cancellation unblocks the GUI dialog.
func ctxAwareEmailPrompt(ctx context.Context, prompt EmailCodePrompter) EmailCodePrompter {
	if prompt == nil || ctx == nil {
		return prompt
	}
	return func(email string) (string, error) {
		type result struct {
			code string
			err  error
		}
		ch := make(chan result, 1)
		go func() {
			code, err := prompt(email)
			ch <- result{code, err}
		}()
		select {
		case <-ctx.Done():
			cancelPendingOTP()
			go func() { _ = <-ch }()
			return "", ctx.Err()
		case r := <-ch:
			if ctx.Err() != nil {
				return "", ctx.Err()
			}
			return r.code, r.err
		}
	}
}
