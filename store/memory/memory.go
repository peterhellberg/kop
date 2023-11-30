package memory

import (
	"strings"

	"github.com/peterhellberg/kop/list"
)

func Store(vals ...string) list.Store {
	return newSet(fn(vals, strings.ToUpper)...)
}

func fn[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))

	for i, t := range ts {
		result[i] = fn(t)
	}

	return result
}
