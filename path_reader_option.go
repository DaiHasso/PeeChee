package peechee

import (
)

type option func(*PathReader)

// WithFilesystem adds filesystem compatibility to the PathReader.
func WithFilesystem() option {
    return func(pathReader *PathReader) {
        pathReader.disk = filesystemPathInterpreter{}
    }
}

func applyOptions(pathReader *PathReader, options []option) {
    for _, opt := range options {
        opt(pathReader)
    }
}
