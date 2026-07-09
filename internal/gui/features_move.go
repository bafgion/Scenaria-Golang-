package gui

func (s *Service) MoveFeature(src, destDir string) (string, error) {
	return s.fileOperator().MoveFeature(src, destDir)
}

func (s *Service) ImportFeatures(destDir string, paths []string) ([]string, error) {
	return s.fileOperator().ImportFeatures(destDir, paths)
}
