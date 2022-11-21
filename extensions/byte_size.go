package extensions

import (
	"github.com/docker/go-units"
	"github.com/spf13/pflag"
)

var _ pflag.Value = new(ByteSizeValue)

// ByteSizeValue is used to pass human-readable byte sizes as flags.
type ByteSizeValue uint64

// NewByteSizeValue builds a ByteSizeValue.
func NewByteSizeValue(defaultValue uint64, p *uint64) *ByteSizeValue {
	*p = defaultValue

	return (*ByteSizeValue)(p)
}

// String method for a bytesize (pflag value and stringer interface)
func (b ByteSizeValue) String() string {
	return byteSizeFormatter(b)
}

// Set the value of this bytesize (pflag value interfaces)
func (b *ByteSizeValue) Set(value string) error {
	return b.set(value)
}

// MarshalFlag implements go-flags Marshaler interface
func (b ByteSizeValue) MarshalFlag() (string, error) {
	return byteSizeFormatter(b), nil
}

// MarshalText implements encoding.TextMarshaler interface
func (b ByteSizeValue) MarshalText() ([]byte, error) {
	return []byte(byteSizeFormatter(b)), nil
}

// UnmarshalFlag implements go-flags Unmarshaler interface
func (b *ByteSizeValue) UnmarshalFlag(value string) error {
	return b.set(value)
}

// UnmarshalText implements encoding.TextUnmarshaler interface
func (b ByteSizeValue) UnmarshalText(text []byte) error {
	return b.set(string(text))
}

func (b *ByteSizeValue) set(value string) error {
	sz, err := byteSizeParser(value)
	if err != nil {
		return err
	}

	*b = ByteSizeValue(sz)

	return nil
}

// Type returns the type of the pflag value (pflag value interface)
func (b *ByteSizeValue) Type() string {
	return "byte-size"
}

// ByteSizeVar defines a byte-size flag wih name, default value and usage string.
// The flag is set on the default pflag.CommandLine flagset.
//
// The flag value is stored at address p.
func ByteSizeVar(p *uint64, name string, defaultValue uint64, usage string) {
	ByteSizeVarP(p, name, "", defaultValue, usage)
}

// ByteSizeVarP is like ByteSize, and takes a shorthand for the flag name.
func ByteSizeVarP(p *uint64, name, shorthand string, defaultValue uint64, usage string) {
	v := NewByteSizeValue(defaultValue, p)
	_ = pflag.CommandLine.VarPF(v, name, shorthand, usage)
}

// ByteSize defines an uint64 flag with the specified name, default value, and usage string.
//
// The return value is the address of an uint64 variable that stores the value of the flag.
func ByteSize(name string, defaultValue uint64, usage string) *uint64 {
	return ByteSizeP(name, "", defaultValue, usage)
}

// ByteSizeP is like ByteSize, but accepts a shorthand letter that can be used after a single dash.
func ByteSizeP(name, shorthand string, defaultValue uint64, usage string) *uint64 {
	p := new(uint64)
	ByteSizeVarP(p, name, shorthand, defaultValue, usage)

	return p
}

func byteSizeFormatter(in ByteSizeValue) string {
	return units.HumanSize(float64(in))
}

func byteSizeParser(value string) (uint64, error) {
	sz, err := units.FromHumanSize(value)
	if err != nil {
		return 0, err
	}

	return uint64(sz), nil
}
