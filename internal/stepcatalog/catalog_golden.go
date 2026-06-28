package stepcatalog

import (
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
	out := make([]Entry, 0, 64)

	add := func(category, template, help string) {
		if template == "" || seen[template] {
			return
		}
		seen[template] = true
		out = append(out, Entry{Category: category, Template: template, Help: help})
	}

	cases, err := stepdsl.LoadGoldenCases()
	if err == nil {
		for _, c := range cases {
			add(categoryForKind(c.Kind), c.Text, helpForKind(c.Kind))
		}
	}

	for _, entry := range manualEntries() {
		add(entry.Category, entry.Template, entry.Help)
	}
	return out
}

func manualEntries() []Entry {
	return []Entry{
		{Category: "Структура", Template: "Если вижу \"#modal\"", Help: "Условный блок"},
		{Category: "Структура", Template: "Повторяю 3 раза", Help: "Цикл с числом итераций"},
		{Category: "Структура", Template: "Пока вижу \"#next\"", Help: "Цикл пока элемент виден"},
		{Category: "Структура", Template: "Для каждого \"#row\" как \"item\"", Help: "Цикл по элементам"},
		{Category: "Ожидание", Template: "жду 2 с", Help: "Пауза в секундах"},
		{Category: "Генераторы", Template: "ввожу случайный огрнип в \"#ogrnip\"", Help: "Случайный ОГРНИП"},
		{Category: "Генераторы", Template: "ввожу случайное фио в \"#name\"", Help: "Случайное ФИО"},
	}
}

func categoryForKind(kind string) string {
	switch kind {
	case "goto", "reload", "go-back", "close-browser":
		return "Навигация"
	case "click", "double-click", "hover", "fill", "select", "upload", "clear", "check", "uncheck", "press", "press-in", "download-click", "draw-signature", "fill-generated", "prompt-email-code":
		return "Действия"
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
		return "Открыть URL"
	case "click":
		return "Клик по элементу"
	case "double-click":
		return "Двойной клик"
	case "hover":
		return "Навести курсор"
	case "fill":
		return "Ввод текста"
	case "select":
		return "Выбор в списке"
	case "upload":
		return "Загрузка файла"
	case "clear":
		return "Очистить поле"
	case "check", "uncheck":
		return "Чекбокс"
	case "assert-visible", "assert-hidden":
		return "Видимость элемента"
	case "assert-text":
		return "Текст в элементе"
	case "assert-url", "assert-url-contains":
		return "Проверка URL"
	case "wait-visible", "wait-hidden":
		return "Ожидание элемента"
	case "wait":
		return "Пауза"
	case "reload":
		return "Перезагрузить страницу"
	case "go-back":
		return "Назад в истории"
	case "close-browser":
		return "Закрыть браузер"
	case "remember-text", "remember-field", "remember-url":
		return "Сохранить в переменную"
	case "switch-tab":
		return "Переключение вкладки"
	case "close-tab":
		return "Закрыть вкладку"
	case "assert-tab-count":
		return "Количество вкладок"
	case "scroll-to":
		return "Прокрутка к элементу"
	case "download-click", "assert-download-contains":
		return "Скачивание файла"
	case "fill-generated":
		return "Случайное значение"
	case "draw-signature":
		return "Рисование подписи на canvas"
	case "press", "press-in":
		return "Нажатие клавиши"
	case "prompt-email-code":
		return "Код из почты"
	default:
		return kind
	}
}
