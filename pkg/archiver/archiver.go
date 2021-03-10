package archiver

import (
	"fmt"
    "path/filepath"
    "os"
)

func Archive(path string) {
    filepath.Walk(path, func (path string, file os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        fmt.Printf("%+v\n\n", file)
        return nil
    })
}
