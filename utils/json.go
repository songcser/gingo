package utils

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

func JSONToMap(content interface{}) []map[string]any {
	var name []map[string]any
	if marshalContent, err := json.Marshal(content); err != nil {
		fmt.Println(err)
	} else {
		//if err := json.Unmarshal(marshalContent, &name); err != nil {
		//	fmt.Println(err)
		//}
		d := json.NewDecoder(bytes.NewReader(marshalContent))
		d.UseNumber()
		if err := d.Decode(&name); err != nil {
			fmt.Println(err)
		}
	}
	return name
}

var timeFormat = "2006-01-02 15:04:05"

// JsonTime format json time field by myself
type JsonTime time.Time

func (t *JsonTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = JsonTime(time.Time{})
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), loc)
	*t = JsonTime(now)
	return
}

// MarshalJSON on JsonTime format Time field with Y-m-d H:i:s
func (t JsonTime) MarshalJSON() ([]byte, error) {
	if time.Time(t).IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", time.Time(t).Format(timeFormat))
	return []byte(formatted), nil
}

// Value insert timestamp into mysql need this function.
func (t JsonTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if time.Time(t).UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return time.Time(t), nil
}

// Scan value of time.Time
func (t *JsonTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = JsonTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *JsonTime) Time() time.Time {
	return time.Time(*t)
}
