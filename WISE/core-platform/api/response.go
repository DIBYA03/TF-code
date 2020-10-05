/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type APIResponse struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
	Request           APIRequest          `json:"request"`
}

type PagedResponse struct {
	Count int         `json:"count"`
	Rows  interface{} `json:"items"`
}

type ResponseError struct {
	Code      int    `json:"code"`
	CodeDesc  string `json:"codeDesc"`
	ErrorDesc string `json:"errorDesc"`
	RawError  error  `json:"error"`
}

func (r *ResponseError) Error() string {
	return r.RawError.Error()
}

func (resp *APIResponse) Duration() int64 {
	duration := time.Since(resp.Request.StartedAt) / time.Millisecond
	return int64(duration)
}

func NewJSONResponseHeader() map[string]string {
	return map[string]string{
		"Content-Type":                "application/json",
		"Cache-Control":               "private, max-age=0",
		"Access-Control-Allow-Origin": "*",
	}
}

func Success(request APIRequest, responseBody string, isBase64Encoded bool) (APIResponse, error) {

	return APIResponse{
			StatusCode:      http.StatusOK,
			Headers:         NewJSONResponseHeader(),
			Body:            responseBody,
			IsBase64Encoded: isBase64Encoded,
		},
		nil
}

func logError(err error) error {
	if err == nil {
		return nil
	}

	b, e := json.Marshal(err)
	if e != nil {
		log.Println(string(b))
	} else {
		log.Println(err)
	}

	// Temp fix for PLATFORM-218
	errText := strings.ToLower(err.Error())
	if strings.Contains(errText, "sql") || strings.Contains(errText, "pq") || strings.Contains(errText, "pg") {
		return errors.New(http.StatusText(http.StatusInternalServerError))
	}

	return err
}

// HTTP 400 Bad Request
func BadRequestError(request APIRequest, err error) (APIResponse, error) {
	err = logError(err)
	b, _ := json.Marshal(
		ResponseError{
			Code:      http.StatusBadRequest,
			CodeDesc:  http.StatusText(http.StatusBadRequest),
			ErrorDesc: err.Error(),
			RawError:  err,
		},
	)

	return APIResponse{
			StatusCode:      http.StatusBadRequest,
			Headers:         NewJSONResponseHeader(),
			Body:            string(b),
			IsBase64Encoded: false,
		},
		nil
}

func BadRequest(request APIRequest, err error) (APIResponse, error) {
	return BadRequestError(request, err)
}

// HTTP 401 Unauthorized
func UnauthorizedError(request APIRequest, err error) (APIResponse, error) {
	err = logError(err)
	return APIResponse{
			StatusCode:      http.StatusUnauthorized,
			Headers:         NewJSONResponseHeader(),
			Body:            fmt.Sprintf("[\"%s\"]", http.StatusText(http.StatusUnauthorized)),
			IsBase64Encoded: false,
		},
		nil
}

// HTTP 403 Forbidden
func ForbiddenError(request APIRequest, err error) (APIResponse, error) {
	err = logError(err)
	return APIResponse{
			StatusCode:      http.StatusForbidden,
			Headers:         NewJSONResponseHeader(),
			Body:            fmt.Sprintf("[\"%s\"]", http.StatusText(http.StatusForbidden)),
			IsBase64Encoded: false,
		},
		nil
}

// HTTP 404 Not Found
func NotFoundError(request APIRequest, err error) (APIResponse, error) {
	err = logError(err)
	return APIResponse{
			StatusCode: http.StatusNotFound,
			Headers:    NewJSONResponseHeader(),
			Body:       fmt.Sprintf("[\"%s\"]", http.StatusText(http.StatusNotFound)),
		},
		nil
}

func NotFound(request APIRequest) (APIResponse, error) {
	return NotFoundError(request, nil)
}

// HTTP 405 Not Allowed
func NotAllowedError(request APIRequest, err error) (APIResponse, error) {
	err = logError(err)
	return APIResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Headers:    NewJSONResponseHeader(),
			Body:       fmt.Sprintf("[\"%s\"]", http.StatusText(http.StatusMethodNotAllowed)),
		},
		nil
}

func NotSupported(request APIRequest) (APIResponse, error) {
	return NotAllowedError(request, nil)
}

// HTTP 500 Internal Server Error
func InternalServerError(request APIRequest, err error) (APIResponse, error) {
	err = logError(err)
	b, _ := json.Marshal(
		ResponseError{
			Code:      http.StatusInternalServerError,
			CodeDesc:  http.StatusText(http.StatusInternalServerError),
			ErrorDesc: err.Error(),
			RawError:  err,
		},
	)

	return APIResponse{
			StatusCode:      http.StatusInternalServerError,
			Headers:         NewJSONResponseHeader(),
			Body:            string(b),
			IsBase64Encoded: false,
		},
		nil
}

// HTTP 501 Not Implemented
func NotImplementedError(request APIRequest, err error) (APIResponse, error) {
	err = logError(err)
	return APIResponse{
			StatusCode: http.StatusNotImplemented,
			Headers:    NewJSONResponseHeader(),
			Body:       fmt.Sprintf("[\"%s\"]", http.StatusText(http.StatusNotImplemented)),
		},
		nil
}

//ProxyErrorResponse is a custom api response to use for ther proxy, we want to send the actual
//status code and any response BBVA might give us back when an error accurs
func ProxyErrorResponse(request APIRequest, body []byte, statusCode int) (APIResponse, error) {
	return APIResponse{
		StatusCode:      statusCode,
		Headers:         NewJSONResponseHeader(),
		Body:            string(body),
		IsBase64Encoded: false,
	}, nil
}
