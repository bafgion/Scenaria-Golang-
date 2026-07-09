package gui

import (
	"fmt"
	"strings"

	"github.com/bafgion/scenaria-golang/internal/httpauth"
	"github.com/bafgion/scenaria-golang/internal/selector"
	"github.com/bafgion/scenaria-golang/internal/settings"
)

// SettingsService owns persisted app settings access.
type SettingsService struct {
	store *settings.Store
}

func NewSettingsService(store *settings.Store) *SettingsService {
	return &SettingsService{store: store}
}

func (s *SettingsService) loadAppSettings() (*settings.AppSettings, error) {
	if s == nil || s.store == nil {
		return &settings.AppSettings{Browser: "chromium"}, nil
	}
	cfg, err := s.store.Load()
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &settings.AppSettings{Browser: "chromium"}
	}
	return cfg, nil
}

func (s *SettingsService) saveAppSettings(cfg *settings.AppSettings) error {
	if s == nil || s.store == nil {
		return nil
	}
	return s.store.Update(func(current *settings.AppSettings) error {
		*current = *cfg
		return nil
	})
}

func (s *SettingsService) LoadRecents() RecentsDTO {
	cfg, err := s.loadAppSettings()
	if err != nil || cfg == nil {
		return RecentsDTO{}
	}
	return RecentsDTO{
		Projects: trimRecents(cfg.RecentProjects),
		Features: trimRecents(cfg.RecentFeatures),
	}
}

func (s *SettingsService) RememberRecentProject(path string) error {
	return s.rememberRecent(path, true)
}

func (s *SettingsService) RememberRecentFeature(path string) error {
	return s.rememberRecent(path, false)
}

func (s *SettingsService) rememberRecent(itemPath string, project bool) error {
	itemPath = strings.TrimSpace(itemPath)
	if itemPath == "" {
		return nil
	}
	if s == nil || s.store == nil {
		return nil
	}
	return s.store.Update(func(cfg *settings.AppSettings) error {
		if project {
			cfg.RecentProjects = pushRecent(cfg.RecentProjects, itemPath)
		} else {
			cfg.RecentFeatures = pushRecent(cfg.RecentFeatures, itemPath)
		}
		return nil
	})
}

func (s *SettingsService) ListHTTPAuthHosts() ([]string, error) {
	cfg, err := s.loadAppSettings()
	if err != nil {
		return nil, err
	}
	return httpauth.ListHosts(cfg), nil
}

func (s *SettingsService) HTTPAuthForHost(host string) (HTTPAuthCredentials, error) {
	cfg, err := s.loadAppSettings()
	if err != nil {
		return HTTPAuthCredentials{}, err
	}
	username, password := httpauth.CredentialsForHost(host, cfg)
	return HTTPAuthCredentials{Username: username, HasPassword: password != ""}, nil
}

func (s *SettingsService) SaveHTTPAuth(req HTTPAuthRequest) error {
	host := strings.TrimSpace(req.Host)
	if host == "" {
		return fmt.Errorf("host is required")
	}
	if s == nil || s.store == nil {
		return fmt.Errorf("settings store unavailable")
	}
	return s.store.Update(func(cfg *settings.AppSettings) error {
		password := req.Password
		if strings.TrimSpace(password) == "" {
			if _, existing := httpauth.CredentialsForHost(host, cfg); existing != "" {
				password = existing
			}
		}
		httpauth.StoreHostCredentials(host, req.Username, password, cfg)
		return nil
	})
}

func (s *SettingsService) RemoveHTTPAuth(host string) error {
	if s == nil || s.store == nil {
		return nil
	}
	return s.store.Update(func(cfg *settings.AppSettings) error {
		httpauth.RemoveHostCredentials(host, cfg)
		return nil
	})
}

func (s *SettingsService) PrepareRecordURL(url string) (string, error) {
	if s == nil || s.store == nil {
		return url, nil
	}
	var clean string
	err := s.store.Update(func(cfg *settings.AppSettings) error {
		clean = httpauth.ApplyURLCredentials(url, cfg)
		return nil
	})
	if err != nil {
		return url, err
	}
	return clean, nil
}

func (s *SettingsService) LoadSettings() (AppSettingsDTO, error) {
	if s == nil || s.store == nil {
		return defaultAppSettingsDTO(), nil
	}
	cfg, err := s.store.Load()
	if err != nil || cfg == nil {
		return defaultAppSettingsDTO(), nil
	}
	return appSettingsFromCfg(cfg), nil
}

func (s *SettingsService) SaveSettings(dto AppSettingsDTO) error {
	if s == nil || s.store == nil {
		return nil
	}
	height := dto.StepsPanelHeight
	if height < 80 {
		height = 160
	}
	return s.store.Update(func(cfg *settings.AppSettings) error {
		existing := *cfg
		navWaitUntil := strings.TrimSpace(dto.NavWaitUntil)
		sessionProject := strings.TrimSpace(dto.SessionProject)
		activeTab := strings.TrimSpace(dto.ActiveTab)
		startURL := strings.TrimSpace(dto.StartURL)
		next := &settings.AppSettings{
			Browser:                 dto.Browser,
			Headless:                dto.Headless,
			ParallelWorkers:         maxInt(1, dto.ParallelWorkers),
			SlowMo:                  maxInt(0, dto.SlowMo),
			MaxLoopIterations:       maxInt(1, dto.MaxLoopIterations),
			NavWaitUntil:            navWaitUntil,
			RecordingFilterMode:     dto.FilterRecording,
			NavOnlyRecording:        dto.NavOnlyRecording,
			RecordingHoverMode:      dto.HoverRecord,
			ToolbarCompact:          dto.ToolbarCompact,
			StepsPanelVisible:       dto.StepsPanelVisible,
			StepsPanelHeight:        height,
			SidebarWidth:            clampSidebarWidth(dto.SidebarWidth),
			RecentProjects:          trimRecents(dto.RecentProjects),
			RecentFeatures:          trimRecents(dto.RecentFeatures),
			SessionProject:          sessionProject,
			OpenTabs:                trimRecents(dto.OpenTabs),
			ActiveTab:               activeTab,
			UntitledTabs:            untitledTabsToCfg(dto.UntitledTabs),
			ScrollBeforeClick:       dto.ScrollBeforeClick,
			HoverRecordMinMs:        normalizeHoverRecordMinMs(dto.HoverRecordMinMs),
			SelectorClickStrategies: selector.NormalizeClickStrategies(dto.SelectorClickStrategies),
			SelectorInputStrategies: selector.NormalizeInputStrategies(dto.SelectorInputStrategies),
			Editor:                  settings.NormalizeEditorSettings(dto.Editor),
			ChecklistDismissed:      dto.ChecklistDismissed,
			WelcomePlayedSuccess:    dto.WelcomePlayedSuccess,
			OnboardingCompleted:     dto.OnboardingCompleted,
			OnboardingDismissed:     dto.OnboardingDismissed,
			OnboardingVersion:       dto.OnboardingVersion,
			StartURL:                startURL,
			RunDialogConfirmed:      dto.RunDialogConfirmed,
			PickerDuringRecording:   dto.PickerDuringRecording,
			UILocale:                normalizeUILocale(dto.UILocale),
		}
		checkUpdates := dto.CheckUpdatesOnStartup
		next.CheckUpdatesOnStartup = &checkUpdates
		next.HTTPAuth = existing.HTTPAuth
		if len(next.RecentProjects) == 0 {
			next.RecentProjects = existing.RecentProjects
		}
		if len(next.RecentFeatures) == 0 {
			next.RecentFeatures = existing.RecentFeatures
		}
		sidebarW := dto.SidebarWidth
		if sidebarW < 120 && existing.SidebarWidth >= 120 {
			sidebarW = existing.SidebarWidth
		}
		next.SidebarWidth = clampSidebarWidth(sidebarW)
		*cfg = *next
		return nil
	})
}
