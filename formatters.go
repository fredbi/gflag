package gflag

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"time"

	"golang.org/x/exp/constraints"
)

func formatInt[T constraints.Signed](in T) string {
	return strconv.FormatInt(int64(in), 10)
}

func formatUint[T constraints.Unsigned](in T) string {
	return strconv.FormatUint(uint64(in), 10)
}

func floatFormatter[T constraints.Float](bits int) func(T) string {
	return func(in T) string {
		return strconv.FormatFloat(float64(in), 'g', -1, bits)
	}
}

func complexFormatter[T constraints.Complex](bits int) func(T) string {
	return func(in T) string {
		return strconv.FormatComplex(complex128(in), 'g', -1, bits)
	}
}

func bytesFormatter(in []byte, encoding bytesSemantics) string {
	switch encoding {
	case bytesIsBase64:
		return base64.StdEncoding.EncodeToString(in)
	default:
		return hex.EncodeToString(in)
	}
}

func timeFormatter(format string) func(time.Time) string {
	return func(ts time.Time) string {
		return ts.Format(format)
	}
}

func formatStringer[T fmt.Stringer](in T) string {
	return in.String()
}

func formatIPnet(in net.IPNet) string {
	return in.String()
}
