package settings

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func TestClientsDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".scenaria", "test_clients")
}

// ListTestClientNames returns sorted JSON test client names (without extension).
func ListTestClientNames(projectRoot string) ([]string, error) {
	if projectRoot == "" {
		return nil, fmt.Errorf("project root is required")
	}
	dir := TestClientsDir(projectRoot)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".json") {
			continue
		}
		if strings.HasSuffix(strings.ToLower(name), ".json.example") {
			continue
		}
		names = append(names, strings.TrimSuffix(name, filepath.Ext(name)))
	}
	sortStrings(names)
	return names, nil
}

func LoadTestClientByName(projectRoot, name string) (*TestClient, error) {
	path, err := TestClientPath(projectRoot, name)
	if err != nil {
		return nil, err
	}
	return LoadTestClient(path)
}

func sortStrings(values []string) {
	for i := 0; i < len(values); i++ {
		for j := i + 1; j < len(values); j++ {
			if values[j] < values[i] {
				values[i], values[j] = values[j], values[i]
			}
		}
	}
}
