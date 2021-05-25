package trigger

import (
	"context"
	"fmt"

	"github.com/danesparza/fxtrigger/data"
)

// HandleAndProcess handles system context calls and channel events to fire trigger
func HandleAndProcess(systemctx context.Context, firetrigger chan data.Trigger) {

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
