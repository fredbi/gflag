package gflag

import "github.com/spf13/pflag"

// AddValueFlag registers a new gflag.Value to a pflag.FlagSet and return the passed flag as well as its pflag.Flag version.
func AddValueFlag[T FlaggablePrimitives | FlaggableTypes](fs *pflag.FlagSet, flag *Value[T], name, shorthand, usage string) (*Value[T], *pflag.Flag) {
	return flag, fs.VarPF(flag, name, shorthand, usage)
}
func AddSliceValueFlag[T FlaggablePrimitives | FlaggableTypes](fs *pflag.FlagSet, flag *SliceValue[T], name, shorthand, usage string) (*SliceValue[T], *pflag.Flag) {
	return flag, fs.VarPF(flag, name, shorthand, usage)
}

func AddMapValueFlag[T FlaggablePrimitives | FlaggableTypes](fs *pflag.FlagSet, flag *MapValue[T], name, shorthand, usage string) (*MapValue[T], *pflag.Flag) {
	return flag, fs.VarPF(flag, name, shorthand, usage)
}
