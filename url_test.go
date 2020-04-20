package g_util

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/url"
	"testing"
)

func TestUrl(t *testing.T) {
	_, e := url.ParseRequestURI("www.nbaidu.com/abdkd/ia?a=2")
	if e != nil {
		t.Errorf("%v", e)
	}
	u, e := url.Parse("http://user:password@www.baidu.com")
	if e != nil {
		t.Errorf("%v", e)
	} else {
		t.Logf("%s, %s, %s", u.Scheme, u.Host, u.RequestURI())
		t.Logf("%s://%s%s", u.Scheme, u.Host, u.RequestURI())
		c := fasthttp.Client{}
		statusCode, body, err := c.Get(nil, fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, u.RequestURI()))
		t.Logf("%d, %s, %v", statusCode, string(body), err)
	}
}

func TestJ(t *testing.T) {
	time := "2019-12-12"
	body, _ := jsoniter.MarshalToString([]string{"kdaf", "iefa"})
	t.Log(`{"time":"` + time + `", "words":` + body + `}`)

	s := `["Expand","your","time","range","스스로","결정하고선택하는","제도","One","more","indices","you’re","looking","contains","Your","query","may","match","anything","current","any","data","all","currently","selected","You","can","try","changing","one","which","다양한","분야의","창작자들이","새로운","콘텐츠와","비즈니스를","만들어세상에","그들의","이름을","결재만","며칠이","걸리는","네이버에는","없습니다","결재의","대부분을","조직장이","승인하는","번거로운","과정을","줄였으며","노트북이나","PC","등의","업무기기도","개인의","업무","스타일에","맞게스스로","선택할","있습니다"]`
	var sl []string
	if err := jsoniter.UnmarshalFromString(s, &sl); err != nil {
		t.Logf("%v", err)
	} else {
		t.Logf("%v", sl)
	}
}
