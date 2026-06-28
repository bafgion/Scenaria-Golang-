//go:build desktop

package desktop

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/bafgion/scenaria-golang/internal/featurehighlight"
)

type syntaxEditor struct {
	widget.BaseWidget
	entry     *widget.Entry
	rich      *widget.RichText
	onChanged func(string)
}

func newSyntaxEditor() *syntaxEditor {
	editor := &syntaxEditor{}
	editor.ExtendBaseWidget(editor)
	editor.entry = widget.NewMultiLineEntry()
	editor.entry.TextStyle = fyne.TextStyle{Monospace: true}
	editor.rich = widget.NewRichText()
	editor.rich.Wrapping = fyne.TextWrapOff
	editor.entry.OnChanged = editor.syncText
	return editor
}

func (e *syntaxEditor) syncText(text string) {
	e.applyHighlight(text)
	if e.onChanged != nil {
		e.onChanged(text)
	}
}

func (e *syntaxEditor) applyHighlight(text string) {
	spans := featurehighlight.Highlight(text)
	segments := make([]widget.RichTextSegment, 0, len(spans))
	for _, span := range spans {
		segments = append(segments, &widget.TextSegment{
			Text:  span.Text,
			Style: richStyleFor(span.Kind),
		})
	}
	e.rich.Segments = segments
	e.rich.Refresh()
}

func (e *syntaxEditor) SetOnChanged(fn func(string)) {
	e.onChanged = fn
}

func (e *syntaxEditor) SetText(text string) {
	e.entry.SetText(text)
	e.applyHighlight(text)
}

func (e *syntaxEditor) Text() string {
	return e.entry.Text
}

func (e *syntaxEditor) CreateRenderer() fyne.WidgetRenderer {
	preview := container.NewScroll(e.rich)
	preview.SetMinSize(fyne.NewSize(120, 0))
	split := container.NewHSplit(e.entry, preview)
	split.SetOffset(0.64)
	return widget.NewSimpleRenderer(split)
}

func richStyleFor(kind featurehighlight.Kind) widget.RichTextStyle {
	style := widget.RichTextStyle{TextStyle: fyne.TextStyle{Monospace: true}}
	switch kind {
	case featurehighlight.KindComment:
		style.ColorName = theme.ColorNameDisabled
	case featurehighlight.KindTag:
		style.ColorName = theme.ColorNameHyperlink
	case featurehighlight.KindGherkinKeyword, featurehighlight.KindStepKeyword, featurehighlight.KindBlockKeyword:
		style.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		style.ColorName = theme.ColorNamePrimary
	case featurehighlight.KindString:
		style.ColorName = theme.ColorNameSuccess
	case featurehighlight.KindTestClient:
		style.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}
		style.ColorName = theme.ColorNameWarning
	case featurehighlight.KindError:
		style.ColorName = theme.ColorNameError
	default:
		style.ColorName = theme.ColorNameForeground
	}
	return style
}
