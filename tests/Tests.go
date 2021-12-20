package tests

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type HttpServeAble interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type Request struct {
	Router HttpServeAble
	Method string
	Path string
	Body interface{}
	Response interface{}
}

func NewTestRequest(t *testing.T,req Request)  {

	bodyBytes , err := json.Marshal(req.Body)
	if !assert.Nil(t, err) {
		return
	}

	response, err := http.NewRequest(http.MethodPost, "/member/login", bytes.NewReader(bodyBytes))
	if !assert.Nil(t, err) {
		return
	}


	recorder := httptest.NewRecorder()

	req.Router.ServeHTTP(recorder, response)
	//fmt.Println(recorder.Body.String())


	err = json.NewDecoder(recorder.Body).Decode(req.Response)
	if !assert.Nil(t, err){
		return
	}
	//fmt.Println("mp: ",mp)
}