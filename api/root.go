package api

import (
	"encoding/json"
	"fmt"
	data2 "github.com/danesparza/fxtrigger/internal/data"
	"net/http"
	"time"
)

// Service encapsulates API service operations
type Service struct {
	DB         *data2.Manager
	StartTime  time.Time
	HistoryTTL time.Duration

	// FireTrigger signals a trigger should be fired
	FireTrigger chan data2.Trigger

	// AddMonitor signals a trigger should be added to the list of monitored triggers
	AddMonitor chan data2.Trigger

	// RemoveMonitor signals a trigger id should not be monitored anymore
	RemoveMonitor chan string
}

// CreateTriggerRequest is a request to create a new trigger
type CreateTriggerRequest struct {
	Name                          string          `json:"name"`                          // The trigger name
	Description                   string          `json:"description"`                   // Additional information about the trigger
	GPIOPin                       int             `json:"gpiopin"`                       // The GPIO pin the sensor or button is on
	WebHooks                      []data2.WebHook `json:"webhooks"`                      // The webhooks to send when triggered
	MinimumSecondsBeforeRetrigger int             `json:"minimumsecondsbeforeretrigger"` // Minimum time (in seconds) before a retrigger
}

// UpdateTriggerRequest is a request to update a trigger
type UpdateTriggerRequest struct {
	ID                            string          `json:"id"`                            // Unique Trigger ID
	Enabled                       bool            `json:"enabled"`                       // Trigger enabled or not
	Name                          string          `json:"name"`                          // The trigger name
	Description                   string          `json:"description"`                   // Additional information about the trigger
	GPIOPin                       int             `json:"gpiopin"`                       // The GPIO pin the sensor or button is on
	WebHooks                      []data2.WebHook `json:"webhooks"`                      // The webhooks to send when triggered
	MinimumSecondsBeforeRetrigger int             `json:"minimumsecondsbeforeretrigger"` // Minimum time (in seconds) before a retrigger
}

// SystemResponse is a response for a system request
type SystemResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ErrorResponse represents an API response
type ErrorResponse struct {
	Message string `json:"message"`
}

// Used to send back an error:
func sendErrorResponse(rw http.ResponseWriter, err error, code int) {
	//	Our return value
	response := ErrorResponse{
		Message: "Error: " + err.Error()}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(code)
	json.NewEncoder(rw).Encode(response)
}

// ShowUI redirects to the /ui/ url path
func ShowUI(rw http.ResponseWriter, req *http.Request) {
	// http.Redirect(rw, req, "/ui/", 301)
	fmt.Fprintf(rw, "Hello, world - UI")
}

// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
