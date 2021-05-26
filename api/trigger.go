package api

import (
	"encoding/json"
	"fmt"
	"net/http"
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
		Message: fmt.Sprintf("%v file(s)", len(retval)),
		Data:    retval,
	}

	//	Serialize to JSON & return the response:
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(response)
}
