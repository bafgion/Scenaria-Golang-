package stepdsl

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bafgion/scenaria-golang/internal/gherkin"
)

type Action struct {
	Kind   string
	Value1 string
	Value2 string
	Value3 string
	Mode   string
	IntVal int
}

var quoted = `((?:\\.|[^"])*)`

var stepPatterns = []struct {
	re   *regexp.Regexp
	kind string
	mapFn func(groups []string) Action
}{
	{regexp.MustCompile(`(?i)^запоминаю текст "` + quoted + `" как "` + quoted + `"$`), "remember-text", two("remember-text")},
	{regexp.MustCompile(`(?i)^запоминаю значение поля "` + quoted + `" как "` + quoted + `"$`), "remember-field", two("remember-field")},
	{regexp.MustCompile(`(?i)^запоминаю url как "` + quoted + `"$`), "remember-url", one("remember-url")},
	{regexp.MustCompile(`(?i)^скачиваю по клику на "` + quoted + `"$`), "download-click", one("download-click")},
	{regexp.MustCompile(`(?i)^проверяю что скачанный файл содержит "` + quoted + `"$`), "assert-download-contains", one("assert-download-contains")},
	{regexp.MustCompile(`(?i)^открыт[а]?\s+"` + quoted + `"$`), "goto", one("goto")},
	{regexp.MustCompile(`(?i)^(?:открываю|перехожу на|перейти на)\s+"` + quoted + `"$`), "goto", one("goto")},
	{regexp.MustCompile(`(?i)^дважды нажимаю "` + quoted + `"$`), "double-click", one("double-click")},
	{regexp.MustCompile(`(?i)^нажимаю клавишу "` + quoted + `" в "` + quoted + `"$`), "press-in", two("press-in")},
	{regexp.MustCompile(`(?i)^нажимаю клавишу "` + quoted + `"$`), "press", one("press")},
	{regexp.MustCompile(`(?i)^(?:нажимаю|кликаю) "` + quoted + `"$`), "click", one("click")},
	{regexp.MustCompile(`(?i)^навожу "` + quoted + `"$`), "hover", one("hover")},
	{regexp.MustCompile(`(?i)^ввожу код из почты с клавиатуры "` + quoted + `" в (\d+) пол(?:ей|я) "` + quoted + `"$`), "prompt-email-code", emailCodeFields("keyboard")},
	{regexp.MustCompile(`(?i)^ввожу код из почты заполнением "` + quoted + `" в (\d+) пол(?:ей|я) "` + quoted + `"$`), "prompt-email-code", emailCodeFields("fill")},
	{regexp.MustCompile(`(?i)^ввожу код из почты с клавиатуры в (\d+) пол(?:ей|я) "` + quoted + `"$`), "prompt-email-code", emailCodeFieldsNoEmail("keyboard")},
	{regexp.MustCompile(`(?i)^ввожу код из почты заполнением в (\d+) пол(?:ей|я) "` + quoted + `"$`), "prompt-email-code", emailCodeFieldsNoEmail("fill")},
	{regexp.MustCompile(`(?i)^ввожу код из почты с клавиатуры "` + quoted + `" в "` + quoted + `"$`), "prompt-email-code", emailCodeSingle("keyboard")},
	{regexp.MustCompile(`(?i)^ввожу код из почты заполнением "` + quoted + `" в "` + quoted + `"$`), "prompt-email-code", emailCodeSingle("fill")},
	{regexp.MustCompile(`(?i)^ввожу код из почты с клавиатуры в "` + quoted + `"$`), "prompt-email-code", emailCodeSingleNoEmail("keyboard")},
	{regexp.MustCompile(`(?i)^ввожу код из почты заполнением в "` + quoted + `"$`), "prompt-email-code", emailCodeSingleNoEmail("fill")},
	{regexp.MustCompile(`(?i)^ввожу код из почты "` + quoted + `" в "` + quoted + `"$`), "prompt-email-code", emailCodeSingle("fill")},
	{regexp.MustCompile(`(?i)^ввожу код из почты в "` + quoted + `"$`), "prompt-email-code", emailCodeSingleNoEmail("fill")},
	{regexp.MustCompile(`(?i)^ввожу случайный телефон в "` + quoted + `"$`), "fill-generated", fillGen("phone")},
	{regexp.MustCompile(`(?i)^ввожу случайное имя в "` + quoted + `"$`), "fill-generated", fillGen("first_name")},
	{regexp.MustCompile(`(?i)^ввожу случайную фамилию в "` + quoted + `"$`), "fill-generated", fillGen("last_name")},
	{regexp.MustCompile(`(?i)^ввожу случайное отчество в "` + quoted + `"$`), "fill-generated", fillGen("patronymic")},
	{regexp.MustCompile(`(?i)^ввожу случайный адрес в "` + quoted + `"$`), "fill-generated", fillGen("address")},
	{regexp.MustCompile(`(?i)^ввожу случайный инн в "` + quoted + `"$`), "fill-generated", fillGen("inn")},
	{regexp.MustCompile(`(?i)^ввожу случайный расчётный счёт в "` + quoted + `"$`), "fill-generated", fillGen("bank_account")},
	{regexp.MustCompile(`(?i)^ввожу случайный расчетный счет в "` + quoted + `"$`), "fill-generated", fillGen("bank_account")},
	{regexp.MustCompile(`(?i)^ввожу случайный огрнип в "` + quoted + `"$`), "fill-generated", fillGen("ogrnip")},
	{regexp.MustCompile(`(?i)^ввожу случайный (.+?) в "` + quoted + `"$`), "fill-generated", fillGenAlias},
	{regexp.MustCompile(`(?i)^ввожу "` + quoted + `" в "` + quoted + `"$`), "fill", two("fill")},
	{regexp.MustCompile(`(?i)^выбираю "` + quoted + `" в "` + quoted + `"$`), "select", two("select")},
	{regexp.MustCompile(`(?i)^загружаю файл "` + quoted + `" в "` + quoted + `"$`), "upload", two("upload")},
	{regexp.MustCompile(`(?i)^очищаю "` + quoted + `"$`), "clear", one("clear")},
	{regexp.MustCompile(`(?i)^рисую подпись в "` + quoted + `"$`), "draw-signature", one("draw-signature")},
	{regexp.MustCompile(`(?i)^снимаю отметку с "` + quoted + `"$`), "uncheck", one("uncheck")},
	{regexp.MustCompile(`(?i)^отмечаю "` + quoted + `"$`), "check", one("check")},
	{regexp.MustCompile(`(?i)^не вижу "` + quoted + `"$`), "assert-hidden", one("assert-hidden")},
	{regexp.MustCompile(`(?i)^вижу "` + quoted + `"$`), "assert-visible", one("assert-visible")},
	{regexp.MustCompile(`(?i)^проверяю текст "` + quoted + `" в "` + quoted + `"$`), "assert-text", two("assert-text")},
	{regexp.MustCompile(`(?i)^проверяю url "` + quoted + `"$`), "assert-url", one("assert-url")},
	{regexp.MustCompile(`(?i)^(?:url содержит|адрес содержит) "` + quoted + `"$`), "assert-url-contains", one("assert-url-contains")},
	{regexp.MustCompile(`(?i)^скроллю к "` + quoted + `"$`), "scroll-to", one("scroll-to")},
	{regexp.MustCompile(`(?i)^обновляю страницу$`), "reload", none("reload")},
	{regexp.MustCompile(`(?i)^возвращаюсь назад$`), "go-back", none("go-back")},
	{regexp.MustCompile(`(?i)^закрываю браузер$`), "close-browser", none("close-browser")},
	{regexp.MustCompile(`(?i)^переключаюсь на вкладку (\d+)$`), "switch-tab", switchTabIndex},
	{regexp.MustCompile(`(?i)^переключаюсь на вкладку "` + quoted + `"$`), "switch-tab", switchTabTitle},
	{regexp.MustCompile(`(?i)^переключаюсь на вкладку с url "` + quoted + `"$`), "switch-tab", switchTabURL},
	{regexp.MustCompile(`(?i)^переключаюсь на первую вкладку$`), "switch-tab", switchTabFirst},
	{regexp.MustCompile(`(?i)^переключаюсь на новую вкладку$`), "switch-tab", switchTabNew},
	{regexp.MustCompile(`(?i)^закрываю текущую вкладку$`), "close-tab", none("close-tab")},
	{regexp.MustCompile(`(?i)^проверяю что открыто (\d+) вкладк[аиуы]?$`), "assert-tab-count", assertTabCount},
	{regexp.MustCompile(`(?i)^жду исчезновения "` + quoted + `"$`), "wait-hidden", one("wait-hidden")},
	{regexp.MustCompile(`(?i)^жду появления "` + quoted + `"$`), "wait-visible", one("wait-visible")},
	{regexp.MustCompile(`(?i)^жду (\d+)\s*(?:мс|мсек)\s*$`), "wait", waitMillis},
	{regexp.MustCompile(`(?i)^жду ([\d.,]+)\s*с\s*$`), "wait", waitSeconds},
	{regexp.MustCompile(`(?i)^жду ([\d.,]+)\s*(?:сек|секунд|сек)\s*$`), "wait", waitSeconds},
	{regexp.MustCompile(`(?i)^(?:жду|ожидаю) "` + quoted + `"$`), "wait", waitQuoted},
}

func one(kind string) func([]string) Action {
	return func(groups []string) Action {
		return Action{Kind: kind, Value1: unquote(groups[1])}
	}
}

func two(kind string) func([]string) Action {
	return func(groups []string) Action {
		return Action{Kind: kind, Value1: unquote(groups[1]), Value2: unquote(groups[2])}
	}
}

func none(kind string) func([]string) Action {
	return func([]string) Action {
		return Action{Kind: kind}
	}
}

func fillGen(generator string) func([]string) Action {
	return func(groups []string) Action {
		return Action{Kind: "fill-generated", Value1: generator, Value2: unquote(groups[1])}
	}
}

func fillGenAlias(groups []string) Action {
	return Action{Kind: "fill-generated", Value1: strings.TrimSpace(unquote(groups[1])), Value2: unquote(groups[2])}
}

func emailCodeFields(method string) func([]string) Action {
	return func(groups []string) Action {
		digits, _ := strconv.Atoi(groups[2])
		return Action{
			Kind:   "prompt-email-code",
			Value1: unquote(groups[1]),
			Value2: unquote(groups[3]),
			Value3: method,
			IntVal: digits,
		}
	}
}

func emailCodeFieldsNoEmail(method string) func([]string) Action {
	return func(groups []string) Action {
		digits, _ := strconv.Atoi(groups[1])
		return Action{
			Kind:   "prompt-email-code",
			Value2: unquote(groups[2]),
			Value3: method,
			IntVal: digits,
		}
	}
}

func emailCodeSingle(method string) func([]string) Action {
	return func(groups []string) Action {
		return Action{
			Kind:   "prompt-email-code",
			Value1: unquote(groups[1]),
			Value2: unquote(groups[2]),
			Value3: method,
			IntVal: 1,
		}
	}
}

func emailCodeSingleNoEmail(method string) func([]string) Action {
	return func(groups []string) Action {
		return Action{
			Kind:   "prompt-email-code",
			Value2: unquote(groups[1]),
			Value3: method,
			IntVal: 1,
		}
	}
}

func switchTabIndex(groups []string) Action {
	userIndex, _ := strconv.Atoi(groups[1])
	return Action{Kind: "switch-tab", Mode: "index", IntVal: userIndex - 1}
}

func switchTabTitle(groups []string) Action {
	return Action{Kind: "switch-tab", Mode: "title", Value1: unquote(groups[1])}
}

func switchTabURL(groups []string) Action {
	return Action{Kind: "switch-tab", Mode: "url", Value1: unquote(groups[1])}
}

func switchTabFirst([]string) Action {
	return Action{Kind: "switch-tab", Mode: "first"}
}

func switchTabNew([]string) Action {
	return Action{Kind: "switch-tab", Mode: "new"}
}

func assertTabCount(groups []string) Action {
	count, _ := strconv.Atoi(groups[1])
	return Action{Kind: "assert-tab-count", IntVal: count}
}

func waitSeconds(groups []string) Action {
	raw := strings.ReplaceAll(groups[1], ",", ".")
	seconds, _ := strconv.ParseFloat(raw, 64)
	if seconds < 0 {
		seconds = 0
	}
	return Action{Kind: "wait", Value1: fmt.Sprintf("%dms", int(seconds*1000))}
}

func waitMillis(groups []string) Action {
	ms, _ := strconv.Atoi(groups[1])
	if ms < 0 {
		ms = 0
	}
	return Action{Kind: "wait", Value1: fmt.Sprintf("%dms", ms)}
}

func waitQuoted(groups []string) Action {
	return Action{Kind: "wait", Value1: unquote(groups[1])}
}

func Parse(step gherkin.Step) (Action, error) {
	body := strings.TrimSpace(step.Text)
	for _, pattern := range stepPatterns {
		if groups := pattern.re.FindStringSubmatch(body); groups != nil {
			return pattern.mapFn(groups), nil
		}
	}
	return Action{}, fmt.Errorf("line %d: unsupported step text %q", step.Line, step.Text)
}

func ResolveURL(raw string, baseURL string) string {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return trimmed
	}
	if strings.TrimSpace(baseURL) == "" {
		return trimmed
	}
	base := strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if strings.HasPrefix(trimmed, "/") {
		return base + trimmed
	}
	return base + "/" + strings.TrimLeft(trimmed, "/")
}

func ParseWaitDuration(raw string) (time.Duration, error) {
	trimmed := strings.TrimSpace(raw)
	if strings.HasSuffix(trimmed, "ms") {
		ms, err := strconv.Atoi(strings.TrimSuffix(trimmed, "ms"))
		if err != nil {
			return 0, fmt.Errorf("invalid wait duration %q: %w", raw, err)
		}
		return time.Duration(ms) * time.Millisecond, nil
	}
	return time.ParseDuration(trimmed)
}

func unquote(s string) string {
	return strings.ReplaceAll(s, `\"`, `"`)
}
