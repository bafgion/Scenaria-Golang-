package gui

import (
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
		next := &settings.AppSettings{
			Browser:                 dto.Browser,
			Headless:                dto.Headless,
			ParallelWorkers:         maxInt(1, dto.ParallelWorkers),
			SlowMo:                  maxInt(0, dto.SlowMo),
			MaxLoopIterations:       maxInt(1, dto.MaxLoopIterations),
			NavWaitUntil:            dto.NavWaitUntil,
			RecordingFilterMode:     dto.FilterRecording,
			NavOnlyRecording:        dto.NavOnlyRecording,
			RecordingHoverMode:      dto.HoverRecord,
			ToolbarCompact:          dto.ToolbarCompact,
			StepsPanelVisible:       dto.StepsPanelVisible,
			StepsPanelHeight:        height,
			SidebarWidth:            clampSidebarWidth(dto.SidebarWidth),
			RecentProjects:          trimRecents(dto.RecentProjects),
			RecentFeatures:          trimRecents(dto.RecentFeatures),
			SessionProject:          dto.SessionProject,
			OpenTabs:                trimRecents(dto.OpenTabs),
			ActiveTab:               dto.ActiveTab,
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
			StartURL:                dto.StartURL,
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
