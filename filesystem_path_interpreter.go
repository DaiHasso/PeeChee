package peechee

import (
    "bufio"
    "io"
    "os"

    "github.com/pkg/errors"
)

type writeAtWrapper struct {
    io.WriterAt

    cur int64
}

func (self *writeAtWrapper) Write(p []byte) (n int, err error) {
    n, err = self.WriteAt(p, self.cur)
    if err != nil {
        return
    }
    self.cur += int64(n)

    return
}

type filesystemPathInterpreter struct {}

func (self filesystemPathInterpreter) read(
    path string, wa io.WriterAt,
) error {
    wrapper := &writeAtWrapper{
        WriterAt: wa, cur: 0,
    }
    handler, err := os.Open(path)
    defer handler.Close()
    if err != nil {
        return errors.Wrap(err, "Error while opening file")
    }
    r := bufio.NewReader(handler)
    w := bufio.NewWriter(wrapper)
    buf := make([]byte, 1024)
    for {
        // read a chunk
        n, err := r.Read(buf)
        if err != nil && err != io.EOF {
            return errors.Wrap(
                err, "Error while reading file data chunk",
            )
        }
        if n == 0 {
            break
        }

        // write a chunk
        if _, err := w.Write(buf[:n]); err != nil {
            return errors.Wrap(
                err, "Error while writing file data chunk to output",
            )
        }
    }

    if err = w.Flush(); err != nil {
        return errors.Wrap(err, "Error while writing file data to output")
    }

    return nil
}
