package inputs

import (
	"fmt"
	"io/ioutil"
)

type FileInput struct {
	Path string
}

func (f *FileInput) Read() string {
	file, err := ioutil.ReadFile(f.Path)
	if err == nil {
		return string(file)
	}

	panic(fmt.Sprintf("%s is not a valid file", f.Path))
}
