package gflag

import "time"

type (
	Option func(*options)

	options struct {
		semanticsForBytes       bytesSemantics
		semanticsForInt         intSemantics
		semanticsForSliceString sliceStringSemantics
		timeFormats             []string
	}

	bytesSemantics       uint8
	intSemantics         uint8
	sliceStringSemantics uint8
)

const (
	bytesIsHex bytesSemantics = iota
	bytesIsBase64

	intIsInt intSemantics = iota
	intIsCount

	sliceStringIsSlice sliceStringSemantics = iota
	sliceStringIsArray
)

func defaultOptions(opts []Option) *options {
	o := &options{
		semanticsForBytes:       bytesIsHex,
		semanticsForInt:         intIsInt,
		semanticsForSliceString: sliceStringIsSlice,
		timeFormats:             []string{time.RFC3339Nano, time.RFC1123Z},
	}
	for _, apply := range opts {
		apply(o)
	}

	return o
}

// BytesIsBase64 adopts base64 encoding for flags of type []byte.
// The default is to encode []byte as an hex-encoded string.
func BytesIsBase64(enabled bool) Option {
	return func(o *options) {
		if enabled {
			o.semanticsForBytes = bytesIsBase64
		} else {
			o.semanticsForBytes = bytesIsHex
		}
	}
}

// IntIsCount adopts count semantics for flags of type int.
//
// Count semantics mean that multiple flags without a given value increment a counter.
// The default is to use plain integer.
func IntIsCount(enabled bool) Option {
	return func(o *options) {
		if enabled {
			o.semanticsForInt = intIsCount
		} else {
			o.semanticsForInt = intIsInt
		}
	}
}

// StringSliceIsArray adopts "array" semantics for flags of type []string,
// that is, we accumulate multiple instances of the flag, each with a single value rather
// than parsing a CSV list of values in a single flag.
//
// The default is to use a CSV list of values ("slice" semantics).
func StringSliceIsArray(enabled bool) Option {
	return func(o *options) {
		if enabled {
			o.semanticsForSliceString = sliceStringIsArray
		} else {
			o.semanticsForSliceString = sliceStringIsSlice
		}
	}
}

// WithTimeFormats define the formats supported to parse time.
//
// The first format specified will be used to render time values as strings.
//
// Defaults are: time.RFC3339Nano, time.RFC1123Z.
func WithTimeFormats(formats ...string) Option {
	return func(o *options) {
		if len(formats) > 0 {
			o.timeFormats = formats
		}
	}
}
