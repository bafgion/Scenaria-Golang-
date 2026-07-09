package gui

func (s *Service) RenameFeature(path, newName string) (string, error) {
	return s.fileOperator().RenameFeature(path, newName)
}
