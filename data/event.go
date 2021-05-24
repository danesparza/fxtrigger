package data

import "time"

// Event represents an event in the system.  These events
// can be logged or passed (as meta information) to other systems
type Event struct {
	ID          string    `json:"id"`          // Unique Event ID
	Created     time.Time `json:"created"`     // Event creation time
	SourceIP    string    `json:"ip"`          // Source IP address of the event
	EventType   string    `json:"eventtype"`   // One of: System startup, Trigger created, Trigger fired, Trigger deleted, System shutdown
	TriggerType string    `json:"triggertype"` // The type of trigger involved: Motion, Button, Time
	Details     string    `json:"details"`     // Additional information (like the trigger name involved)
}
