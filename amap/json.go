package amap

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type dbJSON map[string]interface{}

func (j dbJSON) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *dbJSON) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &j)
}

func (j *dbJSON) FromDB(bytes []byte) error {
	return json.Unmarshal(bytes, j)
}

func (j *dbJSON) ToDB() (bytes []byte, err error) {
	bytes, err = json.Marshal(j)
	return
}

// FromDB FromDB
func (p Poi) FromDB(bytes []byte) error {
	return json.Unmarshal(bytes, &p)
}

// ToDB ToDB
func (p Poi) ToDB() (bytes []byte, err error) {
	bytes, err = json.Marshal(p)
	return
}

// FromDB FromDB
func (p *Detail) FromDB(bytes []byte) error {
	return json.Unmarshal(bytes, &p)
}

// ToDB ToDB
func (p *Detail) ToDB() (bytes []byte, err error) {
	bytes, err = json.Marshal(p)
	return
}
