package gflag

import (
	"encoding"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/spf13/pflag"
	"golang.org/x/exp/constraints"
)

var (
	// type guards: Value implements pflag.Value
	_ pflag.Value              = &Value[string]{}
	_ pflag.Value              = &Value[time.Duration]{}
	_ encoding.TextMarshaler   = &Value[string]{}
	_ encoding.TextUnmarshaler = &Value[string]{}
)

type (
	// FlaggableTypes is a type constraint that holds all types supported by pflag, besides primitive types.
	FlaggableTypes interface {
		time.Duration |
			time.Time |
			net.IP |
			net.IPNet |
			net.IPMask |
			// for extended types
			~struct{}
	}

	// FlaggablePrimitives is a type constraint that holds all primitive types supported by pflag, and then some.
	FlaggablePrimitives interface {
		constraints.Integer |
			constraints.Float |
			constraints.Complex |
			~string |
			~bool |
			~[]byte // aka: ~[]uint8
	}

	// Value is a generic type that implements github.com/spf13/pflag.Value.
	//
	// The underlying value, as T, may be retrieved using GetValue().
	Value[T FlaggablePrimitives | FlaggableTypes] struct {
		Value       *T
		NoOptDefVal string
		*options
	}
)

// NewFlagValue constructs a generic flag compatible with github.com/spf13/pflag.Value.
//
// Since the flag type is inferred from the underlying data type, some flexibility allowed by pflag is not
// always possible at this point.
//
// For example, when T = []byte, NewFlagValue adopts by default the semantics of the pflag.BytesHex flag.
//
// Similarly, when T = int, we adopt by default the semantics of pflag.Int and not pflag.Count.
//
// In order to cover the full range of semantics offered by the pflag package, some options are available.
func NewFlagValue[T FlaggablePrimitives | FlaggableTypes](addr *T, defaultValue T, opts ...Option) *Value[T] {
	if addr == nil {
		addr = new(T)
	}

	m := &Value[T]{
		Value:   addr,
		options: defaultOptions(opts),
	}
	*m.Value = defaultValue

	// bool flag implies true when set without arg
	v := any(m.Value)
	if _, isBool := v.(*bool); isBool {
		m.NoOptDefVal = "true"
	}
	if _, isInt := v.(*int); isInt && m.semanticsForInt == intIsCount {
		m.NoOptDefVal = "+1"
	}

	return m
}

// GetValue returns the underlying value of the flag.
func (m Value[T]) GetValue() T {
	return *m.Value
}

// GetNoOptDefVal returns the default value to consider whenever there is no argument added to the flag.
//
// Example:
// for a Value[bool], NoOptDefVal defaults to "true".
func (m Value[T]) GetNoOptDefVal() string {
	asAny := any(m.Value)
	if withNoOpt, ok := asAny.(interface{ GetNoOptDefVal() string }); ok {
		return withNoOpt.GetNoOptDefVal()
	}

	return m.NoOptDefVal
}

// IsBoolFlag indicates to the "flag" standard lib package that this is a boolean flag.
func (m Value[T]) IsBoolFlag() bool {
	asAny := any(m.Value)
	switch v := asAny.(type) {
	case interface{ IsBoolFlag() bool }:
		return v.IsBoolFlag()
	case *bool:
		return true
	default:
		return false
	}
}

// Get implements the flag.Getter interface for programs using this Value from the
// standard "flag" package.
func (m Value[T]) Get() any {
	return m.Value
}

// String knows how to yield a string representation of type T.
func (m Value[T]) String() string {
	if m.Value == nil {
		return ""
	}

	asAny := any(m.Value)
	switch v := asAny.(type) {
	case pflag.Value:
		return v.String()
	case *string:
		return *v
	case *int:
		return formatInt(*v)
	case *int8:
		return formatInt(*v)
	case *int16:
		return formatInt(*v)
	case *int32:
		return formatInt(*v)
	case *int64:
		return formatInt(*v)
	case *uint:
		return formatUint(*v)
	case *uint8:
		return formatUint(*v)
	case *uint16:
		return formatUint(*v)
	case *uint32:
		return formatUint(*v)
	case *uint64:
		return formatUint(*v)
	case *float32:
		return floatFormatter[float32](32)(*v)
	case *float64:
		return floatFormatter[float64](64)(*v)
	case *complex64:
		return complexFormatter[complex64](64)(*v)
	case *complex128:
		return complexFormatter[complex128](128)(*v)
	case *bool:
		return strconv.FormatBool(*v)
	case *[]byte:
		return bytesFormatter(*v, m.semanticsForBytes)
	case *time.Duration:
		return v.String()
	case *time.Time:
		return timeFormatter(m.timeFormats[0])(*v)
	case *net.IP:
		return v.String()
	case *net.IPNet:
		return v.String()
	case *net.IPMask:
		return v.String()
	case fmt.Stringer:
		return v.String()
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}
}

// Set knows how to config a string representation of the Value into a type T.
func (m *Value[T]) Set(strValue string) error {
	return m.set(strValue)
}

func (m *Value[T]) set(strValue string) error {
	asAny := any(m.Value)
	switch v := asAny.(type) {
	case pflag.Value:
		return v.Set(strValue)
	case *string:
		val := strValue
		*m.Value = *cast[T](&val)
	case *int:
		if m.semanticsForInt == intIsCount {
			m.NoOptDefVal = "+1"

			if strValue == "+1" {
				*v++
				*m.Value = *cast[T](v)

				return nil
			}
		}

		val, err := intParser[int](0)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *int8:
		val, err := intParser[int8](8)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *int16:
		val, err := intParser[int16](16)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *int32:
		val, err := intParser[int32](32)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *int64:
		val, err := intParser[int64](64)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *uint:
		val, err := uintParser[uint](0)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *uint8:
		val, err := uintParser[uint8](8)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *uint16:
		val, err := uintParser[uint16](16)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *uint32:
		val, err := uintParser[uint32](32)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *uint64:
		val, err := uintParser[uint64](64)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *float32:
		val, err := floatParser[float32](32)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *float64:
		val, err := floatParser[float64](64)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *complex64:
		val, err := complexParser[complex64](64)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *complex128:
		val, err := complexParser[complex128](128)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *bool:
		val, err := strconv.ParseBool(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
		m.NoOptDefVal = "true"
	case *[]byte:
		val, err := parseBytes(strings.TrimSpace(strValue), m.semanticsForBytes)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *time.Duration:
		val, err := time.ParseDuration(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *time.Time:
		val, err := timeParser(m.timeFormats)(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)

	case *net.IP:
		val, err := parseIP(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *net.IPMask:
		val, err := parseIPMask(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	case *net.IPNet:
		val, err := parseIPNet(strValue)
		if err != nil {
			return err
		}

		*m.Value = *cast[T](&val)
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}

	return nil
}

func (m *Value[T]) Type() string {
	asAny := any(m.Value)
	switch v := asAny.(type) {
	case pflag.Value:
		return v.Type()
	case *int:
		switch m.semanticsForInt {
		case intIsCount:
			return "count"
		default:
			return "int"
		}
	case *time.Duration:
		return "duration"
	case *time.Time:
		return "time"
	case *net.IP:
		return "ip"
	case *net.IPMask:
		return "ipMask"
	case *net.IPNet:
		return "ipNet"
	case *[]byte:
		switch m.semanticsForBytes {
		case bytesIsBase64:
			return "bytesBase64"
		default:
			return "bytesHex"
		}
	default:
		return fmt.Sprintf("%T", *m.Value)
	}
}

// MarshalFlag implements github.com/jessevdk/go-flags.Marshaler interface
func (m *Value[T]) MarshalFlag() (string, error) {
	return m.String(), nil
}

func (m *Value[T]) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// UnmarshalFlag implements github.com/jessevdk/go-flags.Unmarshaler interface
func (m *Value[T]) UnmarshalFlag(value string) error {
	return m.set(value)
}

func (m *Value[T]) UnmarshalText(text []byte) error {
	return m.set(string(text))
}

func cast[U any, T any](v *T) *U {
	return (*U)(unsafe.Pointer(v))
}
