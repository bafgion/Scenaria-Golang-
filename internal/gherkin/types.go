package gherkin

type Feature struct {
	Title      string
	Line       int
	Tags       []string
	Background []Step
	Scenarios  []Scenario
}

type Scenario struct {
	Title        string
	Line         int
	Tags         []string
	IsOutline    bool
	Steps        []Step
	Examples     []Example
	ExamplesLine int
}

type Step struct {
	Keyword   string
	Text      string
	Line      int
	DocString string
	Table     [][]string
}

type Example struct {
	Line int
	Rows [][]string
}

type Issue struct {
	Line    int
	Message string
}
