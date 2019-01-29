package cathand

import (
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func MakeDirectory(outdir string) {
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		os.MkdirAll(outdir, os.ModePerm)
	}
}

func RemoveFile(outdir string) {
	if _, err := os.Stat(outdir); err != nil {
		return
	}

	err := os.RemoveAll(outdir)
	if err != nil {
		log.Fatal("cannot remove directory " + outdir)
	}
}

func ListUpFilePathes(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	results := []string{}

	for _, file := range files {
		if file.IsDir() {
			continue;
		}
		results = append(results, path.Join(dir, file.Name()))
	}

	return results, nil
}

func ListUpFileNamesWithoutExtension(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)

	if err != nil {
		return nil, err
	}

	results := []string{}

	for _, file := range files {
		if file.IsDir() {
			continue;
		}
		results = append(results, TrimExtension(file.Name()))
	}

	return results, nil
}

func TrimExtension(filename string) string {
	index := strings.LastIndex(filename,".")
	if index == -1 {
		return filename
	}

	return filename[0:index]
}

func FileExists(directory string) bool {
	_, err := os.Stat(directory)
	return err == nil
}
