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
		traces: make(chan map[int][]Span, 1), // buffer = 1, because we need one map to be in the chan without blocking it 
	}
	t.traces<- make(map[int][]Span)
	return t
}

type TelemetryService struct {
	traces chan map[int][]Span // Shared resource, use chan so we can manage its access
}

func (t *TelemetryService) AddSpan(id int, action string) {
	// Gets traces from chan
	traces := <-t.traces
	trace := traces[id] // Gets trace from traces map
	
	newSpan := Span{
		id: id,
		ts: time.Now(),
		action: action,
	}

	// Append to slice
	trace = append(trace, newSpan)
	traces[newSpan.id] = trace // replace with new slice
	t.traces<- traces // send to chan for another goroutine to use it
}

func (t *TelemetryService) PrintTelemetry() {
	traces := <-t.traces
	//TODO: Melhorar o print
	for _, t := range traces {
		for _, s := range t {
			fmtText := fmt.Sprintf("id: %d, timestamp (UnixMilli): %d, ts.DateTime: %s, action: %s", s.id, s.ts.UnixMilli(), s.ts.Format(time.DateTime), s.action)
			fmt.Println(fmtText)
		}
	}
}