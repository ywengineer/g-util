package es

import (
	"context"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	jsoniter "github.com/json-iterator/go"
	"github.com/ywengineer/g-util/util"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestNewESClient(t *testing.T) {
	l := util.NewLogger("../logs/log.log", 32, 64, 7, zapcore.DebugLevel, true)
	es := NewESClient([]string{"http://8.129.217.210:9200"}, l)
	req := esapi.IndicesGetSettingsRequest{
		Index:             []string{"log_order"},
		Name:              nil,
		AllowNoIndices:    nil,
		ExpandWildcards:   "",
		FlatSettings:      nil,
		IgnoreUnavailable: nil,
		IncludeDefaults:   nil,
		Local:             nil,
		MasterTimeout:     0,
		Pretty:            false,
		Human:             true,
		ErrorTrace:        false,
		FilterPath:        nil,
		Header:            nil,
	}
	// Perform the request with the client.
	res, err := req.Do(context.Background(), es)
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
