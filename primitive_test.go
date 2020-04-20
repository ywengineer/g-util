package g_util

import (
	"github.com/ywengineer/g-util/util"
	"testing"
)

func TestFormat(t *testing.T) {
	m := map[string]interface{}{"id": 77965241413144576, "reply": "这是文本"}
	if id, ok := m["id"]; ok {
		docID, e := util.Int2String(id)
		t.Log(docID, e)
	}
}
