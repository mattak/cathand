package cathand

import (
	"log"
	"strings"
	"sync"
)

func AdbRsync(fromLocalFolder string, toDeviceFolder string) {
	// NOTE: adb push raise error secure_mkdirs on nested directory. So do `find . -type d | xargs -n1 mkdir` before running
	result := string(RunWaitWrite("find", fromLocalFolder, "-type", "d"))
	localDirs := strings.Split(result, "\n")

	arguments := []string{"shell", "mkdir", "-p"}
	for _, dir := range localDirs {
		deviceDir := strings.Replace(dir, fromLocalFolder, toDeviceFolder, -1)
		arguments = append(arguments, deviceDir)
	}

	RunWait("adb", arguments...)
	RunWait("adb", "push", fromLocalFolder+"/", toDeviceFolder)
}

func CommandPlay(projectName string) {
	playProject := NewProject(projectName, ".play", "")
	playDeviceProject := NewProject(projectName, ".play", "/data/local/tmp/")
	resultProject := NewProject(projectName, ".result", "")
	resultDeviceProject := NewProject(projectName, ".result", "/data/local/tmp/")

	// 0. verify
	{
		if !FileExists(playProject.RootDir) {
			log.Fatalln("Not found play directory: ", playProject.RootDir)
		}
	}

	// 1. push
	{
		//RunWait("adb", "push", playProject.RootDir+"/", playDeviceProject.RootDir)
		AdbRsync(playProject.RootDir, playDeviceProject.RootDir)
		RunWait("adb", "shell", "mkdir", "-p", resultDeviceProject.VideoDir)
	}

	// 2. record & play
	{
		var stopTrigger sync.WaitGroup
		stopTrigger.Add(1)
		sdcardVideoFilesCh := make(chan []string, 1)

		go RecordContinuouslyWithStopTrigger(&resultDeviceProject, sdcardVideoFilesCh, &stopTrigger)

		RunWait("adb", "shell", "sh", playDeviceProject.RunShellFile)

		stopTrigger.Done()
		<-sdcardVideoFilesCh

		// 3. aggregate results
		RemoveFile(resultProject.RootDir)
		RunWait("adb", "pull", resultDeviceProject.RootDir, resultProject.RootDir)
	}

	// 3. cleanup
	{
		RunWait("adb", "shell", "rm", "-r", resultDeviceProject.RootDir)
		RunWait("adb", "shell", "rm", "-r", playDeviceProject.RootDir)
	}
}
