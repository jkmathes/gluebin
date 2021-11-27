package lib

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"debug/elf"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func bail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func IsInstrumented() (bool, []byte) {
	me, err := os.Executable()
	bail(err)

	e, err := elf.Open(me)
	bail(err)
	defer func(e *elf.File) {
		err := e.Close()
		bail(err)
	}(e)

	j := e.Section("gluebin_payload")
	if j == nil {
		return false, nil
	}

	b, err := j.Data(); bail(err)
	return true, b
}

func ProxyExecutable(payload []byte) {
	home, err := os.UserHomeDir(); bail(err)
	dir := home + "/.gluebin/" + "blah"
	extractPayload(payload, dir)
}

func extractPayload(payload []byte, dir string) {
	err := os.MkdirAll(dir, os.ModePerm)
	bail(err)
	reader := bytes.NewReader(payload)
	gr, err := gzip.NewReader(reader)
	bail(err)
	tr := tar.NewReader(gr)
	bail(err)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else {
			bail(err)
		}

		fmt.Printf("%q\n", header)

		target := dir + "/" + header.Name
		if header.Typeflag == tar.TypeDir {
			err = os.MkdirAll(target, os.FileMode(header.Mode))
			bail(err)
		} else {
			err = os.MkdirAll(filepath.Dir(target), os.ModePerm)
			bail(err)
			f, err := os.Create(target)
			bail(err)
			_, err = io.Copy(f, tr)
			bail(err)
			err = f.Close()
			bail(err)
		}
		err = os.Chmod(target, os.FileMode(header.Mode))
		bail(err)
		err = os.Chtimes(target, header.AccessTime, header.ModTime)
		bail(err)
	}
}

func addFileToPayload(tw *tar.Writer, file string, prefix string) {
	f, err := os.Open(file)
	bail(err)
	defer func(f *os.File) {
		bail(f.Close())
	}(f)

	stat, err := f.Stat(); bail(err)
	h, err := tar.FileInfoHeader(stat, stat.Name()); bail(err)
	h.Name = prefix + stat.Name()
	bail(tw.WriteHeader(h))
	_, err = io.Copy(tw, f); bail(err)
}

func CreatePayload(dir string, clone string) string {
	files, err := ioutil.ReadDir(dir + "/libs")
	bail(err)

	t, err := ioutil.TempFile("", "payload*")

	gw := gzip.NewWriter(t)
	tw := tar.NewWriter(gw)

	for _, file := range files {
		fmt.Printf("Adding %s\n", file.Name())
		addFileToPayload(tw, dir + "/libs/" + file.Name(), "libs/")
	}

	addFileToPayload(tw, dir + "/" + clone, "")

	bail(tw.Close())
	bail(gw.Close())

	return t.Name()
}

func AttachPayload(payload string, orig string, target string) {
	// TODO Write this in go!
	cmd := exec.Command("objcopy", "--add-section", "gluebin_payload=" + payload, orig, target)
	bail(cmd.Run())
}
