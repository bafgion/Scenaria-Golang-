package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bafgion/scenaria-golang/internal/settings"
)

func RunInit(args []string) error {
	target := "."
	if len(args) > 0 {
		target = args[0]
	}
	info, err := os.Stat(target)
	if err != nil || !info.IsDir() {
		return fmt.Errorf("init target must be an existing directory")
	}

	cfg := settings.DefaultProjectConfig()
	if err := settings.SaveProjectConfig(target, cfg); err != nil {
		return err
	}

	scenariaDir := filepath.Join(target, ".scenaria")
	_ = os.MkdirAll(filepath.Join(scenariaDir, "test_clients"), 0o755)
	_ = os.MkdirAll(filepath.Join(scenariaDir, "downloads"), 0o755)

	vanessaExample := filepath.Join(scenariaDir, "vanessa.json.example")
	_ = os.WriteFile(vanessaExample, []byte(`{
  "platform_executable": "C:\\Program Files\\1cv8\\bin\\1cv8.exe",
  "epf_path": "C:\\path\\to\\vanessa-automation.epf",
  "ib_connection_string": "/F\"C:\\base\\1cv8\""
}
`), 0o644)

	vaParamsExample := filepath.Join(scenariaDir, "va-params.base.json.example")
	_ = os.WriteFile(vaParamsExample, []byte(`{
  "ЗапускатьСценарии": true,
  "ЗакрытьБраузерПослеВыполненияСценариев": true
}
`), 0o644)

	demoClient := filepath.Join(scenariaDir, "test_clients", "DemoUser.json.example")
	_ = os.WriteFile(demoClient, []byte(`{
  "name": "DemoUser",
  "base_url": "https://example.com",
  "cookies": [],
  "local_storage": {}
}
`), 0o644)

	fmt.Printf("Initialized Scenaria project in %s\n", scenariaDir)
	fmt.Println("  .scenaria/project.json")
	fmt.Println("  .scenaria/vanessa.json.example")
	fmt.Println("  .scenaria/va-params.base.json.example")
	fmt.Println("  .scenaria/test_clients/DemoUser.json.example")
	return nil
}
