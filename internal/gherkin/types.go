package gherkin

type Feature struct {
	Title           string
	Line            int
	Tags            []string
	TestClient      string
	HasContextBlock bool
	Background      []Step
	Scenarios       []Scenario
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

type Condition struct {
	Type     string
	Selector string
	Value    string
}

type Step struct {
	Keyword         string
	Text            string
	Line            int
	Indent          int
	DocString       string
	Table           [][]string
	Block           string
	Condition       *Condition
	RepeatCount     int
	ForEachSelector string
	ForEachVariable string
	Children        []Step
}

type Example struct {
	Line int
	Rows [][]string
}

type Issue struct {
	Line    int
	Message string
}

const (
	BlockIf      = "if"
	BlockRepeat  = "repeat"
	BlockWhile   = "while"
	BlockForEach = "for_each"
)
