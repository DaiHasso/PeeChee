// +build !nos3

package peechee

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "sync"
    "testing"

    "github.com/djherbis/buffer"
    "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3iface"
	gm "github.com/onsi/gomega"
)

var testBucket = "foo-config"
var fakeRequestId = "fake"
var reqDone bool
var fakeS3Bytes = []byte(`It worked!`)
var s3RW = new(sync.RWMutex)

type fakeS3 struct {
	s3iface.S3API
}

func (self fakeS3) GetObjectWithContext(
	aws.Context, *s3.GetObjectInput, ...request.Option,
) (*s3.GetObjectOutput, error) {
    s3RW.Lock()
    defer s3RW.Unlock()
	if reqDone {
		return nil, awserr.NewRequestFailure(
			nil, http.StatusRequestedRangeNotSatisfiable, fakeRequestId,
		)
	}
	reqDone = true
	output := new(s3.GetObjectOutput)
	output.Body = ioutil.NopCloser(bytes.NewReader(fakeS3Bytes))
	return output, nil
}


func TestTryS3Path(t *testing.T) {
	g := gm.NewGomegaWithT(t)
    fakeS3I := fakeS3{}
    defer func() {
        s3RW.Lock()
        reqDone = false
        s3RW.Unlock()
    }()
    path := fmt.Sprintf("s3://%s/fake-file.txt", testBucket)
    wa := buffer.New(100 * 1024)
    err := tryS3Path(fakeS3I, path, wa)
    g.Expect(err).To(gm.BeNil())
    bs, err := ioutil.ReadAll(wa)
    g.Expect(err).To(gm.BeNil())

    g.Expect(bs).To(gm.Equal(fakeS3Bytes))
}

func TestTryParseS3Uri(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    bucket := "foo-bucket"
    path := "foo/test.txt"

    expectedInput := &s3.GetObjectInput{
        Bucket: &bucket,
        Key:    &path,
    }
    actualInput, err := parseS3Uri(fmt.Sprintf("s3://%s/%s", bucket, path))
    g.Expect(err).To(gm.BeNil())

    g.Expect(actualInput).To(gm.Equal(expectedInput))
}

func TestTryParseS3UriFailure(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    bucket := "foo-bucket"
    s3Uri := fmt.Sprintf("s3://%s", bucket)

    actualInput, err := parseS3Uri(s3Uri)
    g.Expect(actualInput).To(gm.BeNil())
    g.Expect(err).ToNot(gm.BeNil())
    g.Expect(err).To(gm.MatchError(fmt.Sprintf(
        "Provided S3 uri '%s' did not match expected pattern", s3Uri,
    )))
}

func TestS3PathInterpreter(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    bucket := "foo-bucket"
    path := "foo/test.txt"

    fakeS3I := fakeS3{}
    defer func() {
        s3RW.Lock()
        reqDone = false
        s3RW.Unlock()
    }()
    s3Interp := s3PathInterpreter{
        s3Client: fakeS3I,
    }

    wa := buffer.New(100 * 1024)
    err := s3Interp.read(fmt.Sprintf("s3://%s/%s", bucket, path), wa)
    g.Expect(err).To(gm.BeNil())

    resBytes, err := ioutil.ReadAll(wa)
    g.Expect(err).To(gm.BeNil())

    g.Expect(resBytes).To(gm.Equal(fakeS3Bytes))
}

func TestPathReaderS3(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    fakeS3I := fakeS3{}
    defer func() {
        s3RW.Lock()
        reqDone = false
        s3RW.Unlock()
    }()
    pathReader := NewPathReader(WithS3(fakeS3I))

    res, err := pathReader.Read("s3://foo/bar.txt")
    g.Expect(err).To(gm.BeNil())

    resBytes, err := ioutil.ReadAll(res)
    g.Expect(err).To(gm.BeNil())

    g.Expect(resBytes).To(gm.Equal(fakeS3Bytes))
}

func TestPathReaderS3Option(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    fakeS3I := fakeS3{}
    pathReader := NewPathReader(WithS3(fakeS3I))

    g.Expect(pathReader.s3).ToNot(gm.BeNil())
}
