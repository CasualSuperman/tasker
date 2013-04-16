package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type apiUserResponse struct {
	Successful bool   `json:"successful"`
	Err        string `json:"err,omitempty"`
	code       int
}

func (r apiUserResponse) Json() []byte {
	resp, _ := json.Marshal(r)
	return resp
}

func (r apiUserResponse) Type() string {
	return "text/plain"
}

func (r apiUserResponse) Code() int {
	return r.code
}

func (r *apiUserResponse) Succeed() {
	r.Successful = true
	r.code = http.StatusOK
}

func (r *apiUserResponse) Fail(err error) {
	r.code = http.StatusInternalServerError
	r.Err = err.Error()
}

func defaultUserResponse() apiUserResponse {
	return apiUserResponse{
		Successful: false,
		Err:        "",
		code:       500,
	}
}

type eventsList struct {
	Events []Event      `json:"events"`
	StartDate time.Time `json:"startDate"`
	EndDate time.Time   `json:"endDate"`
}

func (e *eventsList) Json() []byte {
	data, _ := json.Marshal(e)
	return data
}

func (e *eventsList) Type() string {
	return "text/plain"
}

func (e *eventsList) Code() int {
	return http.StatusOK
}
