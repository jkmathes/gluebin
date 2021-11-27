package lib

import (
	"embed"
	_ "embed"
	"log"
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
