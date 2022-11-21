// Package gflag exposes generic types to deal with flags.
//
// Supported types are: all primitive golang types (e.g. int, uint, float64, complex128 ...) as well as a few commonly
// used types, which supported by pflag: time.Duration, net.IP, net.IPNet, net.IPMask, plus time.Time.
//
// The new types can be readily used as extensions to the github.com/spf13/flag package.
//
// All types currently provided by pflag can be obtained from gflag using generic types.
//
// gflag provides a more consistent interface across types, and can build single-valued flags,
// slice of flags as well as maps of flags for all underlying types.
//
// Notice that the []byte type is considered single-valued and that we don't support []uint8 as a slice of unsigned integer values (because []byte and []uint8 are aliases).
//
// There are a few attention points, though. Take a look at the examples to deal with these.
//   - []byte semantics as base64-encoded string require the option BytesIsBase64(true) to be passed to the flag constructor NewFlagValue().
//     The default behavior is to use hex-encoded strings.
//   - int semantics an incremental count require the option IntIsCount(true) to be passed to the flag constructor.
//   - []string semantics as "slice" (values passed as a list of CSV value) is the default. Adopting "array" semantics (each string value passed with
//     another instance of the flag) requires the StringSliceValueIsArray(true) option.
//
// Note to users of the github.com/spf13/viper package:
//
// Consuming flags from viper remain still possible.
//
// All pflag types currently supported by viper will work seamlessly with their gflag counterpart. Other types will be converted to a
// their string representation by viper.
package gflag
