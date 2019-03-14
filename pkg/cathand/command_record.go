package cathand

import (
	"io/ioutil"
	"sync"
)

func CommandRecord(project Project, recordOption RecordOption) {
	sdcardProject := NewProject(project.Name, "/sdcard/")
	sdcardVideoFilesCh := make(chan []string)

	// 0. remove directory: sample.record
	RemoveFile(project.RootDir)

	// 1. create directory: sample.record
	MakeDirectory(project.RootDir)
	RunWait("adb", "shell", "rm", "-rf", sdcardProject.VideoDir)
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

	// 3.1. go getevent.log
	eventTextResult := make(chan []byte)
	go RunWaitKillWrite(eventTextResult, &eventWg, "adb", "shell", "getevent", "-t")

	// 3.2. go recording
	go RecordContinuously(&sdcardProject, sdcardVideoFilesCh, recordOption)

	// 4. receive signal & stop 3 shells
	<-sdcardVideoFilesCh
	eventWg.Done()

	// 5. write results
	ioutil.WriteFile(project.EventFile, NormalizeLineFeedBytes(<-eventTextResult), 0644)

	// 6. pull files
	RunWait("adb", "pull", sdcardProject.VideoDir + "/", project.VideoDir + "/")
	RunWait("adb", "shell", "rm", "-r", sdcardProject.RootDir)
}
