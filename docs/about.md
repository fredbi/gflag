# About CLI flags

There are many existing CLI flag handling libraries out there.

The most popular one by a mile is `github.com/spf13/pflag` (originally forked from `githb.com/ogier/pflag`).

It is used by other very popular packages `spf13/viper` and `spf13/cobra`.

The second most popular (by number of imports by public repositories) is `github.com/jessevdk/go-flags`.

The approach proposed by our package is built on top `github.com/spf13/pflag` and remains interoperable with other 
great CLI-building libraries such as `viper` and `cobra`.

[Various packages and approaches to dealing with flags](./approaches.md)
