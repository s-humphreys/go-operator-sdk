package k8s

import (
	"fmt"

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
func NewEvent(crd runtime.Object, recorder record.EventRecorder, event Event) {
	recorder.Event(crd, EventTypeMap[event.EventType], event.Reason, event.Message)
}

// NewCreatedEvent creates a new Kubernetes resource created event on the CRD.
func NewCreatedEvent(crd runtime.Object, recorder record.EventRecorder, kind string, name string) {
	NewEvent(crd, recorder, Event{
		EventType: EventTypeNormal,
		Reason:    kind + "Created",
		Message:   fmt.Sprintf("%s %s has been created successfully", kind, name),
	})
}

// NewCreateErrorEvent creates a new Kubernetes resource creation error event on the CRD.
func NewCreateErrorEvent(crd runtime.Object, recorder record.EventRecorder, kind string, name string) {
	NewEvent(crd, recorder, Event{
		EventType: EventTypeWarning,
		Reason:    kind + "CreateError",
		Message:   fmt.Sprintf("An error occurred whilst creating %s %s", kind, name),
	})
}

// NewUpdatedEvent creates a new Kubernetes resource updated event on the CRD.
func NewUpdatedEvent(crd runtime.Object, recorder record.EventRecorder, kind string, name string) {
	NewEvent(crd, recorder, Event{
		EventType: EventTypeNormal,
		Reason:    kind + "Updated",
		Message:   fmt.Sprintf("%s %s has been updated successfully", kind, name),
	})
}

// NewUpdateErrorEvent creates a new Kubernetes resource update error event on the CRD.
func NewUpdateErrorEvent(crd runtime.Object, recorder record.EventRecorder, kind string, name string) {
	NewEvent(crd, recorder, Event{
		EventType: EventTypeWarning,
		Reason:    kind + "UpdateError",
		Message:   fmt.Sprintf("An error occurred whilst updating %s %s", kind, name),
	})
}

// NewOutOfSyncEvent creates a new Kubernetes resource out of sync event on the CRD.
func NewOutOfSyncEvent(crd runtime.Object, recorder record.EventRecorder, kind string, name string) {
	NewEvent(crd, recorder, Event{
		EventType: EventTypeWarning,
		Reason:    kind + "OutOfSync",
		Message:   fmt.Sprintf("%s %s is out of sync with the desired spec", kind, name),
	})
}
