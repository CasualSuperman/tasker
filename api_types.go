package main

import (
	"encoding/json"
	"net/http"
)

type apiUserResponse struct {
	successful bool
	err        string `json:",omitempty"`
	code       int    `json:"-"`
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
	r.successful = true
	r.code = http.StatusOK
}

func defaultUserResponse() apiUserResponse {
	return apiUserResponse{
		successful: false,
		err:        "",
		code:       500,
	}
}
