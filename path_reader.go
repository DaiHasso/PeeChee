package peechee

import (
    "io"
    "bytes"

    "github.com/pkg/errors"
)

type PathReader struct {
    s3,
    disk pathInterpreter
}

func (self PathReader) fromS3(path string, wa io.WriterAt) error {
    if self.s3 == nil {
        return errors.New("PathReader doesn't have S3 enabled")
    }
    return self.s3.read(path, wa)
}

func (self PathReader) fromDisk(path string, wa io.WriterAt) error {
    if self.disk == nil {
        return errors.New("PathReader doesn't have filesystem enabled")
    }
    return self.disk.read(path, wa)
}

// Read reads a path to a io.Reader for consumption later.
func (self PathReader) Read(path string) (io.Reader, error) {
    bufBytes := make([]byte, 0)
    wa := NewWriteAtBuffer(bufBytes[:])
    err := readPath(self, path, wa)
    if err != nil {
        return nil, errors.Wrap(err, "Error while reading path")
    }
    return bytes.NewReader(wa.Bytes()), nil
}

func (self *PathReader) AddOption(options ...option) {
    applyOptions(self, options)
}

// ReadTo reads a path to a provided io.WriterAt instance.
func (self PathReader) ReadTo(path string, wa io.WriterAt) error {
    return readPath(self, path, wa)
}

// NewPathReader creates a path reader with the provided options.
func NewPathReader(options ...option) *PathReader {
    pathReader := &PathReader{}

    applyOptions(pathReader, options)

    return pathReader
}
