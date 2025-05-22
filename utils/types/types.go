package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
	"time"
)

type CustomDate time.Time

func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*cd = CustomDate(t)
	return nil
}

func (cd CustomDate) String() string {
	return time.Time(cd).Format("2006-01-02")
}

func (cd CustomDate) ToTime() time.Time {
	return time.Time(cd)
}

type JSONNullString sql.NullString

func (jns JSONNullString) MarshalJSON() ([]byte, error) {
	if !jns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(jns.String)
}

func (jns *JSONNullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		jns.String = ""
		jns.Valid = false
		return nil
	}
	err := json.Unmarshal(data, &jns.String)
	jns.Valid = err == nil
	return err
}

func (jns *JSONNullString) Scan(value interface{}) error {
	var ns sql.NullString
	if err := ns.Scan(value); err != nil {
		return err
	}
	jns.String = ns.String
	jns.Valid = ns.Valid
	return nil
}

func (jns JSONNullString) Value() (driver.Value, error) {
	if !jns.Valid {
		return nil, nil
	}
	return jns.String, nil
}
