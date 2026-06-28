package stepcatalog

import "strings"

type Entry struct {
	Category string
	Template string
	Help     string
}

func Entries() []Entry {
	return []Entry{
		{Category: "Навигация", Template: `открыт "https://example.com"`, Help: "Открыть URL"},
		{Category: "Навигация", Template: "обновляю страницу", Help: "Перезагрузить страницу"},
		{Category: "Действия", Template: `нажимаю "#submit"`, Help: "Клик по элементу"},
		{Category: "Действия", Template: `ввожу "value" в "#field"`, Help: "Ввод текста"},
		{Category: "Проверки", Template: `вижу "#result"`, Help: "Элемент виден"},
		{Category: "Проверки", Template: `проверяю текст "ok" в "#result"`, Help: "Текст в элементе"},
		{Category: "Условия", Template: `Если вижу "#modal"`, Help: "Условный блок"},
		{Category: "Циклы", Template: "Повторяю 3 раза", Help: "Повтор шагов"},
		{Category: "Переменные", Template: `запоминаю url как "current_url"`, Help: "Сохранить URL в переменную"},
		{Category: "Вкладки", Template: "переключаюсь на вкладку 2", Help: "Переключение по индексу (с 1)"},
		{Category: "Вкладки", Template: `переключаюсь на вкладку с url "dashboard"`, Help: "Переключение по URL"},
		{Category: "Генераторы", Template: `ввожу случайный телефон в "#phone"`, Help: "Случайное значение"},
		{Category: "Файлы", Template: `скачиваю по клику на "#export"`, Help: "Скачивание по клику"},
		{Category: "Сессия", Template: "закрываю браузер", Help: "Закрыть браузер"},
	}
}

func Search(query string) []Entry {
	query = stringsToLower(query)
	if query == "" {
		return Entries()
	}
	out := make([]Entry, 0)
	for _, entry := range Entries() {
		if containsFold(entry.Template, query) || containsFold(entry.Help, query) || containsFold(entry.Category, query) {
			out = append(out, entry)
		}
	}
	return out
}

func stringsToLower(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func containsFold(haystack, needle string) bool {
	return strings.Contains(stringsToLower(haystack), needle)
}
