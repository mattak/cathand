package cathandtest

import (
	"github.com/mattak/cathand/pkg/cathand"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
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

func TestListUpFilePathes(t *testing.T) {
	context := FileUtilTestContext{}
	context.setup()
	defer context.tearDown()

	os.MkdirAll("/tmp/test", os.ModePerm)
	ioutil.WriteFile("/tmp/test/1.txt", []byte("empty"), 0777)
	ioutil.WriteFile("/tmp/test/2.txt", []byte("empty"), 0777)

	files, err := cathand.ListUpFilePathes("/tmp/test")
	assert.NoError(t, err, "ListUpFiles does not raise error")
	assert.Equal(t, 2, len(files))
	assert.Equal(t, "/tmp/test/1.txt", files[0])
	assert.Equal(t, "/tmp/test/2.txt", files[1])
}

func TestTrimExtension(t *testing.T) {
	context := FileUtilTestContext{}
	context.setup()
	defer context.tearDown()

	assert.Equal(t, "/tmp/hoge", cathand.TrimExtension("/tmp/hoge.txt"))
	assert.Equal(t, "/tmp/hoge", cathand.TrimExtension("/tmp/hoge"))
}

func TestListUpFileNamesWithoutExtension(t *testing.T) {
	context := FileUtilTestContext{}
	context.setup()
	defer context.tearDown()

	os.MkdirAll("/tmp/test", os.ModePerm)
	ioutil.WriteFile("/tmp/test/1.txt", []byte("empty"), 0777)
	ioutil.WriteFile("/tmp/test/2.txt", []byte("empty"), 0777)

	files, err := cathand.ListUpFileNamesWithoutExtension("/tmp/test")
	assert.NoError(t, err, "ListUpFiles does not raise error")
	assert.Equal(t, 2, len(files))
	assert.Equal(t, "1", files[0])
	assert.Equal(t, "2", files[1])
}
