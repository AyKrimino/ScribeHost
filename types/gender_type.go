package types

import (
	"database/sql/driver"
	"fmt"
)

type GenderType string

const (
	GenderMale   GenderType = "male"
	GenderFemale GenderType = "female"
)

func (g *GenderType) Scan(value interface{}) error {
	if value == nil {
		*g = GenderMale
		return nil
	}

	switch v := value.(type) {
	case string:
		*g = GenderType(v)
	case []byte:
		*g = GenderType(v)
	default:
		return fmt.Errorf("cannot scan %T into RoleType", value)
	}
	return nil
}

func (g GenderType) Value() (driver.Value, error) {
	return string(g), nil
}

func (GenderType) GormDataType() string {
	return "enum('male','female')"
}
