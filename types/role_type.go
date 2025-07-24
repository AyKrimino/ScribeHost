package types

import (
	"database/sql/driver"
	"fmt"
)

type RoleType string

const (
	AdminRole  RoleType = "admin"
	AuthorType RoleType = "author"
	ReaderRole RoleType = "reader"
)

func (r *RoleType) Scan(value interface{}) error {
	if value == nil {
		*r = AuthorType
		return nil
	}

	switch v := value.(type) {
	case string:
		*r = RoleType(v)
	case []byte:
		*r = RoleType(v)
	default:
		return fmt.Errorf("cannot scan %T into RoleType", value)
	}
	return nil
}

func (r RoleType) Value() (driver.Value, error) {
	return string(r), nil
}

func (RoleType) GormDataType() string {
	return "enum('admin','author','reader')"
}
