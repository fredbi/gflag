package gflag

import (
	"bytes"
	"encoding"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

var (
	_ mapOfValues              = &MapValue[string]{}
	_ mapOfValues              = &MapValue[time.Duration]{}
	_ encoding.TextMarshaler   = &MapValue[string]{}
	_ encoding.TextUnmarshaler = &MapValue[string]{}
)

type (
	// MapValue is a generic type that implements github.com/spf13/pflag.Value and MapValue.
	//
	// It implements flags as maps of type map[string]T.
	//
	// The underlying value, as map[string]T, may be retrieved using GetMapValue().
	MapValue[T FlaggablePrimitives | FlaggableTypes] struct {
		Value   *map[string]T
		changed bool
		*options
	}

	mapOfValues interface {
		pflag.Value
		MapValuable
	}

	MapValuable interface {
		GetMap() map[string]string
	}
)

// NewFlagMapValue constructs a generic flag compatible with github.com/spf13/pflag.Value.
//
// It replaces pflag.StringToInt(), pflag.StringToInt64(), pflag.StringToString(). Being generic, it can build a map
// of any type, e.g. map[string]net.IP, map[string]time.Duration, map[string]float64...
func NewFlagMapValue[T FlaggablePrimitives | FlaggableTypes](addr *map[string]T, defaultValue map[string]T, opts ...Option) *MapValue[T] {
	m := &MapValue[T]{
		Value:   addr,
		options: defaultOptions(opts),
	}
	*m.Value = defaultValue

	return m
}

// GetValue returns the underlying value of the flag.
func (m MapValue[T]) GetMapValue() map[string]T {
	return *m.Value
}

func (m *MapValue[T]) String() string {
	return writeMapAsString(m.GetMap())
}

// Set knows how to config a string representation of the Value into a type T.
func (m *MapValue[T]) Set(strValue string) error {
	if !m.changed {
		// reset any default value
		initial := make(map[string]T)
		*m.Value = initial
		if err := m.set(strValue); err != nil {
			return err
		}

		m.changed = true

		return nil
	}

	// handle multiple occurences of the same flag with append semantics
	if err := m.set(strValue); err != nil {
		return err
	}

	return nil
}

func (m *MapValue[T]) Type() string {
	asAny := any(m.Value)
	switch v := asAny.(type) {
	case pflag.Value:
		return v.Type()
	case *map[string]string:
		return "stringToString"
	case *map[string]bool:
		return "stringToBool"
	case *map[string]int:
		return "stringToInt"
	case *map[string]int8:
		return "stringToInt8"
	case *map[string]int16:
		return "stringToInt16"
	case *map[string]int32:
		return "stringToInt32"
	case *map[string]int64:
		return "stringToInt64"
	case *map[string]uint:
		return "stringToUint"
	case *map[string]uint16:
		return "stringToUint16"
	case *map[string]uint32:
		return "stringToUint32"
	case *map[string]uint64:
		return "stringToUint64"
	case *map[string]float32:
		return "stringToFloat32"
	case *map[string]float64:
		return "stringToFloat64"
	case *map[string]complex64:
		return "stringToComplex64"
	case *map[string]complex128:
		return "stringToComplex128"
	case *map[string]time.Duration:
		return "stringToDuration"
	case *map[string]time.Time:
		return "stringToTime"
	case *map[string]net.IP:
		return "stringToIp"
	case *map[string]net.IPNet:
		return "stringToIpNet"
	case *map[string]net.IPMask:
		return "stringToIpMask"
	default:
		return fmt.Sprintf("%T", *m.Value)
	}
}

func (m *MapValue[T]) set(strValues string) error {
	asAny := any(m.Value)

	switch v := asAny.(type) {
	case *map[string]string:
		mapValues, err := readAsMap(strValues)
		if err != nil {
			return err
		}
		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]bool:
		mapValues, err := buildMapFromParser(strValues, strconv.ParseBool)
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]int:
		mapValues, err := buildMapFromParser(strValues, intParser[int](0))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]int8:
		mapValues, err := buildMapFromParser(strValues, intParser[int8](8))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]int16:
		mapValues, err := buildMapFromParser(strValues, intParser[int16](16))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]int32:
		mapValues, err := buildMapFromParser(strValues, intParser[int32](32))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]int64:
		mapValues, err := buildMapFromParser(strValues, intParser[int64](64))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]uint:
		mapValues, err := buildMapFromParser(strValues, uintParser[uint](0))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]uint16:
		mapValues, err := buildMapFromParser(strValues, uintParser[uint16](16))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]uint32:
		mapValues, err := buildMapFromParser(strValues, uintParser[uint32](32))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]uint64:
		mapValues, err := buildMapFromParser(strValues, uintParser[uint64](64))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]float32:
		mapValues, err := buildMapFromParser(strValues, floatParser[float32](32))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]float64:
		mapValues, err := buildMapFromParser(strValues, floatParser[float64](64))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]complex64:
		mapValues, err := buildMapFromParser(strValues, complexParser[complex64](64))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]complex128:
		mapValues, err := buildMapFromParser(strValues, complexParser[complex128](128))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]time.Duration:
		mapValues, err := buildMapFromParser(strValues, time.ParseDuration)
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]time.Time:
		mapValues, err := buildMapFromParser(strValues, timeParser(m.timeFormats))
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]net.IP:
		mapValues, err := buildMapFromParser(strValues, parseIP)
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]net.IPNet:
		mapValues, err := buildMapFromParser(strValues, parseIPNet)
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	case *map[string]net.IPMask:
		mapValues, err := buildMapFromParser(strValues, parseIPMask)
		if err != nil {
			return err
		}

		*v = mapValues
		*m.Value = *cast[map[string]T](v)
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}

	return nil
}

// GetMap return a map[string]string representation of the slice values.
func (m *MapValue[T]) GetMap() map[string]string {
	asAny := any(m.Value)

	switch v := asAny.(type) {
	case MapValuable:
		return v.GetMap()
	case *map[string]string:
		return *v
	case *map[string]bool:
		return buildMapWithFormatter(*v, strconv.FormatBool)
	case *map[string]int:
		return buildMapWithFormatter(*v, formatInt[int])
	case *map[string]int8:
		return buildMapWithFormatter(*v, formatInt[int8])
	case *map[string]int16:
		return buildMapWithFormatter(*v, formatInt[int16])
	case *map[string]int32:
		return buildMapWithFormatter(*v, formatInt[int32])
	case *map[string]int64:
		return buildMapWithFormatter(*v, formatInt[int64])
	case *map[string]uint:
		return buildMapWithFormatter(*v, formatUint[uint])
	case *map[string]uint16:
		return buildMapWithFormatter(*v, formatUint[uint16])
	case *map[string]uint32:
		return buildMapWithFormatter(*v, formatUint[uint32])
	case *map[string]uint64:
		return buildMapWithFormatter(*v, formatUint[uint64])
	case *map[string]float32:
		return buildMapWithFormatter(*v, floatFormatter[float32](32))
	case *map[string]float64:
		return buildMapWithFormatter(*v, floatFormatter[float64](64))
	case *map[string]complex64:
		return buildMapWithFormatter(*v, complexFormatter[complex64](64))
	case *map[string]complex128:
		return buildMapWithFormatter(*v, complexFormatter[complex128](128))
	case *map[string]time.Duration:
		return buildMapWithFormatter(*v, formatStringer[time.Duration])
	case *map[string]time.Time:
		return buildMapWithFormatter(*v, timeFormatter(m.timeFormats[0]))
	case *map[string]net.IP:
		return buildMapWithFormatter(*v, formatStringer[net.IP])
	case *map[string]net.IPNet:
		return buildMapWithFormatter(*v, formatIPnet)
	case *map[string]net.IPMask:
		return buildMapWithFormatter(*v, formatStringer[net.IPMask])
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}
}

func (m *MapValue[T]) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// MarshalFlag implements go-flags Marshaller interface
func (m *MapValue[T]) MarshalFlag() (string, error) {
	return m.String(), nil
}

func (m *MapValue[T]) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte(`[`))
	text = bytes.TrimSuffix(text, []byte(`]`))

	return m.set(string(text))
}

// UnmarshalFlag implements go-flags Unmarshaller interface
func (m *MapValue[T]) UnmarshalFlag(value string) error {
	value = strings.TrimPrefix(value, `[`)
	value = strings.TrimSuffix(value, `]`)

	return m.set(value)
}

func buildMapFromParser[T any](strValue string, parser func(string) (T, error)) (map[string]T, error) {
	mapValues, e := readAsMap(strValue)
	if e != nil {
		return nil, e
	}
	out := make(map[string]T, len(mapValues))

	for k, v := range mapValues {
		parsed, err := parser(v)
		if err != nil {
			return nil, err
		}
		out[k] = parsed
	}

	return out, nil
}

func buildMapWithFormatter[T any](mapValues map[string]T, formatter func(T) string) map[string]string {
	out := make(map[string]string, len(mapValues))

	for k, v := range mapValues {
		out[k] = formatter(v)
	}

	return out
}

func readAsMap(in string) (map[string]string, error) {
	strValues := rmQuote.Replace(in)
	keyValuePairs := strings.Split(strValues, ",")
	out := make(map[string]string, len(keyValuePairs))

	for _, pair := range keyValuePairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return nil, fmt.Errorf("%s must be formatted as key=value", pair)
		}

		out[kv[0]] = kv[1]
	}

	return out, nil
}

func writeMapAsString(mapValues map[string]string) string {
	var (
		buf        bytes.Buffer
		afterFirst bool
	)

	buf.WriteRune('[')

	for k, v := range mapValues {
		if !afterFirst {
			buf.WriteRune(',')
			afterFirst = true
		}
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(v)
	}

	buf.WriteRune(']')

	return buf.String()
}
