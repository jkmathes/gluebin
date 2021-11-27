package lib

import (
	_ "embed"
	"io"
	"os"
)

func CopyFile(src string, dest string) {
	d, err := os.Create(dest)
	bail(err)
	defer func(d *os.File) {
		bail(d.Close())
	}(d)

	s, err := os.Open(src)
	bail(err)
	defer func(d *os.File) {
		bail(d.Close())
	}(s)

	_, err = io.Copy(d, s)
	bail(err)
}
