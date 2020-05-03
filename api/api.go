package api

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

type Result struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (ar *Result) IsSuccess() bool {
	return ar.Code == SUCCESS
}

func (ar *Result) ReadFrom(res *http.Response) error {
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("status: [%d], skip read body", res.StatusCode)
	}
	if b, e := ioutil.ReadAll(res.Body); e != nil {
		return e
	} else {
		if e := jsoniter.Unmarshal(b, ar); e != nil {
			return e
		}
		return nil
	}
}

type Code string

const SUCCESS Code = "SUCCESS"
const ERR400 Code = "ERR400"
const ERR500 Code = "ERR500"
