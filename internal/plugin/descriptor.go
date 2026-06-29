package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Descriptor struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Commands    []string `json:"commands"`
}

func DescriptorPath(projectRoot, pluginName string) string {
	return filepath.Join(projectRoot, "addons", pluginName, "plugin.json")
}

func LoadDescriptor(projectRoot, pluginName string) (Descriptor, error) {
	path := DescriptorPath(projectRoot, pluginName)
	payload, err := os.ReadFile(path)
	if err != nil {
		return Descriptor{}, fmt.Errorf("read plugin descriptor %q: %w", path, err)
	}
	var desc Descriptor
	if err := json.Unmarshal(payload, &desc); err != nil {
		return Descriptor{}, fmt.Errorf("decode plugin descriptor %q: %w", path, err)
	}
	if desc.ID == "" {
		desc.ID = pluginName
	}
	return desc, nil
}

func IsVanessa(desc Descriptor) bool {
	if strings.EqualFold(desc.ID, "vanessa") {
		return true
	}
	for _, cmd := range desc.Commands {
		lower := strings.ToLower(strings.TrimSpace(cmd))
		if lower == "va run" || strings.HasPrefix(lower, "va run ") {
			return true
		}
	}
	return false
}

func IsRunnable(desc Descriptor) bool {
	return IsVanessa(desc) || len(desc.Commands) > 0
}
