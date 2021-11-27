package main

import (
    "errors"
    "fmt"
    "github.com/jkmathes/gluebin/lib"
    "github.com/urfave/cli/v2"
    "io"
    "log"
    "os"
    "io/ioutil"
    "path/filepath"
)
import _ "github.com/urfave/cli/v2"

func bail(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
    app := &cli.App{
        Name: "gluebin",
        Usage: "Create a static binary from a dynamic binary",
        Action: func(c *cli.Context) error {
            if c.Args().Len() < 1 {
                return errors.New("missing: binary file to convert")
            }
            fmt.Printf("Converting %q\n", c.Args().Get(0))
            bin := c.Args().Get(0)
            deps := lib.GetDependencies(bin)
            dir, err := ioutil.TempDir("", "gluebin")
            bail(err)
            fmt.Printf("Temp dir [%s]\n", dir)

            fmt.Printf("%q\n", lib.GetDependencies(bin))
            for _, dep := range deps {
                base := filepath.Base(dep)
                dest, err := os.Create(dir + "/" + base)
                bail(err)

                src, err := os.Open(dep)
                bail(err)

                _, err = io.Copy(dest, src)
                bail(err)


                bail(dest.Sync())
                bail(src.Close())
                bail(dest.Close())
            }
            return nil
        },
    }

    err := app.Run(os.Args)
    if err != nil {
        log.Fatal(err)
    }
}