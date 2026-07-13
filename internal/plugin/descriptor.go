package plugin

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Descriptor struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Version        string        `json:"version"`
	Description    string        `json:"description"`
	Commands       []string      `json:"commands"`
	StructuredRuns []CommandSpec `json:"structuredRuns,omitempty"`
	CommandSpecs   []CommandSpec `json:"commandSpecs,omitempty"`
}

type CommandSpec struct {
	Runner string   `json:"runner"`
	Args   []string `json:"args,omitempty"`
}

func DescriptorPath(projectRoot, pluginName string) string {
	path, err := descriptorPath(projectRoot, pluginName)
	if err != nil {
		return ""
	}
	return path
}

func LoadDescriptor(projectRoot, pluginName string) (Descriptor, error) {
	path, err := descriptorPath(projectRoot, pluginName)
	if err != nil {
		return Descriptor{}, err
	}
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

func descriptorPath(projectRoot, pluginName string) (string, error) {
	dir, err := addonPath(projectRoot, pluginName)
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "plugin.json"), nil
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
	for _, cmd := range structuredCommands(desc) {
		if strings.EqualFold(strings.TrimSpace(cmd.Runner), "va") {
			return true
		}
	}
	return false
}

func IsRunnable(desc Descriptor) bool {
	return IsVanessa(desc) || len(desc.Commands) > 0 || len(structuredCommands(desc)) > 0
}

func structuredCommands(desc Descriptor) []CommandSpec {
	if len(desc.StructuredRuns) > 0 {
		return desc.StructuredRuns
	}
	return desc.CommandSpecs
}
