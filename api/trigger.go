package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/danesparza/fxtrigger/event"
	"github.com/danesparza/fxtrigger/triggertype"
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
	service.DB.AddEvent(event.TriggerCreated, triggertype.Unknown, fmt.Sprintf("%+v", request), GetIP(req), service.HistoryTTL)

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

	//	If we don't have any webhooks associated, make sure we indicate that's not valid
	if len(request.WebHooks) < 1 {
		sendErrorResponse(rw, fmt.Errorf("at least one webhook must be included"), http.StatusBadRequest)
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

	//	Update the trigger:
	trigUpdate.Name = request.Name
	trigUpdate.Description = request.Description
	trigUpdate.Enabled = request.Enabled
	trigUpdate.GPIOPin = request.GPIOPin
	trigUpdate.MinimumSecondsBeforeRetrigger = request.MinimumSecondsBeforeRetrigger
	trigUpdate.WebHooks = request.WebHooks

	//	Create the new trigger:
	updatedTrigger, err := service.DB.UpdateTrigger(trigUpdate)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Record the event:
	service.DB.AddEvent(event.TriggerUpdated, triggertype.Unknown, fmt.Sprintf("%+v", request), GetIP(req), service.HistoryTTL)

	//	If we have a state change, make sure to add/remove monitoring and record that event as well
	if shouldAddMonitoring {
		service.AddMonitor <- trigUpdate
		service.DB.AddEvent(event.TriggerUpdated, triggertype.Unknown, fmt.Sprintf("Trigger ID %s is now marked as 'enabled' ... adding to active event monitors", trigUpdate.ID), GetIP(req), service.HistoryTTL)
	}

	if shouldRemoveMonitoring {
		service.RemoveMonitor <- trigUpdate.ID
		service.DB.AddEvent(event.TriggerUpdated, triggertype.Unknown, fmt.Sprintf("Trigger ID %s is no longer marked as 'enabled' ... removing from active event monitors", trigUpdate.ID), GetIP(req), service.HistoryTTL)
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
	service.DB.AddEvent(event.TriggerDeleted, triggertype.Unknown, vars["id"], GetIP(req), service.HistoryTTL)

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
	service.DB.AddEvent(event.TriggerFired, triggertype.Unknown, fmt.Sprintf("Trigger ID: %s / Name: %s", trigger.ID, trigger.Name), GetIP(req), service.HistoryTTL)

	//	Construct our response
	response := SystemResponse{
		Message: "Trigger fired",
		Data:    trigger,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
