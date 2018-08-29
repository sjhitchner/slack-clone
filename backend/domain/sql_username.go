// DO NOT EDIT - Auto generated by sqltype for Username

package domain

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

func (t *Username) Scan(value interface{}) error {

	if value == nil {
		return nil
	}

	switch v := value.(type) {

	case []byte:
		*t = Username(string(v))

	case string:
		*t = Username(v)

	case *string:
		*t = Username(*v)

	default:
		return fmt.Errorf("%s Can't convert '%v' to string", reflect.TypeOf(t), value)
	}

	return nil
}

func (t Username) Value() (driver.Value, error) {

	return string(t), nil

}
