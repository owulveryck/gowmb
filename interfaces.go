package gowmb

// Messager is an interface.
// The message can be setted via the set method.
// The message is encoded by the Serialize method and sent through the sockets
type Messager interface {
	Set(n []byte) error
	Serialize() ([]byte, error)
}

// The Tag for channel separation
type Tag interface {
	Parse(s string) error
}
