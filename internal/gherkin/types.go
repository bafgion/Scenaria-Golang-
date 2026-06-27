package gherkin

type Feature struct {
	Title      string
	Background []Step
	Scenarios  []Scenario
}

type Scenario struct {
	Title string
	Steps []Step
}

type Step struct {
	Keyword string
	Text    string
	Line    int
}

type Issue struct {
	Line    int
	Message string
}
