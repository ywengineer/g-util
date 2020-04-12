package es

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/ywengineer/g-util/util"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestJSON(t *testing.T) {
	jsonApi := jsoniter.Config{UseNumber: true, EscapeHTML: true}.Froze()
	m := map[string]interface{}{
		"a": 1,
		"b": "123",
	}
	text, _ := jsonApi.MarshalToString(m)
	t.Log(text + "\n")
}

func TestNewESClient(t *testing.T) {
	l := util.NewLogger("../logs/log.log", 32, 64, 7, zapcore.DebugLevel, true)
	es := NewESClient([]string{"http://8.129.217.210:9200"}, l)
	gs := es.Indices.GetSettings
	res, err := gs(gs.WithIndex("log_order"), gs.WithHuman(), gs.WithContext(context.Background()))
	//
	if err != nil {
		t.Fatalf("Error getting response: %v", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		t.Logf("[%s] Error get index setting", res.Status())
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := jsoniter.NewDecoder(res.Body).Decode(&r); err != nil {
			t.Logf("Error parsing the response body: %v", err)
		} else if tx, err := jsoniter.MarshalToString(r); err != nil {
			// Print the response status and indexed document version.
			t.Logf("[%s] %v", res.Status(), err)
		} else {
			t.Logf("[%s] %s", res.Status(), tx)
		}
	}
}
