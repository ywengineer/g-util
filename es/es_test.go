package es

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/ywengineer/g-util/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"testing"
)

var tl = util.NewLogger("../logs/log.log", 32, 64, 7, zapcore.DebugLevel, true)
var es = NewESClient([]string{"http://8.129.217.210:9200"}, tl)

func TestJSON(t *testing.T) {
	jsonApi := jsoniter.Config{UseNumber: true, EscapeHTML: true}.Froze()
	m := map[string]interface{}{
		"a": 1,
		"b": "123",
	}
	text, _ := jsonApi.MarshalToString(m)
	t.Log(text + "\n")
}

func TestIndicesAnalyzer(t *testing.T) {
	analyzer := es.Indices.Analyze
	res, err := analyzer(analyzer.WithIndex("log_player_create"),
		analyzer.WithContext(context.Background()),
		analyzer.WithBody(strings.NewReader(`{"analyzer":"text_analyzer","text":["我是一个大坏蛋", "하오지식 전문가 | 지식Q&A | 지식공유 | 고민있어요 조선족 사이트 ,한민족 사이트. 지구인"]}`)))
	//
	if err != nil {
		t.Fatalf("Error getting analyzer response: %v", err)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := jsoniter.NewDecoder(res.Body).Decode(&r); err != nil {
			t.Logf("Error parsing the analyzer response body: %v", err)
		} else {
			var words []string
			//
			if tokens, ok := r["tokens"]; ok {
				for _, token := range tokens.([]interface{}) {
					words = append(words, token.(map[string]interface{})["token"].(string))
				}
			}
			//
			tl.Info("analyzer res", zap.String("status", res.Status()), zap.Any("words", words))
		}
	}
}

func TestNewESClient(t *testing.T) {
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
