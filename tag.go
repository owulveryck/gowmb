package gowmb

// The Tag for channel separation
type Tag interface {
	Parse(s string) error
}
