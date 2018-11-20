package cathandtest

import (
	"github.com/mattak/cathand/pkg/cathand"
	"log"
	"os"
	"testing"
)

func Setup() {
	os.RemoveAll("/tmp/test")
}

func TearDown() {
	os.RemoveAll("/tmp/test")
}

func isExists(path string) bool {
	r, err := os.Stat(path)
	return err == nil && r != nil
}

func checkMkdirp(path string) {
	if isExists(path) {
		log.Fatal("already exists " + path)
	}

	cathand.MakeDirectory(path)

	if !isExists(path) {
		log.Fatal("no created " + path)
	}
}

func TestMkdirp(t *testing.T) {
	Setup()
	defer TearDown()

	checkMkdirp("/tmp/test")
	checkMkdirp("/tmp/test/a/b/c")
}

func TestRemoveFile(t *testing.T) {
	Setup()
	defer TearDown()

	os.MkdirAll("/tmp/test/a", os.ModePerm)
	os.MkdirAll("/tmp/test/b", os.ModePerm)

	if !isExists("/tmp/test/a") || !isExists("/tmp/test/b") {
		log.Fatal("error creating directory")
	}

	cathand.RemoveFile("/tmp/test")

	if isExists("/tmp/test/a") || isExists("/tmp/test/b") {
		log.Fatal("error cannot remove directory")
	}
}

func TestFileExists(t *testing.T) {
	Setup()
	defer TearDown()

	if cathand.FileExists("/tmp/test") {
		log.Fatal("error directory already exists")
	}

	os.MkdirAll("/tmp/test", os.ModePerm)

	if !cathand.FileExists("/tmp/test") {
		log.Fatal("error directory not exists")
	}
}