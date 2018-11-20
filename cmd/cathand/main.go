package main

import (
	"github.com/mattak/cathand/pkg/cathand"
	"os"
)

func usage() {
	panic(`
usage:
  record  <project>
  compose <project>
  play    <project>
example: 
  record  sample
  compose sample
  play    sample
`)
}

func main() {
	if len(os.Args) < 3 {
		usage()
	}

	command := os.Args[1]
	projectName := os.Args[2]

	if command == "record" {
		cathand.CommandRecord(projectName)
	} else if command == "play" {
		cathand.CommandPlay(projectName)
	} else if command == "compose" {
		cathand.CommandCompose(projectName)
	} else {
		usage()
	}
}
