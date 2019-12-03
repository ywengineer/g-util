package sql

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Replace(s, old string, new []interface{}) string {
	// Compute number of replacements.
	if m := strings.Count(s, old); m == 0 || len(new) == 0 {
		return s // avoid allocation
	} else {
		if len(new) != m {
			return s
		}
		// Apply replacements to buffer.
		buf := bytes.NewBuffer(make([]byte, 0, len(s)+22*len(new)))
		//
		start := 0
		for i := 0; i < m; i++ {
			j := start
			if len(old) == 0 {
				if i > 0 {
					_, wid := utf8.DecodeRuneInString(s[start:])
					j += wid
				}
			} else {
				j += strings.Index(s[start:], old)
			}
			_, _ = buf.WriteString(s[start:j])
			//w += wd
			_, _ = buf.WriteString(ToString(new[i]))
			//w += wd
			start = j + len(old)
		}
		//
		buf.WriteString(s[start:])
		//
		return buf.String() //string(t[0:w])
	}
}

func ToString(v interface{}) string {
	if v == nil {
		return "''"
	}
	switch v.(type) {
	case uintptr:
		panic(fmt.Sprintf("ptr not supported: %v : %v\n", reflect.TypeOf(v).Elem(), reflect.ValueOf(v).Elem()))
	case int:
		return strconv.Itoa(v.(int))
	case int32:
		return fmt.Sprintf("%d", v.(int32))
	case int64:
		return fmt.Sprintf("%d", v.(int64))
	case bool:
		if v.(bool) == true {
			return "1"
		}
		return "0"
	case string:
		return fmt.Sprintf("'%s'", v.(string))
	case NullString:
		vv := v.(NullString)
		if vv.Valid {
			return fmt.Sprintf("'%s'", vv.String)
		}
	case sql.NullString:
		vv := v.(sql.NullString)
		if vv.Valid {
			return fmt.Sprintf("'%s'", vv.String)
		}
	case NullTime:
		vv := v.(NullTime)
		if vv.Valid {
			return fmt.Sprintf("'%s'", vv.String())
		}
	}
	return "''"
}
