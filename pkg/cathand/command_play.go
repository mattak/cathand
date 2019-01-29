package cathand

import (
	"log"
	"strings"
	"sync"
)

func AdbRsync(fromLocalFolder string, toDeviceFolder string) {
	// NOTE: adb push raise error secure_mkdirs on nested directory. So do `find . -type d | xargs -n1 mkdir` before running
	{
		result := string(RunWaitWrite("find", fromLocalFolder, "-type", "d"))
		localDirs := strings.Split(result, "\n")

		arguments := []string{"shell", "mkdir", "-p"}
		for _, dir := range localDirs {
			deviceDir := strings.Replace(dir, fromLocalFolder, toDeviceFolder, -1)
			arguments = append(arguments, deviceDir)
		}
		RunWait("adb", arguments...)
	}

	{
		result := string(RunWaitWrite("find", fromLocalFolder, "-type", "f"))
		localFiles := strings.Split(result, "\n")

		for _, localFile := range localFiles {
			if len(localFile) < 1 {
				continue;
			}
			deviceFile := strings.Replace(localFile, fromLocalFolder, toDeviceFolder, -1)
			RunWait("adb", "push", localFile, deviceFile)
		}
	}
}

func CommandPlay(playProject Project, resultProject Project, recordOption RecordOption) {
	playDeviceProject := NewProject(playProject.Name, "/data/local/tmp/")
	resultDeviceProject := NewProject(resultProject.Name, "/data/local/tmp/")

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

		go RecordContinuouslyWithStopTrigger(&resultDeviceProject, sdcardVideoFilesCh, &stopTrigger, recordOption)

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
