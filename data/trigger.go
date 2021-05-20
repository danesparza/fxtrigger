package data

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/xid"
	"github.com/tidwall/buntdb"
)

// Trigger represents sensor/button trigger information.
type Trigger struct {
	ID                          string    `json:"id"`                          // Unique Trigger ID
	Created                     time.Time `json:"created"`                     // File create time
	Name                        string    `json:"name"`                        // The trigger name
	Description                 string    `json:"description"`                 // Additional information about the trigger
	GPIOPin                     string    `json:"gpiopin"`                     // The GPIO pin the sensor or button is on
	WebHooks                    []WebHook `json:"webhooks"`                    // The webhooks to send when triggered
	MinimumSleepBeforeRetrigger int       `json:"minimumsleepbeforeretrigger"` // Minimum sleep time (in seconds) before a retrigger
}

// WebHook represents a notification message sent to an endpoint
type WebHook struct {
	URL         string `json:"url"`         // The URL to connect to
	ContentType string `json:"contenttype"` // The requested content type of the response (usually application/json)
	HTTPVerb    string `json:"httpverb"`    // HTTP verb (GET/PUT/POST/DELETE/etc)
	HTTPHeaders string `json:"httpheaders"` // The HTTP headers to send
	HTTPBody    []byte `json:"httpbody"`    // The HTTP body to send.  This can be empty
}

// AddTrigger adds a trigger to the system
func (store Manager) AddTrigger(name, description, gpiopin string, webhooks []WebHook, minimumsleep int) (Trigger, error) {

	//	Our return item
	retval := Trigger{}

	newFile := Trigger{
		ID:                          xid.New().String(), // Generate a new id
		Created:                     time.Now(),
		Name:                        name,
		Description:                 description,
		GPIOPin:                     gpiopin,
		WebHooks:                    webhooks,
		MinimumSleepBeforeRetrigger: minimumsleep,
	}

	//	Serialize to JSON format
	encoded, err := json.Marshal(newFile)
	if err != nil {
		return retval, fmt.Errorf("problem serializing the data: %s", err)
	}

	//	Save it to the database:
	err = store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(GetKey("Trigger", newFile.ID), string(encoded), &buntdb.SetOptions{})
		return err
	})

	//	If there was an error saving the data, report it:
	if err != nil {
		return retval, fmt.Errorf("problem saving the trigger: %s", err)
	}

	//	Set our retval:
	retval = newFile

	//	Return our data:
	return retval, nil
}

// GetTrigger gets information about a single trigger in the system based on its id
func (store Manager) GetTrigger(id string) (Trigger, error) {
	//	Our return item
	retval := Trigger{}

	//	Find the item:
	err := store.systemdb.View(func(tx *buntdb.Tx) error {

		val, err := tx.Get(GetKey("Trigger", id))
		if err != nil {
			return err
		}

		if len(val) > 0 {
			//	Unmarshal data into our item
			if err := json.Unmarshal([]byte(val), &retval); err != nil {
				return err
			}
		}

		//	If we get to this point and there is no error...
		return nil
	})

	//	If there was an error, report it:
	if err != nil {
		return retval, fmt.Errorf("problem getting the trigger: %s", err)
	}

	//	Return our data:
	return retval, nil
}
