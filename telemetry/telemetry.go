package telemetry

import (
	"fmt"
	"time"
)

type Telemetry interface {
	AddSpan(newSpan Span)
	PrintTelemetry()
}

type Span struct {
	id 		int
	ts 		time.Time
	action 	string
}

func NewTelemetryService() TelemetryService {
	t := TelemetryService{
		traces: make(map[int][]Span),
	}
	return t
}

type TelemetryService struct {
	traces map[int][]Span
}

func (t *TelemetryService) AddSpan(id int, ts time.Time, action string) {
	newSpan := Span{
		id: id,
		ts: ts,
		action: action,
	}
	trace := t.traces[newSpan.id]

	trace = append(trace, newSpan)
	t.traces[newSpan.id] = trace
}

func (t *TelemetryService) PrintTelemetry() {
	//TODO: Melhorar o print
	for _, t := range t.traces {
		for _, s := range t {
			fmtText := fmt.Sprintf("id: %d, timestamp (Unix): %d, ts.DateTime: %s, action: %s", s.id, s.ts.Unix(), s.ts.Format(time.DateTime), s.action)
			fmt.Println(fmtText)
		}
	}
}