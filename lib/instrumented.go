package lib

import (
	"archive/tar"
	"compress/gzip"
	"debug/elf"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func bail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func IsInstrumented() bool {
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
		return false
	}

	b, err := j.Data()
	bail(err)

	fmt.Printf("%s\n", string(b[:]))
	return true
}

func ProxyExecutable() {

}

func addFile(tw *tar.Writer, file string, prefix string) {
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
		addFile(tw, dir + "/libs/" + file.Name(), "libs/")
	}

	addFile(tw, dir + "/" + clone, "")

	bail(tw.Close())
	bail(gw.Close())

	return t.Name()
}

func AttachPayload(payload string, orig string, target string) {
	// TODO Write this in go!
	cmd := exec.Command("objcopy", "--add-section", "gluebin_payload=" + payload, orig, target)
	bail(cmd.Run())
}
