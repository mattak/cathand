package cathandtest

import (
	"github.com/mattak/cathand/pkg/cathand"
	"log"
	"os"
	"testing"
)

type FileUtilTestContext struct{}

func (FileUtilTestContext) setup() {
	os.RemoveAll("/tmp/test")
}

func (FileUtilTestContext) tearDown() {
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
	context := FileUtilTestContext{}
	context.setup()
	defer context.tearDown()

	checkMkdirp("/tmp/test")
	checkMkdirp("/tmp/test/a/b/c")
}

func TestRemoveFile(t *testing.T) {
	context := FileUtilTestContext{}
	context.setup()
	defer context.tearDown()

	os.MkdirAll("/tmp/test/a", os.ModePerm)
	os.MkdirAll("/tmp/test/b", os.ModePerm)

	if !isExists("/tmp/test/a") || !isExists("/tmp/test/b") {
		t.Fatal("error creating directory")
	}

	cathand.RemoveFile("/tmp/test")

	if isExists("/tmp/test/a") || isExists("/tmp/test/b") {
		t.Fatal("error cannot remove directory")
	}
}

func TestFileExists(t *testing.T) {
	context := FileUtilTestContext{}
	context.setup()
	defer context.tearDown()

	if cathand.FileExists("/tmp/test") {
		t.Fatal("error directory already exists")
	}

	os.MkdirAll("/tmp/test", os.ModePerm)

	if !cathand.FileExists("/tmp/test") {
		t.Fatal("error directory not exists")
	}
}
