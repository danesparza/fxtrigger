package api

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// ListAllTriggers godoc
// @Summary List all triggers in the system
// @Description List all triggers in the system
// @Tags triggers
// @Accept  json
// @Produce  json
// @Success 200 {object} api.SystemResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /triggers [get]
func (service Service) ListAllTriggers(rw http.ResponseWriter, req *http.Request) {

	//	Get a list of files
	retval, err := service.DB.GetAllTriggers()
	if err != nil {
		err = fmt.Errorf("error getting a list of triggers: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Construct our response
	response := SystemResponse{
		Message: fmt.Sprintf("%v triggers(s)", len(retval)),
		Data:    retval,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// CreateTrigger godoc
// @Summary Create a new trigger
// @Description Create a new trigger
// @Tags triggers
// @Accept  json
// @Produce  json
// @Param trigger body api.CreateTriggerRequest true "The trigger to create"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /triggers [post]
func (service Service) CreateTrigger(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request
	request := CreateTriggerRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	If we don't have any webhooks associated, make sure we indicate that's not valid
	if len(request.WebHooks) < 1 {
		sendErrorResponse(rw, fmt.Errorf("at least one webhook must be included"), http.StatusBadRequest)
		return
	}

	//	Create the new trigger:
	newTrigger, err := service.DB.AddTrigger(request.Name, request.Description, request.GPIOPin, request.WebHooks, request.MinimumSecondsBeforeRetrigger)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Record the event:
	log.Debug().Any("request", request).Msg("Trigger created")

	//	Add the new trigger to monitoring:
	service.AddMonitor <- newTrigger

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Trigger created",
		Data:    newTrigger,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// UpdateTrigger godoc
// @Summary Update a trigger
// @Description Update a trigger
// @Tags triggers
// @Accept  json
// @Produce  json
// @Param trigger body api.UpdateTriggerRequest true "The trigger to update.  Must include trigger.id"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /triggers [put]
func (service Service) UpdateTrigger(rw http.ResponseWriter, req *http.Request) {

	//	Some state change instructions
	shouldAddMonitoring := false
	shouldRemoveMonitoring := false

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request
	request := UpdateTriggerRequest{}
	err := json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	If we don't have the trigger.id, make sure we indicate that's not valid
	if strings.TrimSpace(request.ID) == "" {
		sendErrorResponse(rw, fmt.Errorf("the trigger.id is required"), http.StatusBadRequest)
		return
	}

	//	Make sure the id exists
	trigUpdate, _ := service.DB.GetTrigger(request.ID)
	if trigUpdate.ID != request.ID {
		sendErrorResponse(rw, fmt.Errorf("trigger must already exist"), http.StatusBadRequest)
		return
	}

	//	See if 'enabled' has changed
	if trigUpdate.Enabled != request.Enabled {
		if request.Enabled {
			//	If it has, and it's now 'enabled', add the trigger to monitoring
			shouldAddMonitoring = true
		} else {
			//	If it has, and it's now 'disabled', remove the trigger from monitoring
			shouldRemoveMonitoring = true
		}
	}

	//	Only update the name if it's been passed
	if strings.TrimSpace(request.Name) != "" {
		trigUpdate.Name = request.Name
	}

	//	Only update the description if it's been passed
	if strings.TrimSpace(request.Description) != "" {
		trigUpdate.Description = request.Description
	}

	//	Enabled / disabled is always set
	trigUpdate.Enabled = request.Enabled

	//	If the GPIO pin is not zero (the default value of an int) pass it in.  Yes -- GPIO 0 is valid,
	//	but is generally reserved for special uses.  See https://pinout.xyz/pinout/pin27_gpio0#
	if request.GPIOPin != 0 {
		trigUpdate.GPIOPin = request.GPIOPin
	}

	//	This is an int. It's always going to get updated
	trigUpdate.MinimumSecondsBeforeRetrigger = request.MinimumSecondsBeforeRetrigger

	//	Only update webhooks if we've passed some in
	if len(request.WebHooks) > 0 {
		trigUpdate.WebHooks = request.WebHooks
		service.RemoveMonitor <- trigUpdate.ID
		shouldAddMonitoring = true
	}

	//	Create the new trigger:
	updatedTrigger, err := service.DB.UpdateTrigger(trigUpdate)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Record the event:
	log.Debug().Any("request", request).Msg("Trigger updated")

	//	If we have a state change, make sure to add/remove monitoring and record that event as well
	if shouldAddMonitoring {
		service.AddMonitor <- trigUpdate
		log.Debug().Str("id", trigUpdate.ID).Msg("Trigger monitoring enabled")
	}

	if shouldRemoveMonitoring {
		service.RemoveMonitor <- trigUpdate.ID
		log.Debug().Str("id", trigUpdate.ID).Msg("Trigger monitoring disabled")
	}

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Trigger updated",
		Data:    updatedTrigger,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// DeleteTrigger godoc
// @Summary Deletes a trigger in the system
// @Description Deletes a trigger in the system
// @Tags triggers
// @Accept  json
// @Produce  json
// @Param id path string true "The trigger id to delete"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Failure 503 {object} api.ErrorResponse
// @Router /triggers/{id} [delete]
func (service Service) DeleteTrigger(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	vars := mux.Vars(req)
	if vars["id"] == "" {
		err := fmt.Errorf("requires an id of a trigger to delete")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Delete the trigger
	err := service.DB.DeleteTrigger(vars["id"])
	if err != nil {
		err = fmt.Errorf("error deleting file: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Record the event:
	log.Debug().Str("id", vars["id"]).Msg("Trigger deleted")

	//	Remove the trigger from monitoring:
	service.RemoveMonitor <- vars["id"]

	//	Construct our response
	response := SystemResponse{
		Message: "Trigger deleted",
		Data:    vars["id"],
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}

// FireSingleTrigger godoc
// @Summary Fires a trigger in the system
// @Description Fires a trigger in the system
// @Tags triggers
// @Accept  json
// @Produce  json
// @Param id path string true "The trigger id to fire"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /trigger/fire/{id} [post]
func (service Service) FireSingleTrigger(rw http.ResponseWriter, req *http.Request) {

	//	Get the id from the url (if it's blank, return an error)
	vars := mux.Vars(req)
	if vars["id"] == "" {
		err := fmt.Errorf("requires an id of a trigger to fire")
		sendErrorResponse(rw, err, http.StatusBadRequest)
		return
	}

	//	Get the trigger
	trigger, err := service.DB.GetTrigger(vars["id"])
	if err != nil {
		err = fmt.Errorf("error getting trigger: %v", err)
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Call the channel to fire the event:
	service.FireTrigger <- trigger

	//	Record the event:
	log.Debug().Str("id", trigger.ID).Str("name", trigger.Name).Msg("Trigger fired")

	//	Construct our response
	response := SystemResponse{
		Message: "Trigger fired",
		Data:    trigger,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
