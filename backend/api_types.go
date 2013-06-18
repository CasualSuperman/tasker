package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type apiUserResponse struct {
	Successful bool   `json:"success"`
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

type calendar struct {
	Cid   int    `json:"cid"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type calendarList []calendar

func (c *calendarList) Json() []byte {
	data, _ := json.Marshal(c)
	return data
}

func (c calendarList) Type() string {
	return "text/plain"
}

func (c calendarList) Code() int {
	return http.StatusOK
}

type eventsList struct {
	Events    []Event   `json:"events"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
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

type apiFormResponse struct {
	Success   bool     `json:"success"`
	ErrFields []string `json:"errFields,omitempty"`
	ErrMsgs   []string `json:"errMsgs,omitempty"`
}

func (fr *apiFormResponse) Json() []byte {
	data, _ := json.Marshal(fr)
	return data
}
func (fr *apiFormResponse) Type() string {
	return "text/plain"
}
func (fr *apiFormResponse) Code() int {
	return http.StatusOK
}

type apiRawResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func (ar apiRawResponse) Json() []byte {
	data, _ := json.Marshal(ar.Data)
	return data
}
func (ar apiRawResponse) Code() int {
	return http.StatusOK
}
func (ar apiRawResponse) Type() string {
	return "text/plain"
}
