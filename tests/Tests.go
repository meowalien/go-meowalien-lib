package tests

import (
	"bytes"
	"encoding/json"
	"github.com/meowalien/go-meowalien-lib/errs"
	"net/http"
	"net/http/httptest"
)

type HttpServeAble interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type Request struct {
	Router   HttpServeAble
	Method   string
	Path     string
	Body     interface{}
	Response interface{}
}

func NewTestRequest(req Request) (httpCode int, err error) {
	bodyBytes, err := json.Marshal(req.Body)
	if err != nil {
		err = errs.New(err)
		return
	}

	response, err := http.NewRequest(req.Method, req.Path, bytes.NewReader(bodyBytes))
	if err != nil {
		err = errs.New(err)
		return
	}

	recorder := httptest.NewRecorder()
	req.Router.ServeHTTP(recorder, response)
	err = json.NewDecoder(recorder.Body).Decode(req.Response)
	if err != nil {
		err = errs.New(err)
		return
	}
	return recorder.Code, nil
}
