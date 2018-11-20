package cathand

import "fmt"

type Project struct {
	Name         string
	RootDir      string
	VideoDir     string
	InputDir     string
	DeviceDir    string
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
		RunShellFile: root + "/run.sh",
		EventFile:    root + "/event.log",
		SizeFile:     root + "/vm_size.log",
	}
}

func (p *Project) VideoFile(count int) string {
	return p.VideoDir + fmt.Sprintf("/%02d.mp4", count)
}

func (p *Project) InputFile(count int) string {
	return p.InputDir + fmt.Sprintf("/%05d.bytes", count)
}

func (p *Project) InputFileWithoutRootDir(count int) string {
	return fmt.Sprintf("input/%05d.bytes", count)
}

func (p *Project) DeviceFile(deviceName string) string {
	return p.DeviceDir + "/" + deviceName
}
