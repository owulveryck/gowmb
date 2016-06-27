package gowmb

import (
	"strconv"
)

// The Tag for channel separation
type Tag int

func stringToTag(s string) (Tag, error) {
	i, err := strconv.Atoi(s)
	return Tag(i), err
}
