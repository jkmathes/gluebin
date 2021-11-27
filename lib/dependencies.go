package lib

import (
	"bytes"
	"debug/elf"
	"os/exec"
	"regexp"
	"strings"
)

func getLdConfig() map[string]string {
	m := make(map[string]string)
	cmd := exec.Command("ldconfig", "-p")
	buffer := new(bytes.Buffer)
	cmd.Stdout = buffer
	err := cmd.Run()
	bail(err)

	/**
	 * Example ldconfig line:
	 * libc.so.6 (libc6,x86-64, OS ABI: Linux 3.2.0) => /lib/x86_64-linux-gnu/libc.so.6
	 */
	re := regexp.MustCompile(`(?m)\s*(?P<key>[^\(\s]+)[^=]*=>\s*(?P<path>[^\s]+)`)
	lines := strings.Split(buffer.String(), "\n")
	for _, line := range lines {
		names := re.SubexpNames()
		matches := re.FindAllStringSubmatch(line, -1)
		/**
		 * Ignore summary lines, etc
		 */
		if len(matches) <= 0 {
			continue
		}
		matchLut := map[string]string{}
		for mi, match := range matches[0] {
			matchLut[names[mi]] = match
		}
		key := matchLut["key"]
		path := matchLut["path"]
		m[key] = path
	}
	return m
}

func getBinaryDependencies(name string) []string {
	e, err := elf.Open(name)
	bail(err)
	defer func(e *elf.File) {
		err := e.Close()
		bail(err)
	}(e)

	needed, err := e.DynString(elf.DT_NEEDED)
	return needed
}

func GetDependencies(name string) ([]string, []string) {
	ldConfig := getLdConfig()
	deps := map[string]string{}
	work := make([]string, 0)
	missing := make([]string, 0)

	/**
	 * Get the dynamic dependencies of the binary,
	 * then the dynamic dependencies of each dependency, etc
	 */
	work = append(work, getBinaryDependencies(name)...)

	for len(work) > 0 {
		top := len(work) - 1
		depName := work[top]
		work = work[:top]

		if len(deps[depName]) == 0 {
			deps[depName] = ldConfig[depName]
			if len(ldConfig[depName]) == 0 {
				missing = append(missing, depName)
			} else {
				work = append(work, getBinaryDependencies(ldConfig[depName])...)
			}
		}
	}

	r := make([]string, 0, len(deps))
	for _, v := range deps {
		r = append(r, v)
	}
	return r, missing
}
