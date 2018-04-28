package socket

import "io"

type Unpacker interface {
	Unpack(r io.Reader) ([]byte, error)
}
