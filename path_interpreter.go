package peechee

import (
    "io"
)

type pathInterpreter interface {
    read(string, io.WriterAt) error
}
