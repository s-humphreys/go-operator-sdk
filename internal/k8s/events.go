package k8s

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
)

// EventType represents the type of event for a Kubernetes resource.
// The values are concrete and cannot be changed, this is enforced by
// the Kubernetes event recorder.
type EventType int

const (
	EventTypeNormal EventType = iota
	EventTypeWarning
)

var EventTypeMap = map[EventType]string{
	EventTypeNormal:  "Normal",
	EventTypeWarning: "Warning",
}

type Event struct {
	EventType EventType
	Reason    string
	Message   string
}

// Returns an event for the Deployment resource based on the reason.
func NewEvent(object runtime.Object, recorder record.EventRecorder, event Event) {
	recorder.Event(object, EventTypeMap[event.EventType], event.Reason, event.Message)
}
