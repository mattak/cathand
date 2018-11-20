package cathand

import (
	"github.com/labstack/gommon/log"
	"os"
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

func FileExists(directory string) bool {
	_, err := os.Stat(directory)
	return err == nil
}
