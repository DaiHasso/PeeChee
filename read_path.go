package peechee

import (
    "io"
    "regexp"

    "github.com/pkg/errors"
)

var s3ProtocolRegex = regexp.MustCompile(`^s3:\/{2}.*`)

func readPath(pathReader PathReader, path string, wa io.WriterAt) error {
    var (
        err error
    )
    if s3ProtocolRegex.MatchString(path) {
        err = pathReader.fromS3(path, wa)
    } else {
        err = pathReader.fromDisk(path, wa)
    }
    if err != nil {
        return errors.Wrapf(
            err, "Error while reading file from path '%s'", path,
        )
    }

    return nil
}
