package utils

import (
	"strconv"
	"strings"
)

func IdSliceToString(ids []uint64) string {
	var out []string

	for _, id := range ids {
		out = append(out, strconv.FormatUint(id, 10))
	}

	return strings.Join(out, " ")
}
