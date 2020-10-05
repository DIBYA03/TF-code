package golang

import (
	"database/sql/driver"
	"errors"
	"fmt"
)

// ToBool forces conversion from a Boolean to a go bool value
// Unspecified will convert to false
func (b Boolean) ToBool() bool {
	return b == Boolean_B_TRUE
}

func (b Boolean) IsTrue() bool {
	return b == Boolean_B_TRUE
}

func (b Boolean) IsFalse() bool {
	return b == Boolean_B_FALSE
}

func (b Boolean) IsUnspecified() bool {
	return b == Boolean_B_UNSPECIFIED
}

// NewBoolean from go bool
func NewBoolean(b bool) Boolean {
	if b {
		return Boolean_B_TRUE
	}

	return Boolean_B_FALSE
}

// SQL value marshaller
func (b Boolean) Value() (driver.Value, error) {
	// Database should use boolean type
	switch b {
	case Boolean_B_TRUE:
		return true, nil
	case Boolean_B_FALSE:
		return false, nil
	default:
		return false, errors.New("boolean value is unspecified")
	}
}

// SQL scanner
func (b *Boolean) Scan(value interface{}) error {
	// Default to false
	*b = Boolean_B_FALSE

	// Boolean value must be specified
	if value == nil {
		return errors.New("invalid boolean value")
	}

	val, ok := value.(bool)
	if !ok {
		return fmt.Errorf("unable to scan value %v", value)
	}

	if val {
		*b = Boolean_B_TRUE
	}

	return nil
}
