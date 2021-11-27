package lib

import (
	"embed"
	_ "embed"
	"io"
	"log"
	"os"
)

//go:embed resources/*
var resourceFs embed.FS

func LinuxAsm() string {
	return readFile("resources/linux.s")
}

func readFile(f string) string {
	data, err := resourceFs.ReadFile(f)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

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
