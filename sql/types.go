package sql

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
)

type NullTime mysql.NullTime

func (t *NullTime) UnmarshalJSON(data []byte) (err error) {
	if data == nil || len(data) == 0 {
		t.Valid = false
		return nil
	}
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), time.Local)
	t.Time, t.Valid = time.Time(now), true
	return nil
}

func (t *NullTime) MarshalJSON() ([]byte, error) {
	if t == nil || !t.Valid {
		return []byte{}, nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.String())), nil
}

func (t *NullTime) String() string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(TimeFormat)
}

type NullString sql.NullString

func (t *NullString) UnmarshalJSON(data []byte) (err error) {
	if data == nil || len(data) == 0 {
		t.String, t.Valid = "", false
		return nil
	}
	t.String, t.Valid = string(data[:]), true
	return nil
}

func (t *NullString) MarshalJSON() ([]byte, error) {
	if t == nil || !t.Valid || len(t.String) == 0 {
		return []byte{}, nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.String)), nil
}
