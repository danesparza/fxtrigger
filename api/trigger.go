package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/danesparza/fxtrigger/data"
	"github.com/danesparza/fxtrigger/event"
	"github.com/danesparza/fxtrigger/triggertype"
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
// @Param trigger body data.Trigger true "The trigger to create"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /triggers [post]
func (service Service) CreateTrigger(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request
	request := data.Trigger{}
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
// @Param trigger body data.Trigger true "The trigger to update.  Must include trigger.id"
// @Success 200 {object} api.SystemResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /triggers [put]
func (service Service) UpdateTrigger(rw http.ResponseWriter, req *http.Request) {

	//	req.Body is a ReadCloser -- we need to remember to close it:
	defer req.Body.Close()

	//	Decode the request
	request := data.Trigger{}
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

	//	If we don't have any webhooks associated, make sure we indicate that's not valid
	if len(request.WebHooks) < 1 {
		sendErrorResponse(rw, fmt.Errorf("at least one webhook must be included"), http.StatusBadRequest)
		return
	}

	//	Create the new trigger:
	updatedTrigger, err := service.DB.UpdateTrigger(request)
	if err != nil {
		sendErrorResponse(rw, err, http.StatusInternalServerError)
		return
	}

	//	Record the event:
	service.DB.AddEvent(event.TriggerUpdated, triggertype.Unknown, fmt.Sprintf("%+v", request), GetIP(req), service.HistoryTTL)

	//	Create our response and send information back:
	response := SystemResponse{
		Message: "Trigger updated",
		Data:    updatedTrigger,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
