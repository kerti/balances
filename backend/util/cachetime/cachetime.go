package cachetime

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/guregu/null"
)

func toCacheTime(value interface{}) (*CacheTime, error) {
	switch x := value.(type) {
	case string:
		val, err := strconv.ParseInt(x, 10, 64)
		if err != nil {
			return nil, err
		}
		t := time.UnixMilli(val)
		ct := CacheTime(t)
		return &ct, nil
	case int64:
		t := time.UnixMilli(x)
		ct := CacheTime(t)
		return &ct, nil
	case time.Time:
		ct := CacheTime(x)
		return &ct, nil
	case *time.Time:
		if x != nil {
			ct := CacheTime(*x)
			return &ct, nil
		}
		return nil, nil
	case CacheTime:
		return &x, nil
	case *CacheTime:
		return x, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("null: cannot scan type %T into time.Time: %v", value, value)
	}
}

// CacheTime is the data type to represent time in cache, represented by Unix millis.
// It is castable to time.Time
type CacheTime time.Time

// MarshalJSON implements json.Marshaler interface
func (ct CacheTime) MarshalJSON() ([]byte, error) {
	ms := (time.Time)(ct).UnixNano() / int64(time.Millisecond)
	return json.Marshal(ms)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (ct *CacheTime) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	ts, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return
	}
	*ct = CacheTime(time.Unix(0, ts*int64(time.Millisecond)))
	return
}

// Scan implements the Scanner interface.
func (ct *CacheTime) Scan(value interface{}) error {
	newCt, err := toCacheTime(value)
	if err != nil {
		return err
	}
	if newCt != nil {
		*ct = *newCt
	}
	return nil
}

// Value implements the driver Valuer interface.
func (ct CacheTime) Value() (driver.Value, error) {
	return time.Time(ct), nil
}

// Time returns the time object inside the CacheTime
func (ct CacheTime) Time() time.Time { return (time.Time)(ct) }

// NCacheTime is nullable CacheTime.
// It is castable to null.Time
type NCacheTime null.Time

// MarshalJSON implements json.Marshaler interface
func (ct NCacheTime) MarshalJSON() ([]byte, error) {
	if !ct.Valid {
		return []byte("null"), nil
	}
	ms := (time.Time)(ct.Time).UnixNano() / int64(time.Millisecond)
	return json.Marshal(ms)
}

// UnmarshalJSON implements json.Unmarshaler interface
func (ct *NCacheTime) UnmarshalJSON(data []byte) (err error) {
	str := string(data)
	if str == "null" {
		ct.Valid = false
		return
	}
	ts, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return
	}
	ct.Time = time.Unix(0, ts*int64(time.Millisecond))
	ct.Valid = true
	return
}

// Scan implements the Scanner interface.
func (ct *NCacheTime) Scan(value interface{}) error {
	newCt, err := toCacheTime(value)
	if err != nil {
		return err
	}
	if newCt != nil {
		ct.Time = newCt.Time()
		ct.Valid = true
	} else {
		ct.Valid = false
	}
	return nil
}

// Value implements the driver Valuer interface.
func (ct NCacheTime) Value() (driver.Value, error) {
	if !ct.Valid {
		return nil, nil
	}
	return ct.Time, nil
}
