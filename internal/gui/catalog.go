package gui

func (s *Service) SearchSteps(query string) []StepCatalogEntry {
	return s.catalogOps().Search(query)
}

func (s *Service) CompletionsForLine(line string, column int, language string) StepCompletionsDTO {
	return s.catalogOps().CompletionsForLine(line, column, language)
}

func (s *Service) catalogOps() *CatalogService {
	s.mu.RLock()
	ops := s.catalogService
	s.mu.RUnlock()
	if ops != nil {
		return ops
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.catalogService == nil {
		s.catalogService = NewCatalogService()
	}
	return s.catalogService
}
