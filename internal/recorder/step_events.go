package recorder

// RecordStepOp describes how the frontend should reconcile recorder state.
type RecordStepOp string

const (
	RecordStepUpsert   RecordStepOp = "upsert"
	RecordStepDelete   RecordStepOp = "delete"
	RecordStepReset    RecordStepOp = "reset"
	RecordStepSnapshot RecordStepOp = "snapshot"
)

// RecordStepEvent is the canonical live-recording step payload.
type RecordStepEvent struct {
	Op    RecordStepOp
	Index int
	Line  string
	Lines []string
}

type StepNotifier func(event RecordStepEvent)

func notifyUpsert(notify StepNotifier, index int, bareLine string) {
	if notify == nil || bareLine == "" {
		return
	}
	notify(RecordStepEvent{
		Op:    RecordStepUpsert,
		Index: index,
		Line:  FormatRecordedGherkinLine(bareLine, pickRecordedStepKeyword(index)),
	})
}

func notifyDelete(notify StepNotifier, index int) {
	if notify == nil || index < 0 {
		return
	}
	notify(RecordStepEvent{Op: RecordStepDelete, Index: index})
}

func notifyReset(notify StepNotifier) {
	if notify == nil {
		return
	}
	notify(RecordStepEvent{Op: RecordStepReset})
}

func notifySnapshot(notify StepNotifier, steps []RecordedStep) {
	if notify == nil {
		return
	}
	lines := FormatRecordedStepsAsGherkin(steps)
	if len(lines) == 0 {
		return
	}
	notify(RecordStepEvent{Op: RecordStepSnapshot, Lines: lines})
}
