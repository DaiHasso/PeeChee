package peechee


import (
    "io"
    "io/ioutil"
    "testing"

	gm "github.com/onsi/gomega"
)

var fakeFSBytes = []byte(`Fantastic data!`)

type fakeFSPathInterpreter struct {}

func (self fakeFSPathInterpreter) read(path string, wa io.WriterAt) error {
    _, err := wa.WriteAt(fakeFSBytes, 0)
    return err
}


func TestPathReaderBasic(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    pathReader := NewPathReader(func(pathReader *PathReader) {
        pathReader.disk = fakeFSPathInterpreter{}
    })

    res, err := pathReader.Read("fake/file/path.txt")
    g.Expect(err).To(gm.BeNil())

    resBytes, err := ioutil.ReadAll(res)
    g.Expect(err).To(gm.BeNil())

    g.Expect(resBytes).To(gm.Equal(fakeFSBytes))
}

func TestPathReaderFSOption(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    pathReader := NewPathReader(WithFilesystem())

    g.Expect(pathReader.disk).ToNot(gm.BeNil())
}
