package trigger

import (
	"bytes"
	"context"
	"fmt"
	"github.com/danesparza/fxtrigger/internal/data"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/danesparza/go-rpio"
)

// BackgroundProcess encapsulates background processing operations
type BackgroundProcess struct {
	DB         *data.Manager
	HistoryTTL time.Duration

	// FireTrigger signals a trigger should be fired
	FireTrigger chan data.Trigger

	// AddMonitor signals a trigger should be added to the list of monitored triggers
	AddMonitor chan data.Trigger

	// RemoveMonitor signals a trigger id should not be monitored anymore
	RemoveMonitor chan string
}

type monitoredTriggersMap struct {
	m       map[string]func()
	rwMutex sync.RWMutex
}

// HandleAndProcess handles system context calls and channel events to fire triggers
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
						log.Err(err).Str("TriggerID", trigger.ID).Str("HookUrl", hook.URL).Msg("Error creating request for trigger/hook")
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
						log.Err(err).Str("TriggerID", trigger.ID).Str("HookUrl", hook.URL).Msg("Error with response for trigger/hook")
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

// ListenForEvents listens to channel events to add / remove monitors
//
//	and 'fires' triggers when an event (motion / button press / time) occurs from a monitor
func (bp BackgroundProcess) ListenForEvents(systemctx context.Context) {

	//	Track our list of active event monitors.  These could be buttons or sensors
	monitoredTriggers := monitoredTriggersMap{m: make(map[string]func())}

	//	Loop and respond to channels:
	for {
		select {
		case monitorReq := <-bp.AddMonitor:
			//	This should be called when creating a trigger,
			//	when initializing the service,
			//	or when enabling a trigger (that was previously disabled)

			//	If you need to add a monitor, spin up a background goroutine to monitor that pin
			go func(cx context.Context, req data.Trigger) {

				//	Create a cancelable context from the passed (system) context
				ctx, cancel := context.WithCancel(cx)
				defer cancel()

				//	Add an entry to the map with
				//	- key: triggerid
				//	- value: the cancel function (pointer)
				//	(critical section)
				monitoredTriggers.rwMutex.Lock()
				monitoredTriggers.m[req.ID] = cancel
				monitoredTriggers.rwMutex.Unlock()

				if err := rpio.Open(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
				defer rpio.Close()

				pin := rpio.Pin(req.GPIOPin)
				pin.Mode(rpio.Input)

				//	Store the 'last reading'
				//	Initially, set it to the 'low' (no motion) state
				lr := rpio.Low
				lastTrigger := time.Unix(0, 0) // Initialize with 1/1/1970

				log.Debug().Int("GPIOPin", req.GPIOPin).Str("TriggerID", req.ID).Msg("Monitoring started")

				//	Our channel checker and sensor reader
				for {
					select {
					case <-ctx.Done():
						//	Remove ourselves from the map and exit (critical section)
						monitoredTriggers.rwMutex.Lock()
						delete(monitoredTriggers.m, req.ID)
						monitoredTriggers.rwMutex.Unlock()
						return
					case <-time.After(500 * time.Millisecond):
						//	Read from the sensor
						v := pin.Read()

						//	Latch / unlatch check
						if lr != v {
							lr = v
							currentTime := time.Now()
							diff := currentTime.Sub(lastTrigger)

							if lr == rpio.High {
								if diff.Seconds() > float64(req.MinimumSecondsBeforeRetrigger) {
									//	If it's been long enough -- reset the lrTime to now
									//	and actually trigger the item
									lastTrigger = currentTime
									log.Debug().Int("GPIOPin", req.GPIOPin).Str("TriggerID", req.ID).Msg("Motion detected.  Firing event")
									bp.FireTrigger <- req
								} else {
									log.Debug().
										Int("GPIOPin", req.GPIOPin).
										Str("TriggerID", req.ID).
										Int("MinimumSecondsBeforeRetrigger", req.MinimumSecondsBeforeRetrigger).
										Msg("Motion detected, but minimum seconds threshold not met.  Not triggering.")
								}
							}
							if lr == rpio.Low {
								log.Debug().Int("GPIOPin", req.GPIOPin).Str("TriggerID", req.ID).Msg("Motion reset")
							}
						}
					}
				}

			}(systemctx, monitorReq) // Launch the goroutine

		case removeReq := <-bp.RemoveMonitor:
			//	This should be called when removing a trigger (permanently)
			//	or when disabling a trigger

			//	Look up the item in the map and call cancel if the item exists (critical section):
			monitoredTriggers.rwMutex.Lock()
			monitorCancel, exists := monitoredTriggers.m[removeReq]

			if exists {
				log.Debug().Str("TriggerID", removeReq).Msg("Monitoring stopped")

				//	Call the context cancellation function
				monitorCancel()

				//	Remove ourselves from the map and exit
				delete(monitoredTriggers.m, removeReq)
			}
			monitoredTriggers.rwMutex.Unlock()

		case <-systemctx.Done():
			fmt.Println("Stopping trigger processor")
			return
		}
	}
}

// InitializeMonitors starts all monitoring processes
func (bp BackgroundProcess) InitializeMonitors() {

	//	Get all triggers:
	allTriggers, err := bp.DB.GetAllTriggers()
	if err != nil {
		log.Err(err).Msg("Problem getting all triggers to initialze monitors")
	}

	log.Debug().Int("TriggerCount", len(allTriggers)).Msg("Initializing monitoring")

	//	Start monitoring all enabled triggers:
	for _, trigger := range allTriggers {
		if trigger.Enabled {
			bp.AddMonitor <- trigger
		}
	}
}
