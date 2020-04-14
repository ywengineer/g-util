package g_util

import (
	jsoniter "github.com/json-iterator/go"
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	_, e := url.ParseRequestURI("www.nbaidu.com/abdkd/ia?a=2")
	if e != nil {
		t.Errorf("%v", e)
	}
	_, e = url.ParseRequestURI("http://www.nbaidu.com/abdkd/ia")
	if e != nil {
		t.Errorf("%v", e)
	}
}

func TestJ(t *testing.T) {
	time := "2019-12-12"
	body, _ := jsoniter.MarshalToString([]string{"kdaf", "iefa"})
	t.Log(`{"time":"` + time + `", "words":` + body + `}`)
}
