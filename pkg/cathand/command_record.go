package cathand

import (
	"io/ioutil"
	"strings"
	"sync"
)

func GetInputDevices() []string {
	inputsRaw := RunWaitWrite("adb", "shell", "ls", "/dev/input")
	inputsStr := strings.Replace(NormalizeLineFeed(inputsRaw), "\n", " ", -1)

	inputs := strings.Split(inputsStr, " ")
	results := []string{}

	for i := 0; i < len(inputs); i++ {
		if len(inputs[i]) >= 2 {
			results = append(results, inputs[i])
		}
	}

	return results
}

func CommandRecord(projectName string) {
	project := NewProject(projectName, ".record", "")
	sdcardProject := NewProject(projectName, ".record", "/sdcard/")
	sdcardVideoFilesCh := make(chan []string)

	// 0. remove directory: sample.record
	RemoveFile(project.RootDir)

	// 1. create directory: sample.record
	MakeDirectory(project.RootDir)
	MakeDirectory(project.DeviceDir)
	MakeDirectory(project.VideoDir)
	RunWait("adb", "shell", "mkdir", "-p", sdcardProject.VideoDir)

	// 2. save getprop.log
	{
		bytes := RunWaitWrite("adb", "shell", "wm", "size")
		ioutil.WriteFile(project.SizeFile, NormalizeLineFeedBytes(bytes), 0644)
	}

	// 3. wait multi signal
	var eventWg sync.WaitGroup
	var recordWg sync.WaitGroup
	eventWg.Add(1)
	recordWg.Add(1)

	// 3.1. go data.bytes
	inputDevices := GetInputDevices()
	eventDataResults := make([]chan []byte, len(inputDevices))
	for i := 0; i < len(inputDevices); i++ {
		eventDataResults[i] = make(chan []byte)
		go RunWaitKillWrite(eventDataResults[i], &eventWg, "adb", "shell", "cat", "/dev/input/"+inputDevices[i])
	}

	// 3.2. go getevent.log
	eventTextResult := make(chan []byte)
	go RunWaitKillWrite(eventTextResult, &eventWg, "adb", "shell", "getevent", "-lt")

	// 3.3. go recording
	go RecordContinuously(&sdcardProject, sdcardVideoFilesCh)

	// 4. receive signal & stop 3 shells
	<-sdcardVideoFilesCh
	eventWg.Done()

	// 5. write results
	ioutil.WriteFile(project.EventFile, NormalizeLineFeedBytes(<-eventTextResult), 0644)
	for i := 0; i < len(inputDevices); i++ {
		ioutil.WriteFile(project.DeviceFile(inputDevices[i]), <-eventDataResults[i], 0644)
	}

	// 6. pull files
	RunWait("adb", "pull", sdcardProject.VideoDir, project.VideoDir)
	RunWait("adb", "shell", "rm", "-r", sdcardProject.RootDir)
}
