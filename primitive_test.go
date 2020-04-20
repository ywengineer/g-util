package g_util

import (
	"encoding/json"
	"github.com/ywengineer/g-util/util"
	"reflect"
	"testing"
)

func TestFormat(t *testing.T) {

	m := map[string]interface{}{"id": 77965241413144576, "reply": "这是文本"}
	if id, ok := m["id"]; ok {
		docID, e := util.Int2String(id)
		t.Log(docID, e)
	}
	t.Log(reflect.TypeOf(new(interface{})).String())
	t.Log(reflect.TypeOf(10).String())
	t.Log(reflect.TypeOf([]string{}).String())
	t.Log(reflect.TypeOf(json.Number("13")).String())
	t.Log(reflect.TypeOf(json.Number("13")).Name())
}
