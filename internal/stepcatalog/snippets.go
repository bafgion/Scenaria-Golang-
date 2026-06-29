package stepcatalog

import (
	"regexp"
	"strings"
)

type snippetDef struct {
	label       string
	insert      string
	description string
}

var actionRE = regexp.MustCompile(`(?i)action:\s*([\w-]+)`)

func buildSnippetEntries() []Entry {
	out := make([]Entry, 0, len(stepSnippets))
	for _, snip := range stepSnippets {
		action := actionFromDescription(snip.description)
		out = append(out, Entry{
			Label:       snip.label,
			Action:      action,
			Category:    categoryForAction(action),
			Description: plainDescription(snip.description),
			Template:    snip.insert,
			Example:     snip.insert,
			Parameters:  actionParameters(action),
			Help:        plainDescription(snip.description),
		})
	}
	return out
}

func actionFromDescription(description string) string {
	match := actionRE.FindStringSubmatch(description)
	if len(match) < 2 {
		return ""
	}
	return match[1]
}

func plainDescription(description string) string {
	text := strings.TrimSpace(actionRE.ReplaceAllString(description, ""))
	return strings.Trim(text, " .")
}

func categoryForAction(action string) string {
	switch action {
	case "goto", "go_back", "reload", "scroll_to", "remember_url":
		return "Навигация"
	case "click", "double_click", "hover", "fill", "clear", "select", "check", "uncheck", "press", "remember_text", "remember_field", "prompt_email_code":
		return "Формы и ввод"
	case "fill_generated":
		return "Генераторы"
	case "upload":
		return "Файлы"
	case "download_click", "assert_download_contains":
		return "Файлы"
	case "assert_visible", "assert_hidden", "assert_text", "assert_url", "assert_tab_count":
		return "Проверки"
	case "wait", "wait_for", "wait_for_hidden":
		return "Ожидание"
	case "close_browser":
		return "Сессия"
	case "switch_tab", "close_tab":
		return "Вкладки"
	case "draw_signature":
		return "Формы и ввод"
	case "if":
		return "Условия"
	case "repeat", "while", "for_each":
		return "Циклы"
	default:
		return "Шаги"
	}
}

func actionParameters(action string) []string {
	if params, ok := actionParams[action]; ok {
		out := make([]string, len(params))
		copy(out, params)
		return out
	}
	return nil
}

var actionParams = map[string][]string{
	"goto":                      {`url — адрес страницы в кавычках`},
	"click":                     {`selector — CSS/XPath селектор элемента`},
	"double_click":              {`selector — элемент для двойного клика`},
	"hover":                     {`selector — элемент для наведения`},
	"fill":                      {`value — текст; selector — поле ввода`},
	"fill_generated":            {`generator — тип данных; selector — поле ввода`},
	"clear":                     {`selector — поле ввода`},
	"select":                    {`value — значение option; selector — элемент select`},
	"check":                     {`selector — checkbox или radio`},
	"uncheck":                   {`selector — checkbox`},
	"press":                     {`key — имя клавиши (Enter, Tab…); selector — опционально`},
	"prompt_email_code":         {`digits/selector — ячейки OTP; email — опционально`},
	"upload":                    {`file_path — путь к файлу; selector — input[type=file]`},
	"download_click":            {`selector — ссылка или кнопка скачивания`},
	"assert_download_contains":  {`text — подстрока в скачанном файле`},
	"remember_text":             {`value — текст; variable — имя переменной`},
	"remember_field":            {`selector — поле; variable — имя переменной`},
	"remember_url":              {`variable — имя для текущего URL`},
	"draw_signature":            {`selector — элемент canvas`},
	"scroll_to":                 {`selector — элемент, к которому прокрутить`},
	"assert_visible":            {`selector — видимый элемент`},
	"assert_hidden":             {`selector — скрытый элемент`},
	"assert_text":               {`text — ожидаемый текст; selector — контейнер`},
	"assert_url":                {`url — ожидаемый адрес страницы`},
	"wait":                      {`seconds — длительность паузы`},
	"wait_for":                  {`selector — элемент, появление которого ждём`},
	"wait_for_hidden":           {`selector — элемент, исчезновение которого ждём`},
	"switch_tab":                {`mode — title, url, index (номер с 1), first, new`, `value — фрагмент заголовка, URL или индекс (для title/url/index)`},
	"assert_tab_count":          {`count — ожидаемое число вкладок`},
	"if":                        {`condition — вижу / не вижу / url содержит / текст на странице`, `steps — вложенные шаги с отступом (таб или 2 пробела)`},
	"repeat":                    {`count — число повторений`, `steps — тело цикла с отступом`},
	"while":                     {`condition — как в «Если»`, `steps — тело цикла; лимит итераций — max_loop_iterations в настройках`},
	"for_each":                  {`selector — CSS селектор элементов списка`, `variable — имя переменной для {{имя}} во вложенных шагах`, `steps — тело цикла с отступом`},
}

// stepSnippets mirrors app/gherkin_snippets.py STEP_SNIPPETS (Python Scenaria).
var stepSnippets = []snippetDef{
	{"открыт", `открыт "https://site.com"`, "Переход на страницу (action: goto)"},
	{"нажимаю", `нажимаю "button.submit"`, "Клик по элементу (action: click)"},
	{"дважды нажимаю", `дважды нажимаю ".file-item"`, "Двойной клик (action: double_click)"},
	{"навожу", `навожу "nav a:has-text('Услуги')"`, "Наведение перед кликом по подменю (action: hover)"},
	{"ввожу", `ввожу "текст" в "input[name=email]"`, "Ввод текста в поле (action: fill)"},
	{"случайный телефон", `ввожу случайный телефон в "input[type=tel]"`, "Случайный телефон на каждый прогон (action: fill_generated)"},
	{"случайное имя", `ввожу случайное имя в "input[name=firstName]"`, "Случайное имя на каждый прогон (action: fill_generated)"},
	{"случайная фамилия", `ввожу случайную фамилию в "input[name=lastName]"`, "Случайная фамилия на каждый прогон (action: fill_generated)"},
	{"случайное отчество", `ввожу случайное отчество в "input[name=middleName]"`, "Случайное отчество на каждый прогон (action: fill_generated)"},
	{"случайный адрес", `ввожу случайный адрес в "textarea[name=address]"`, "Случайный адрес на каждый прогон (action: fill_generated)"},
	{"случайный инн", `ввожу случайный инн в "input[name=inn]"`, "ИНН физлица (12 цифр) с контрольной суммой (action: fill_generated)"},
	{"расчётный счёт", `ввожу случайный расчётный счёт в "input[name=account]"`, "20-значный расчётный счёт (action: fill_generated)"},
	{"огрнип", `ввожу случайный огрнип в "input[name=ogrnip]"`, "ОГРНИП (15 цифр) с контрольной суммой (action: fill_generated)"},
	{"плейсхолдер", `ввожу "{{first_name}}" в "input[name=firstName]"`, "Подстановка переменных {{first_name}} и др. (action: fill)"},
	{"ввожу код из почты", `ввожу код из почты с клавиатуры в 6 полей "input.pin-digit"`, "OTP через keyboard.type в нескольких полях (action: prompt_email_code)"},
	{"код из почты email", `ввожу код из почты "user@example.com" в "input#code"`, "OTP из письма на указанный ящик (action: prompt_email_code)"},
	{"очищаю", `очищаю "input#search"`, "Очистка поля ввода (action: clear)"},
	{"выбираю", `выбираю "Значение" в "select#country"`, "Выбор значения в списке (action: select)"},
	{"отмечаю", `отмечаю "input#agree"`, "Установка галочки (action: check)"},
	{"снимаю отметку", `снимаю отметку с "input#newsletter"`, "Снятие галочки (action: uncheck)"},
	{"нажимаю клавишу", `нажимаю клавишу "Enter"`, "Клавиша на странице — Enter, Escape, Tab… (action: press)"},
	{"нажимаю клавишу в", `нажимаю клавишу "Tab" в "input[name=email]"`, "Клавиша в конкретном поле (action: press)"},
	{"загружаю файл", `загружаю файл "C:\\data\\doc.pdf" в "input[type=file]"`, "Загрузка файла в поле (action: upload)"},
	{"скачиваю", `скачиваю по клику на "a.export"`, "Скачивание файла по клику (action: download_click)"},
	{"скачанный файл", `проверяю что скачанный файл содержит "Invoice"`, "Проверка содержимого скачанного файла (action: assert_download_contains)"},
	{"запоминаю текст", `запоминаю текст "{{login}}" как "user_login"`, "Сохранить литерал в переменную (action: remember_text)"},
	{"запоминаю url", `запоминаю url как "current_url"`, "Сохранить текущий URL (action: remember_url)"},
	{"запоминаю значение поля", `запоминаю значение поля "input#email" как "user_email"`, "Сохранить значение поля в переменную (action: remember_field)"},
	{"рисую подпись", `рисую подпись в "canvas"`, "Рисование подписи на canvas, ПЭП (action: draw_signature)"},
	{"скроллю к", `скроллю к "section#contacts"`, "Прокрутка к элементу (action: scroll_to)"},
	{"обновляю страницу", "обновляю страницу", "Перезагрузка страницы (action: reload)"},
	{"возвращаюсь назад", "возвращаюсь назад", "Кнопка «Назад» браузера (action: go_back)"},
	{"закрываю браузер", "закрываю браузер", "Закрыть окно браузера (action: close_browser)"},
	{"переключаюсь на вкладку", `переключаюсь на вкладку "Оплата"`, "Активировать вкладку по заголовку (часть title) (action: switch_tab)"},
	{"переключаюсь на вкладку с url", `переключаюсь на вкладку с url "checkout"`, "Активировать вкладку по подстроке в URL (action: switch_tab)"},
	{"переключаюсь на вкладку 2", "переключаюсь на вкладку 2", "Переключиться на вкладку по номеру, с 1 (action: switch_tab)"},
	{"переключаюсь на первую вкладку", "переключаюсь на первую вкладку", "Переключиться на первую вкладку в окне (action: switch_tab)"},
	{"переключаюсь на новую вкладку", "переключаюсь на новую вкладку", "Переключиться на последнюю открытую вкладку (action: switch_tab)"},
	{"закрываю текущую вкладку", "закрываю текущую вкладку", "Закрыть активную вкладку (action: close_tab)"},
	{"проверяю что открыто", "проверяю что открыто 2 вкладки", "Проверить число открытых вкладок (action: assert_tab_count)"},
	{"вижу", `вижу "h1.title"`, "Проверка видимости элемента (action: assert_visible)"},
	{"не вижу", `не вижу ".modal-overlay"`, "Элемент скрыт или отсутствует (action: assert_hidden)"},
	{"проверяю текст", `проверяю текст "Успех" в ".message"`, "Проверка текста в элементе (action: assert_text)"},
	{"проверяю url", `проверяю url "https://site.com/profile"`, "Проверка текущего URL (action: assert_url)"},
	{"жду", "жду 2 сек", "Пауза перед следующим шагом (action: wait)"},
	{"жду мс", "жду 500 мс", "Короткая пауза в миллисекундах (action: wait)"},
	{"жду появления", `жду появления "button.ready"`, "Ожидание появления элемента (action: wait_for)"},
	{"жду исчезновения", `жду исчезновения ".spinner"`, "Ожидание скрытия лоадера или модалки (action: wait_for_hidden)"},
	{"если вижу", "Если вижу \".cookie-banner\"\n\tнажимаю \"button.accept\"", "Условный блок: шаги внутри выполняются, если элемент виден (action: if)"},
	{"если не вижу", "Если не вижу \".spinner\"\n\tнажимаю \"button.next\"", "Условный блок: шаги внутри выполняются, если элемент скрыт (action: if)"},
	{"если url", "Если url содержит \"/shop\"\n\tнажимаю \"button.next\"", "Условный блок: шаги внутри выполняются, если URL содержит подстроку (action: if)"},
	{"если текст", "Если текст на странице \"Готово\"\n\tнажимаю \"button.ok\"", "Условный блок: шаги внутри выполняются, если текст есть на странице (action: if)"},
	{"повторяю", "Повторяю 3 раза\n\tнажимаю \"button.add\"", "Цикл с фиксированным числом повторений (action: repeat)"},
	{"пока", "Пока вижу \"button.load-more\"\n\tнажимаю \"button.load-more\"", "Цикл «пока условие истинно»; лимит итераций — в настройках (action: while)"},
	{"для каждого", "Для каждого \".product-card\" как \"card\"\n\tнажимаю \"{{card}} .buy-btn\"", "Цикл по каждому элементу списка; переменная доступна как {{card}} (action: for_each)"},
}
