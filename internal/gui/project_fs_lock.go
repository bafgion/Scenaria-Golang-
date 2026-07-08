package gui

func (s *Service) withProjectFSReadLock(fn func() error) error {
	s.projectFSMu.RLock()
	defer s.projectFSMu.RUnlock()
	return fn()
}

func (s *Service) withProjectFSWriteLock(fn func() error) error {
	s.projectFSMu.Lock()
	defer s.projectFSMu.Unlock()
	return fn()
}
