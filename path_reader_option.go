package peechee

type option func(*PathReader)

// WithFilesystem adds filesystem compatibility to the PathReader.
func WithFilesystem() option {
    return func(pathReader *PathReader) {
        pathReader.disk = filesystemPathInterpreter{}
    }
}
