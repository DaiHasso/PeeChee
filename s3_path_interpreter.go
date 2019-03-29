// +build !nos3

package peechee

import (
    "io"

    "github.com/pkg/errors"
    "github.com/daihasso/beagle"
	"github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

var s3UriRegex = beagle.MustRegex(
    `^s3:\/{2}(?P<bucket>[^\/]+)\/(?P<filepath>.*)`,
)

func parseS3Uri(s3Uri string) (*s3.GetObjectInput, error) {
    result := s3UriRegex.Match(s3Uri)
    if !result.Matched() {
        return nil, errors.Errorf(
            "Provided S3 uri '%s' did not match expected pattern", s3Uri,
        )
    }
    bucket := result.NamedGroup("bucket")[0]
    path := result.NamedGroup("filepath")[0]
    return &s3.GetObjectInput{
        Bucket: &bucket,
        Key:    &path,
    }, nil
}


func tryS3Path (s3c s3iface.S3API, path string, wa io.WriterAt) error {
    objectInput, err := parseS3Uri(path)
    if err != nil {
        return errors.Wrap(err, "Error while trying to parse S3 path provided")
    }

    downloader := s3manager.NewDownloaderWithClient(s3c)
    _, err = downloader.Download(wa, objectInput)
    if err != nil {
        return errors.Wrap(err, "Error while downloading path in S3")
    }

    return nil
}

type s3PathInterpreter struct {
    s3Client s3iface.S3API
}

func (self s3PathInterpreter) read(path string, wa io.WriterAt) error {
    return tryS3Path(self.s3Client, path, wa)
}


// WithS3 adds S3 compatiblity to the PathReader.
func WithS3(s3Client s3iface.S3API) option {
    return func(pathReader *PathReader) {
        pathReader.s3 = s3PathInterpreter{
            s3Client: s3Client,
        }
    }
}
