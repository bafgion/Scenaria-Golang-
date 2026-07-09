package gui

func (s *Service) SaveFeatureDraft(featurePath, content string) error {
	return s.fileOperator().SaveFeatureDraft(featurePath, content)
}

func (s *Service) LoadFeatureDraft(featurePath string) (string, error) {
	return s.fileOperator().LoadFeatureDraft(featurePath)
}

func (s *Service) ClearFeatureDraft(featurePath string) error {
	return s.fileOperator().ClearFeatureDraft(featurePath)
}
