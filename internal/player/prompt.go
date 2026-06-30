package player

import "sync"

// EmailCodePrompter asks the user for an OTP / email verification code.
type EmailCodePrompter func(email string) (string, error)

var (
	emailPromptMu sync.RWMutex
	emailPrompt   EmailCodePrompter
)

// SetEmailCodePrompt registers the interactive OTP prompt (desktop GUI).
// Safe for concurrent reads during scenario execution.
func SetEmailCodePrompt(fn EmailCodePrompter) {
	emailPromptMu.Lock()
	emailPrompt = fn
	emailPromptMu.Unlock()
}

func emailCodePrompter() EmailCodePrompter {
	emailPromptMu.RLock()
	fn := emailPrompt
	emailPromptMu.RUnlock()
	return fn
}
