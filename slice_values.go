package gflag

import (
	"bytes"
	"encoding"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
)

var (
	_ sliceOfValues            = &SliceValue[string]{}
	_ sliceOfValues            = &SliceValue[time.Duration]{}
	_ encoding.TextMarshaler   = &SliceValue[string]{}
	_ encoding.TextUnmarshaler = &SliceValue[string]{}

	rmQuote = strings.NewReplacer(`"`, "", `'`, "", "`", "")
)

type (
	// SliceValue is a generic type that implements github.com/spf13/pflag.Value and SliceValue.
	SliceValue[T FlaggablePrimitives | FlaggableTypes] struct {
		Value   *[]T
		changed bool
		*options
	}

	sliceOfValues interface {
		pflag.Value
		pflag.SliceValue
	}
)

// NewFlagSliceValue constructs a generic flag compatible with github.com/spf13/pflag.SliceValue.
//
// Since the flag type is inferred from the underlying data type, some flexibility allowed by pflag is not
// always possible at this point.
//
// For example, when T = []string, NewFlagSliceValue adopts the semantics of the pflag.StringSlice (with comma-separated values),
// whereas pflag also supports a StringArray flag.
//
// The underlying value, as []T, may be retrieved using GetSliceValue().
//
// In order to cover the full range of semantics offered by the pflag package, some options are available.
func NewFlagSliceValue[T FlaggablePrimitives | FlaggableTypes](addr *[]T, defaultValue []T, opts ...Option) *SliceValue[T] {
	if addr == nil {
		slice := make([]T, 0)
		addr = &slice
	}

	m := &SliceValue[T]{
		Value:   addr,
		options: defaultOptions(opts),
	}
	*m.Value = defaultValue

	return m
}

// GetValue returns the underlying value of the flag.
func (m SliceValue[T]) GetValue() []T {
	return *m.Value
}

// Get() implements the flag.Getter interface for programs using this Value from the
// standard "flag" package.
func (m SliceValue[T]) Get() any {
	return m.Value
}

func (m SliceValue[T]) String() string {
	if m.Value == nil {
		return ""
	}

	return writeAsSlice(m.GetSlice())
}

// Set knows how to config a string representation of the Value into a type T.
func (m *SliceValue[T]) Set(strValue string) error {
	return m.set(strValue)
}

func (m *SliceValue[T]) set(strValue string) error {
	slice, err := readAsCSV(rmQuote.Replace(strValue))
	if err != nil {
		return err
	}

	if !m.changed {
		// replace any default value
		err = m.Replace(slice)
		if err != nil {
			return err
		}

		m.changed = true

		return nil
	}

	// handle multiple occurences of the same flag with append semantics
	return m.append(slice...)
}

func (m *SliceValue[T]) Type() string {
	asAny := any(m.Value)
	switch v := asAny.(type) {
	case pflag.Value:
		return v.Type()
	case *[]string:
		switch m.semanticsForSliceString {
		case sliceStringIsArray:
			return "stringArray"
		default:
			return "stringSlice"
		}
	case *[]bool:
		return "boolSlice"
	case *[]int:
		return "intSlice"
	case *[]int8:
		return "int8Slice"
	case *[]int16:
		return "int16Slice"
	case *[]int32:
		return "int32Slice"
	case *[]int64:
		return "int64Slice"
	case *[]uint:
		return "uintSlice"
	case *[]uint16:
		return "uint16Slice"
	case *[]uint32:
		return "uint32Slice"
	case *[]uint64:
		return "uint64Slice"
	case *[]float32:
		return "float32Slice"
	case *[]float64:
		return "float64Slice"
	case *[]complex64:
		return "complex64Slice"
	case *[]complex128:
		return "complex128Slice"
	case *[]time.Duration:
		return "durationSlice"
	case *[]time.Time:
		return "timeSlice"
	case *[]net.IP:
		return "ipSlice"
	case *[]net.IPNet:
		return "ipNetSlice"
	case *[]net.IPMask:
		return "ipMaskSlice"
	default:
		return fmt.Sprintf("%T", *m.Value)
	}
}

// Append a single element to the SliceValue, from its string representation.
func (m *SliceValue[T]) Append(strValue string) error {
	asAny := any(m.Value)
	if sliceValue, ok := asAny.(pflag.SliceValue); ok {
		return sliceValue.Append(strValue)
	}

	return m.append(strValue)
}

func (m *SliceValue[T]) Replace(strValues []string) error {
	asAny := any(m.Value)
	if sliceValue, ok := asAny.(pflag.SliceValue); ok {
		return sliceValue.Replace(strValues)
	}

	*m.Value = (*m.Value)[:0]

	return m.append(strValues...)
}

func (m *SliceValue[T]) append(strValues ...string) error {
	asAny := any(m.Value)

	switch v := asAny.(type) {
	case *[]string:
		*v = append(*v, strValues...)
		*m.Value = *cast[[]T](v)
	case *[]bool:
		slice, err := buildSliceFromParser(strValues, strconv.ParseBool)
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]int:
		slice, err := buildSliceFromParser(strValues, intParser[int](0))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]int8:
		slice, err := buildSliceFromParser(strValues, intParser[int8](8))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]int16:
		slice, err := buildSliceFromParser(strValues, intParser[int16](16))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]int32:
		slice, err := buildSliceFromParser(strValues, intParser[int32](32))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]int64:
		slice, err := buildSliceFromParser(strValues, intParser[int64](64))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]uint:
		slice, err := buildSliceFromParser(strValues, uintParser[uint](0))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]uint16:
		slice, err := buildSliceFromParser(strValues, uintParser[uint16](16))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]uint32:
		slice, err := buildSliceFromParser(strValues, uintParser[uint32](32))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]uint64:
		slice, err := buildSliceFromParser(strValues, uintParser[uint64](64))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]float32:
		slice, err := buildSliceFromParser(strValues, floatParser[float32](32))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]float64:
		slice, err := buildSliceFromParser(strValues, floatParser[float64](64))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]complex64:
		slice, err := buildSliceFromParser(strValues, complexParser[complex64](64))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]complex128:
		slice, err := buildSliceFromParser(strValues, complexParser[complex128](128))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]time.Duration:
		slice, err := buildSliceFromParser(strValues, time.ParseDuration)
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]time.Time:
		slice, err := buildSliceFromParser(strValues, timeParser(m.timeFormats))
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]net.IP:
		slice, err := buildSliceFromParser(strValues, parseIP)
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]net.IPNet:
		slice, err := buildSliceFromParser(strValues, parseIPNet)
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	case *[]net.IPMask:
		slice, err := buildSliceFromParser(strValues, parseIPMask)
		if err != nil {
			return err
		}

		*v = append(*v, slice...)
		*m.Value = *cast[[]T](v)
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}

	return nil
}

// GetSlice return a []string representation of the slice values.
func (m *SliceValue[T]) GetSlice() []string {
	asAny := any(m.Value)

	switch v := asAny.(type) {
	case pflag.SliceValue:
		return v.GetSlice()
	case *[]string:
		return *v
	case *[]bool:
		return buildSliceWithFormatter(*v, strconv.FormatBool)
	case *[]int:
		return buildSliceWithFormatter(*v, formatInt[int])
	case *[]int8:
		return buildSliceWithFormatter(*v, formatInt[int8])
	case *[]int16:
		return buildSliceWithFormatter(*v, formatInt[int16])
	case *[]int32:
		return buildSliceWithFormatter(*v, formatInt[int32])
	case *[]int64:
		return buildSliceWithFormatter(*v, formatInt[int64])
	case *[]uint:
		return buildSliceWithFormatter(*v, formatUint[uint])
	case *[]uint16:
		return buildSliceWithFormatter(*v, formatUint[uint16])
	case *[]uint32:
		return buildSliceWithFormatter(*v, formatUint[uint32])
	case *[]uint64:
		return buildSliceWithFormatter(*v, formatUint[uint64])
	case *[]float32:
		return buildSliceWithFormatter(*v, floatFormatter[float32](32))
	case *[]float64:
		return buildSliceWithFormatter(*v, floatFormatter[float64](64))
	case *[]complex64:
		return buildSliceWithFormatter(*v, complexFormatter[complex64](64))
	case *[]complex128:
		return buildSliceWithFormatter(*v, complexFormatter[complex128](128))
	case *[]time.Duration:
		return buildSliceWithFormatter(*v, formatStringer[time.Duration])
	case *[]time.Time:
		return buildSliceWithFormatter(*v, timeFormatter(m.timeFormats[0]))
	case *[]net.IP:
		return buildSliceWithFormatter(*v, formatStringer[net.IP])
	case *[]net.IPNet:
		return buildSliceWithFormatter(*v, formatIPnet)
	case *[]net.IPMask:
		return buildSliceWithFormatter(*v, formatStringer[net.IPMask])
	default:
		panic(fmt.Sprintf("unsupported type: %T", v))
	}
}

func (m *SliceValue[T]) MarshalText() ([]byte, error) {
	return []byte(m.String()), nil
}

// MarshalFlag implements github.com/jessevdk/go-flags.Marshaler interface
func (m *SliceValue[T]) MarshalFlag() (string, error) {
	return m.String(), nil
}

func (m *SliceValue[T]) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte(`[`))
	text = bytes.TrimSuffix(text, []byte(`]`))

	return m.set(string(text))
}

// UnmarshalFlag implements github.com/jessevdk/go-flags.Unmarshaler interface
func (m *SliceValue[T]) UnmarshalFlag(value string) error {
	value = strings.TrimPrefix(value, `[`)
	value = strings.TrimSuffix(value, `]`)

	return m.set(value)
}

func readAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}

	csvReader := csv.NewReader(strings.NewReader(val))

	res, err := csvReader.Read()
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return res, nil
}

func writeAsCSV(values []string) (string, error) {
	b := &bytes.Buffer{}
	csvWriter := csv.NewWriter(b)

	if err := csvWriter.Write(values); err != nil {
		return "", err
	}

	csvWriter.Flush()

	return strings.TrimSuffix(b.String(), "\n"), csvWriter.Error()
}

func writeAsSlice(slice []string) string {
	out, _ := writeAsCSV(slice)

	return fmt.Sprintf("[%s]", out)
}

func buildSliceWithFormatter[T any](slice []T, formatter func(T) string) []string {
	rep := make([]string, 0, len(slice))
	for _, val := range slice {
		rep = append(rep, formatter(val))
	}

	return rep
}

func buildSliceFromParser[T any](slice []string, parser func(string) (T, error)) ([]T, error) {
	rep := make([]T, 0, len(slice))
	for _, strVal := range slice {
		val, err := parser(strVal)
		if err != nil {
			return nil, err
		}

		rep = append(rep, val)
	}

	return rep, nil
}
