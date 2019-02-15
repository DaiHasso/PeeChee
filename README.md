# PeeChee

[PeeChee?](https://en.wikipedia.org/wiki/Pee_Chee_folder)

## Description
PeeChee is a simple library for accessing filepaths (mostly)
backend-agnostically.


## Example

``` go
import (
    "io/ioutil"

    "github.com/daihasso/peechee"
)

var pathReader *peechee.PathReader

func main() {
    myS3Client := MakeMyS3Client()
    pathReader = peechee.NewPathReader(
        peechee.WithFilesystem(), peechee.WithS3(myS3Client),
    )

    localResultReader, err := pathReader.Read("my/local/config.yaml")
    if err != nil {
        panic(err)
    }

    localConfigBytes, err := ioutil.ReadAll(localResultReader)
    if err != nil {
        panic(err)
    }

    cloudResultReader, err := pathReader.Read("s3://my-bucket/config.yaml")
    if err != nil {
        panic(err)
    }

    cloudConfigBytes, err := ioutil.ReadAll(cloudResultReader)
    if err != nil {
        panic(err)
    }

    // Do stuff with all these configs.
}
```
