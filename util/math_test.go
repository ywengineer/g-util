package util

import (
	"reflect"
	"testing"
)

func TestMath(t *testing.T) {
	t.Logf("%v\n", ValidMail("ywengineer@gmail.com"))
	t.Logf("%v\n", ValidMail("akei.ei-@gmail.com"))
	t.Logf("%v\n", ValidMail("iakf_ie.12@gmail.com"))
	t.Logf("%v\n", ValidMail("iakf_ie.12@gmail.cn"))
	t.Logf("%v\n", ValidMail("iakf_ie.12@gmail.online"))
	t.Logf("%v\n", ValidMail("iakf_ie.12@gmail.store"))
}

func TestMap(t *testing.T) {
	m := []map[string]interface{}{
		{"a": 1},
		{"b": 1},
	}
	t.Log(reflect.TypeOf(m[0]).Kind().String())
}
