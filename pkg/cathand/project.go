package cathand

import "fmt"

type Project struct {
	Name         string
	RootDir      string
	VideoDir     string
	InputDir     string
	DeviceDir    string
	BinDir       string
	RunShellFile string
	EventFile    string
	SizeFile     string
}

func NewProject(name string, extension string, prefix string) Project {
	root := prefix + name + extension

	return Project{
		Name:         name,
		RootDir:      root,
		VideoDir:     root + "/video",
		InputDir:     root + "/input",
		DeviceDir:    root + "/device",
		BinDir:       root + "/bin",
		RunShellFile: root + "/run.sh",
		EventFile:    root + "/event.log",
		SizeFile:     root + "/wm_size.log",
	}
}

func (p *Project) VideoFile(count int) string {
	return p.VideoDir + fmt.Sprintf("/%02d.mp4", count)
}

func (p *Project) InputFile(eventDriverName string) string {
	return p.InputDir + fmt.Sprintf("/%s.log", eventDriverName)
}

func (p *Project) InputFileWithoutRootDir(eventDriverName string) string {
	return fmt.Sprintf("input/%s.log", eventDriverName)
}

func (p *Project) DeviceFile(deviceName string) string {
	return p.DeviceDir + "/" + deviceName
}
