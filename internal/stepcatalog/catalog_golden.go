package stepcatalog

import (
	"strings"
	"sync"

	"github.com/bafgion/scenaria-golang/internal/stepdsl"
)

var (
	entriesOnce sync.Once
	allEntries  []Entry
)

func loadEntries() []Entry {
	entriesOnce.Do(func() {
		allEntries = buildEntries()
	})
	return allEntries
}

func buildEntries() []Entry {
	seen := map[string]bool{}
	out := make([]Entry, 0, 96)

	add := func(entry Entry) {
		if entry.Template == "" || seen[entry.Template] {
			return
		}
		seen[entry.Template] = true
		if entry.Example == "" {
			entry.Example = entry.Template
		}
		if entry.Label == "" {
			entry.Label = labelFromTemplate(entry.Template)
		}
		if entry.Description == "" {
			entry.Description = entry.Help
		}
		if entry.Help == "" {
			entry.Help = entry.Description
		}
		out = append(out, entry)
	}

	for _, entry := range buildSnippetEntries() {
		add(entry)
	}

	cases, err := stepdsl.LoadGoldenCases()
	if err == nil {
		for _, c := range cases {
			desc := helpForKind(c.Kind)
			add(Entry{
				Action:      c.Kind,
				Category:    categoryForKind(c.Kind),
				Description: desc,
				Template:    c.Text,
				Example:     c.Text,
				Parameters:  paramsForKind(c.Kind),
				Help:        desc,
			})
		}
	}

	for _, entry := range manualEntries() {
		add(entry)
	}
	return out
}

func labelFromTemplate(template string) string {
	line := strings.Split(template, "\n")[0]
	line = strings.TrimSpace(line)
	if len(line) > 48 {
		return line[:45] + "…"
	}
	return line
}

func manualEntries() []Entry {
	desc := func(text string) string { return text }
	return []Entry{
		entryManual("Сессия и TestClient", `я подключаю TestClient "имя"`, "test_client", desc("Именованный профиль cookies/localStorage из .scenaria/test_clients (блок «Контекст»)")),
		entryManual("Структура", "Если вижу \"#modal\"", "if", desc("Условный блок: шаги внутри выполняются, если элемент виден")),
		entryManual("Структура", "Повторяю 3 раза", "repeat", desc("Цикл с фиксированным числом повторений")),
		entryManual("Структура", "Пока вижу \"#next\"", "while", desc("Цикл «пока условие истинно»")),
		entryManual("Структура", "Для каждого \"#row\" как \"item\"", "for_each", desc("Цикл по каждому элементу списка")),
		entryManual("Ожидание", "жду 2 с", "wait", desc("Пауза в секундах")),
		entryManual("Генераторы", "ввожу случайный огрнип в \"#ogrnip\"", "fill_generated", desc("Случайный ОГРНИП")),
		entryManual("Генераторы", "ввожу случайное фио в \"#name\"", "fill_generated", desc("Случайное ФИО")),
	}
}

func entryManual(category, template, action, description string) Entry {
	return Entry{
		Label:       labelFromTemplate(template),
		Action:      action,
		Category:    category,
		Description: description,
		Template:    template,
		Example:     template,
		Parameters:  actionParameters(action),
		Help:        description,
	}
}

func paramsForKind(kind string) []string {
	switch kind {
	case "goto":
		return actionParameters("goto")
	case "click", "double-click":
		return actionParameters("click")
	case "hover":
		return actionParameters("hover")
	case "fill", "fill-generated":
		return actionParameters("fill")
	case "select":
		return actionParameters("select")
	case "upload":
		return actionParameters("upload")
	case "clear":
		return actionParameters("clear")
	case "check", "uncheck":
		return actionParameters("check")
	case "assert-visible", "assert-hidden":
		return actionParameters("assert_visible")
	case "assert-text":
		return actionParameters("assert_text")
	case "assert-url", "assert-url-contains":
		return actionParameters("assert_url")
	case "wait-visible":
		return actionParameters("wait_for")
	case "wait-hidden":
		return actionParameters("wait_for_hidden")
	case "wait":
		return actionParameters("wait")
	case "remember-text":
		return actionParameters("remember_text")
	case "remember-field":
		return actionParameters("remember_field")
	case "remember-url":
		return actionParameters("remember_url")
	case "switch-tab":
		return actionParameters("switch_tab")
	case "close-tab":
		return actionParameters("close_tab")
	case "assert-tab-count":
		return actionParameters("assert_tab_count")
	case "download-click":
		return actionParameters("download_click")
	case "assert-download-contains":
		return actionParameters("assert_download_contains")
	case "draw-signature":
		return actionParameters("draw_signature")
	case "press", "press-in":
		return actionParameters("press")
	case "prompt-email-code":
		return actionParameters("prompt_email_code")
	case "scroll-to":
		return actionParameters("scroll_to")
	case "reload":
		return actionParameters("reload")
	case "go-back":
		return actionParameters("go_back")
	case "close-browser":
		return actionParameters("close_browser")
	default:
		return nil
	}
}

func categoryForKind(kind string) string {
	switch kind {
	case "goto", "reload", "go-back", "close-browser":
		return "Навигация"
	case "click", "double-click", "hover", "fill", "select", "upload", "clear", "check", "uncheck", "press", "press-in", "download-click", "draw-signature", "fill-generated", "prompt-email-code":
		return "Формы и ввод"
	case "assert-visible", "assert-hidden", "assert-text", "assert-url", "assert-url-contains", "assert-tab-count", "assert-download-contains":
		return "Проверки"
	case "wait-visible", "wait-hidden", "wait":
		return "Ожидание"
	case "remember-text", "remember-field", "remember-url":
		return "Переменные"
	case "switch-tab", "close-tab", "scroll-to":
		return "Вкладки"
	default:
		return "Шаги"
	}
}

func helpForKind(kind string) string {
	switch kind {
	case "goto":
		return "Переход на страницу"
	case "click":
		return "Клик по элементу"
	case "double-click":
		return "Двойной клик"
	case "hover":
		return "Навести курсор"
	case "fill":
		return "Ввод текста в поле"
	case "select":
		return "Выбор значения в списке"
	case "upload":
		return "Загрузка файла в поле"
	case "clear":
		return "Очистка поля ввода"
	case "check", "uncheck":
		return "Чекбокс"
	case "assert-visible", "assert-hidden":
		return "Проверка видимости элемента"
	case "assert-text":
		return "Проверка текста в элементе"
	case "assert-url", "assert-url-contains":
		return "Проверка текущего URL"
	case "wait-visible", "wait-hidden":
		return "Ожидание появления или скрытия элемента"
	case "wait":
		return "Пауза перед следующим шагом"
	case "reload":
		return "Перезагрузка страницы"
	case "go-back":
		return "Кнопка «Назад» браузера"
	case "close-browser":
		return "Закрыть окно браузера"
	case "remember-text", "remember-field", "remember-url":
		return "Сохранить значение в переменную"
	case "switch-tab":
		return "Переключение вкладки"
	case "close-tab":
		return "Закрыть активную вкладку"
	case "assert-tab-count":
		return "Проверить число открытых вкладок"
	case "scroll-to":
		return "Прокрутка к элементу"
	case "download-click", "assert-download-contains":
		return "Скачивание и проверка файла"
	case "fill-generated":
		return "Случайное значение в поле"
	case "draw-signature":
		return "Рисование подписи на canvas"
	case "press", "press-in":
		return "Нажатие клавиши"
	case "prompt-email-code":
		return "Код из почты (OTP)"
	case "drag-drop":
		return "Перетаскивание элемента"
	default:
		return kind
	}
}
