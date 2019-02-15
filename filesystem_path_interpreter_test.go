// +build !nofstests

package peechee

import (
    "testing"
    "io/ioutil"
    "os"

    "github.com/djherbis/buffer"
	gm "github.com/onsi/gomega"
)

func TestFileSystemPathInterpreter(t *testing.T) {
	g := gm.NewGomegaWithT(t)

    fakeBytes := []byte(`Test fs.`)

	tempFile, err := ioutil.TempFile("", "foobar.txt")
    g.Expect(err).To(gm.BeNil())
	defer os.Remove(tempFile.Name())


	err = ioutil.WriteFile(tempFile.Name(), fakeBytes, 0644)
    g.Expect(err).To(gm.BeNil())

    fsInterp := filesystemPathInterpreter{}

    wa := buffer.New(100 * 1024)
    err = fsInterp.read(tempFile.Name(), wa)
    g.Expect(err).To(gm.BeNil())

    resBytes, err := ioutil.ReadAll(wa)
    g.Expect(err).To(gm.BeNil())

    g.Expect(resBytes).To(gm.Equal(fakeBytes))
}
