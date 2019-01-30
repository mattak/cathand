package cathand

import (
	"fmt"
	"path"
)

type Project struct {
	Name         string
	RootDir      string
	VideoDir     string
	ImageDir     string
	InputDir     string
	BinDir       string
	RunShellFile string
	EventFile    string
	SizeFile     string
}

func NewProject(name string, prefix string) Project {
	root := prefix + name

	return Project{
		Name:         name,
		RootDir:      root,
		VideoDir:     root + "/video",
		ImageDir:     root + "/image",
		InputDir:     root + "/input",
		BinDir:       root + "/bin",
		RunShellFile: root + "/run.sh",
		EventFile:    root + "/event.log",
		SizeFile:     root + "/wm_size.log",
	}
}

func (p *Project) VideoFile(count int) string {
	return path.Join(p.VideoDir, fmt.Sprintf("%02d.mp4", count))
}

func (p *Project) InputFile(eventDriverName string) string {
	return path.Join(p.InputDir, fmt.Sprintf("%s.log", eventDriverName))
}

func (p *Project) InputFileWithoutRootDir(eventDriverName string) string {
	return fmt.Sprintf("input/%s.log", eventDriverName)
}

func (p *Project) ImageFileFormat(prefix string) string {
	return path.Join(p.ImageDir, prefix + "_%04d.png")
}
