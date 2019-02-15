// +build nos3

package peechee

import (
    "io"

    "github.com/pkg/errors"
)


type s3PathInterpreter struct {
    s3Client interface{}
}

func (self s3PathInterpreter) read(path string) (io.Reader, error) {
    return nil, errors.New(
        "S3 compatibility is disabled, build without nos3 flag to enable S3 " +
            "compatibility",
    )
}
