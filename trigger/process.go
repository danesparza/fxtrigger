package trigger

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danesparza/fxtrigger/data"
	"github.com/danesparza/fxtrigger/event"
	"github.com/danesparza/fxtrigger/triggertype"
)

// BackgroundProcess encapsulates background processing operations
type BackgroundProcess struct {
	DB         *data.Manager
	HistoryTTL time.Duration

	// FireTrigger signals a trigger should be fired
	FireTrigger chan data.Trigger
}

// HandleAndProcess handles system context calls and channel events to fire trigger
func (bp BackgroundProcess) HandleAndProcess(systemctx context.Context) {

	//	Loop and respond to channels:
	for {
		select {
		case trigReq := <-bp.FireTrigger:
			//	As we get a request on a channel to fire a trigger...
			//	Create a goroutine
			go func(cx context.Context, trigger data.Trigger) {

				//	Loop through the associated webhooks
				for _, hook := range trigger.WebHooks {
					//	Fire each of them...

					//	First, build the initial request with the verb, url and body (if the body exists)
					req, err := http.NewRequestWithContext(systemctx, http.MethodPost, hook.URL, bytes.NewBuffer(hook.Body))
					if err != nil {
						bp.DB.AddEvent(event.TriggerError, triggertype.Unknown, fmt.Sprintf("Error creating request for trigger/hook %s/%s: %v", trigger.ID, hook.URL, err), "", bp.HistoryTTL)
						continue //	Go to the next hook
					}

					//	Then, set our initial content-type header
					req.Header.Set("Content-Type", "application/json")

					//	Next, set any custom headers
					for k, v := range hook.Headers {
						req.Header.Set(k, v)
					}

					//	Finally, send the request
					client := &http.Client{Timeout: time.Second * 10}
					resp, err := client.Do(req)
					if err != nil {
						bp.DB.AddEvent(event.TriggerError, triggertype.Unknown, fmt.Sprintf("Error with response for trigger/hook %s/%s: %v", trigger.ID, hook.URL, err), "", bp.HistoryTTL)
						//	'continue' doesn't really matter here -- we're already at the end of this loop
					}
					defer resp.Body.Close()
				}

			}(systemctx, trigReq) // Launch the goroutine
		case <-systemctx.Done():
			fmt.Println("Stopping trigger processor")
			return
		}
	}
}

// ListenForEvents monitors the collection of triggers and reloads based on changes to the trigger collection,
//	listens for events on those triggers,
//	and 'fires' triggers when an event (motion / button press / time) occurs
func (bp BackgroundProcess) ListenForEvents(systemctx context.Context, firetrigger chan data.Trigger) {

	//	Loop and respond to channels:
	for {
		select {
		case trigReq := <-firetrigger:
			//	As we get a request on a channel to fire a trigger...
			//	Create a goroutine
			go func(cx context.Context, req data.Trigger) {

				//	Loop through the associated webhooks

				//	Fire each of them

			}(systemctx, trigReq) // Launch the goroutine
		case <-systemctx.Done():
			fmt.Println("Stopping trigger processor")
			return
		}
	}
}
