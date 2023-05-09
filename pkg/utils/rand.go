package utils

import (
	"math/rand"
	"strconv"
)

func RandID() string {
	return strconv.FormatUint(uint64(rand.Uint32()), 10)
}
