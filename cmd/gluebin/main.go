package main

import (
	"errors"
	"fmt"
	"github.com/jkmathes/gluebin/lib"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)
import _ "github.com/urfave/cli/v2"

func bail(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	instrumented, payload := lib.IsInstrumented()
	if instrumented {
		lib.ProxyExecutable(payload)
		return
	}
	app := &cli.App{
		Name:  "gluebin",
		Usage: "Create a static binary from a dynamic binary",
		Action: func(c *cli.Context) error {
			if c.Args().Len() < 2 {
				return errors.New("missing: binary file to convert and target binary to write")
			}

			bin := c.Args().Get(0)
			target := c.Args().Get(1)

			// TODO This needs to be more intelligent
			if bin == target {
				log.Fatal("src binary and target binary really shouldn't be the same")
			}

			buildBinary(bin, target)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func buildBinary(bin string, target string) {
	me, err := os.Executable()
	bail(err)
	deps, missing := lib.GetDependencies(bin)
	if len(missing) > 0 {
		log.Fatal("are you sure this binary was compiled here?")
	}

	fmt.Printf("Writing %q\n", target)
	dir, err := ioutil.TempDir("", "gluebin")
	bail(err)
	bail(os.MkdirAll(dir+"/libs", os.ModePerm))

	for _, dep := range deps {
		base := filepath.Base(dep)
		lib.CopyFile(dep, dir+"/libs/"+base)
	}
	lib.CopyFile(bin, dir+"/"+filepath.Base(bin))
	err = os.Chmod(dir+"/"+filepath.Base(bin), 0755)
	bail(err)

	p := lib.CreatePayload(dir, filepath.Base(bin))
	lib.AttachPayload(p, me, target)
	err = os.RemoveAll(dir)
	bail(err)
}
