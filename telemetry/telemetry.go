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
		traces: make(chan map[int][]Span, 1),
	}
	t.traces<- make(map[int][]Span)
	return t
}

type TelemetryService struct {
	traces chan map[int][]Span
}

func (t *TelemetryService) AddSpan(id int, ts time.Time, action string) {
	traces := <-t.traces
	trace := traces[id]
	
	newSpan := Span{
		id: id,
		ts: ts,
		action: action,
	}

	trace = append(trace, newSpan)
	traces[newSpan.id] = trace
	t.traces<- traces
}

func (t *TelemetryService) PrintTelemetry() {
	traces := <-t.traces
	//TODO: Melhorar o print
	for _, t := range traces {
		for _, s := range t {
			fmtText := fmt.Sprintf("id: %d, timestamp (Unix): %d, ts.DateTime: %s, action: %s", s.id, s.ts.Unix(), s.ts.Format(time.DateTime), s.action)
			fmt.Println(fmtText)
		}
	}
}