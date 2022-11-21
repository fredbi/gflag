package gflag

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/exp/constraints"
)

func intParser[T constraints.Signed](bits int) func(in string) (T, error) {
	return func(in string) (T, error) {
		val, err := strconv.ParseInt(in, 0, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

func uintParser[T constraints.Unsigned](bits int) func(string) (T, error) {
	return func(in string) (T, error) {
		val, err := strconv.ParseUint(in, 0, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

func floatParser[T constraints.Float](bits int) func(string) (T, error) {
	return func(in string) (T, error) {
		val, err := strconv.ParseFloat(in, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

func complexParser[T constraints.Complex](bits int) func(string) (T, error) {
	return func(in string) (T, error) {
		val, err := strconv.ParseComplex(in, bits)
		if err != nil {
			return 0, err
		}

		return T(val), nil
	}
}

func parseIP(in string) (net.IP, error) {
	val := net.ParseIP(strings.TrimSpace(in))
	if val == nil {
		return nil, fmt.Errorf("failed to parse IP: %q", in)
	}

	return val, nil
}

func parseIPMask(in string) (net.IPMask, error) {
	val := pflag.ParseIPv4Mask(strings.TrimSpace(in))
	if val == nil {
		return nil, fmt.Errorf("failed to parse IP mask: %q", in)
	}

	return val, nil
}

func parseIPNet(in string) (net.IPNet, error) {
	_, val, err := net.ParseCIDR(strings.TrimSpace(in))
	if val == nil {
		return net.IPNet{}, fmt.Errorf("failed to parse CIDR: %q", in)
	}

	return *val, err
}

func parseBytes(strValue string, encoding bytesSemantics) ([]byte, error) {
	switch encoding {
	case bytesIsBase64:
		return base64.StdEncoding.DecodeString(strValue)
	default:
		return hex.DecodeString(strings.TrimSpace(strValue))
	}
}

func timeParser(formats []string) func(string) (time.Time, error) {
	return func(strValue string) (time.Time, error) {
		for _, format := range formats {
			ts, err := time.Parse(format, strValue)
			if err == nil {
				return ts, nil
			}
		}

		return time.Time{}, fmt.Errorf("invalid time format `%s` must be one of: %v", strValue, formats)
	}
}
