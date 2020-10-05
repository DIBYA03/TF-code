/********************************************************************
 * Copyright 2019 Wise Company
 ********************************************************************/

package api

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (request APIRequest) GetPathParam(key string) string {
	v, _ := request.PathParameters[key]
	return v
}

func (request APIRequest) GetPathParamInt(key string) (int, error) {
	var num int
	_, scanErr := fmt.Sscanf(request.GetPathParam(key), "%d", &num)
	if scanErr != nil {
		return 0, scanErr
	}

	return num, nil
}

func (request APIRequest) GetQueryParam(key string) string {
	v, _ := request.QueryStringParameters[key]
	return v
}

func (request APIRequest) GetQueryIntParam(key string) (int, error) {
	return strconv.Atoi(request.QueryStringParameters[key])
}

func (request APIRequest) GetQueryIntParamWithDefault(key string, d int) (int, error) {
	value, exists := request.QueryStringParameters[key]
	if !exists {
		return d, errors.New("Invalid parameter")
	}

	var intVal int

	intVal, _ = strconv.Atoi(value)

	return intVal, nil
}

func (request APIRequest) SingleHeaderValue(key string) string {
	kl := strings.ToLower(key)
	for k, v := range request.Headers {
		if strings.ToLower(k) == kl {
			return v
		}
	}

	return ""
}
func (request APIRequest) MultiHeaderValue(key string) []string {
	kl := strings.ToLower(key)
	for k, vals := range request.MultiValueHeaders {
		if strings.ToLower(k) == kl {
			return vals
		}
	}

	return []string{}
}
