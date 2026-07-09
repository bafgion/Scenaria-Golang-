package gui

type PluginEntryDTO struct {
	Name        string   `json:"name"`
	Source      string   `json:"source"`
	InstalledAt string   `json:"installedAt"`
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Commands    []string `json:"commands"`
	Runnable    bool     `json:"runnable"`
	Vanessa     bool     `json:"vanessa"`
}

func (s *Service) ListPlugins() ([]PluginEntryDTO, error) {
	return s.pluginOps().List()
}

func (s *Service) RunPlugin(req PluginRunRequest) RunResult {
	return s.pluginOps().Run(req)
}

// RunVanessa keeps backward compatibility for simple dry/run calls.
func (s *Service) RunVanessa(dryRun bool) RunResult {
	return s.pluginOps().RunVanessa(dryRun)
}

func (s *Service) InstallPlugin(name, source string) error {
	return s.pluginOps().Install(name, source)
}

func (s *Service) UninstallPlugin(name string) error {
	return s.pluginOps().Uninstall(name)
}

func (s *Service) pluginOps() *PluginService {
	s.mu.RLock()
	ops := s.pluginService
	s.mu.RUnlock()
	if ops != nil {
		return ops
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pluginService == nil {
		s.pluginService = NewPluginService(s.ProjectPath, func() pluginCLIRunner { return s.cliRunner() })
	}
	return s.pluginService
}
