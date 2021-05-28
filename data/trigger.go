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
	ID                            string    `json:"id"`                            // Unique Trigger ID
	Enabled                       bool      `json:"enabled"`                       // Trigger enabled or not
	Created                       time.Time `json:"created"`                       // Trigger create time
	Name                          string    `json:"name"`                          // The trigger name
	Description                   string    `json:"description"`                   // Additional information about the trigger
	GPIOPin                       string    `json:"gpiopin"`                       // The GPIO pin the sensor or button is on
	WebHooks                      []WebHook `json:"webhooks"`                      // The webhooks to send when triggered
	MinimumSecondsBeforeRetrigger int       `json:"minimumsecondsbeforeretrigger"` // Minimum time (in seconds) before a retrigger
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

	newTrigger := Trigger{
		ID:                            xid.New().String(), // Generate a new id
		Created:                       time.Now(),
		Enabled:                       true,
		Name:                          name,
		Description:                   description,
		GPIOPin:                       gpiopin,
		WebHooks:                      webhooks,
		MinimumSecondsBeforeRetrigger: minimumsleep,
	}

	//	Serialize to JSON format
	encoded, err := json.Marshal(newTrigger)
	if err != nil {
		return retval, fmt.Errorf("problem serializing the data: %s", err)
	}

	//	Save it to the database:
	err = store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(GetKey("Trigger", newTrigger.ID), string(encoded), &buntdb.SetOptions{})
		return err
	})

	//	If there was an error saving the data, report it:
	if err != nil {
		return retval, fmt.Errorf("problem saving the trigger: %s", err)
	}

	//	Set our retval:
	retval = newTrigger

	//	Return our data:
	return retval, nil
}

// AddTrigger adds a trigger to the system
func (store Manager) UpdateTrigger(updatedTrigger Trigger) (Trigger, error) {

	//	Our return item
	retval := Trigger{}

	//	Serialize to JSON format
	encoded, err := json.Marshal(updatedTrigger)
	if err != nil {
		return retval, fmt.Errorf("problem serializing the data: %s", err)
	}

	//	Save it to the database:
	err = store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(GetKey("Trigger", updatedTrigger.ID), string(encoded), &buntdb.SetOptions{})
		return err
	})

	//	If there was an error saving the data, report it:
	if err != nil {
		return retval, fmt.Errorf("problem saving the trigger: %s", err)
	}

	//	Set our retval:
	retval = updatedTrigger

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

// GetAllTriggers gets all triggers in the system
func (store Manager) GetAllTriggers() ([]Trigger, error) {
	//	Our return item
	retval := []Trigger{}

	//	Set our prefix
	prefix := GetKey("Trigger")

	//	Iterate over our values:
	err := store.systemdb.View(func(tx *buntdb.Tx) error {
		tx.Descend(prefix, func(key, val string) bool {

			if len(val) > 0 {
				//	Create our item:
				item := Trigger{}

				//	Unmarshal data into our item
				bval := []byte(val)
				if err := json.Unmarshal(bval, &item); err != nil {
					return false
				}

				//	Add to the array of returned users:
				retval = append(retval, item)
			}

			return true
		})
		return nil
	})

	//	If there was an error, report it:
	if err != nil {
		return retval, fmt.Errorf("problem getting the list of triggers: %s", err)
	}

	//	Return our data:
	return retval, nil
}

// DeleteTrigger deletes a trigger from the system
func (store Manager) DeleteTrigger(id string) error {

	//	Remove it from the database:
	err := store.systemdb.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(GetKey("Trigger", id))
		return err
	})

	//	If there was an error removing the data, report it:
	if err != nil {
		return fmt.Errorf("problem removing the trigger: %s", err)
	}

	//	Return our data:
	return nil
}
