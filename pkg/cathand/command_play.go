package cathand

import (
	"log"
	"sync"
)

func CommandPlay(projectName string) {
	playProject := NewProject(projectName, ".play", "")
	playSdcardProject := NewProject(projectName, ".play", "/sdcard/")
	resultProject := NewProject(projectName, ".result", "")
	resultSdcardProject := NewProject(projectName, ".result", "/sdcard/")

	// 0. verify
	{
		if !FileExists(playProject.RootDir) {
			log.Fatalln("Not found play directory: ", playProject.RootDir)
		}
	}

	// 1. push
	{
		RunWait("adb", "push", playProject.RootDir+"/", playSdcardProject.RootDir)
		RunWait("adb", "shell", "mkdir", "-p", resultSdcardProject.VideoDir)
	}

	// 2. record & play
	{

		var stopTrigger sync.WaitGroup
		stopTrigger.Add(1)
		sdcardVideoFilesCh := make(chan []string, 1)

		go RecordContinuouslyWithStopTrigger(&resultSdcardProject, sdcardVideoFilesCh, &stopTrigger)

		RunWait("adb", "shell", "sh", playSdcardProject.RunShellFile)

		stopTrigger.Done()
		<-sdcardVideoFilesCh

		// 3. aggregate results
		RemoveFile(resultProject.RootDir)
		RunWait("adb", "pull", resultSdcardProject.RootDir, resultProject.RootDir)
	}
}
