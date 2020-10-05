package business

import (
	"encoding/json"

	"github.com/wiseco/core-platform/api"
	"github.com/wiseco/core-platform/services"
	"github.com/wiseco/core-platform/services/csp/business"
	csp "github.com/wiseco/core-platform/services/csp/services"
)

func handleCreateNotes(businessID string, r api.APIRequest) (api.APIResponse, error) {
	var requestBody business.NotesCreate

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}
	requestBody.BusinessID = businessID

	note, err := business.NewNotesService(csp.NewSRRequest(r.CognitoID)).Create(requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	notesJSON, err := json.Marshal(note)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(notesJSON), false)
}

func handleNotesList(businessID string, limit, offset int, r api.APIRequest) (api.APIResponse, error) {

	list, err := business.NewNotesService(csp.NewSRRequest(r.CognitoID)).List(businessID, limit, offset)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	notesJSON, err := json.Marshal(list)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(notesJSON), false)
}

func handleNoteByID(businessID, notesID string, r api.APIRequest) (api.APIResponse, error) {

	note, err := business.NewNotesService(csp.NewSRRequest(r.CognitoID)).ByID(businessID, notesID)

	if err != nil {
		_, isENF := err.(*services.ErrorNotFound)
		if isENF {
			return api.NotFound(r)
		}

		return api.InternalServerError(r, err)
	}

	noteJSON, err := json.Marshal(note)
	if err != nil {
		return api.InternalServerError(r, err)
	}
	return api.Success(r, string(noteJSON), false)
}

func handleNoteUpdate(businessID, notesID string, r api.APIRequest) (api.APIResponse, error) {
	var requestBody business.NotesUpdate

	err := json.Unmarshal([]byte(r.Body), &requestBody)
	if err != nil {
		return api.BadRequest(r, err)
	}

	note, err := business.NewNotesService(csp.NewSRRequest(r.CognitoID)).Update(notesID, requestBody)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	noteJSON, err := json.Marshal(note)
	if err != nil {
		return api.InternalServerError(r, err)
	}

	return api.Success(r, string(noteJSON), false)
}

func handleDeleteNote(businessID, notesID string, r api.APIRequest) (api.APIResponse, error) {
	if err := business.NewNotesService(csp.NewSRRequest(r.CognitoID)).Delete(businessID, notesID); err != nil {
		return api.InternalServerError(r, err)
	}
	status := struct {
		Message string `json:"message"`
	}{Message: "Success"}
	resp, _ := json.Marshal(status)
	return api.Success(r, string(resp), false)
}
