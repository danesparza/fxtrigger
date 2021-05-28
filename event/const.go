package event

const (
	// SystemStartup event is when the system has started up
	SystemStartup = "System startup"

	// TriggerCreated event is when a trigger has been created
	TriggerCreated = "Trigger created"

	// TriggerUpdated event is when a trigger has been updated
	TriggerUpdated = "Trigger updated"

	// TriggerDeleted event is when a trigger has been removed
	TriggerDeleted = "Trigger deleted"

	// TriggerFired event is when a trigger has been fired
	TriggerFired = "Trigger fired"

	// TriggerError event is when there was an error processing a trigger
	TriggerError = "Trigger error"

	// SystemShutdown event is when the system is shutting down
	SystemShutdown = "System Shutdown"
)
