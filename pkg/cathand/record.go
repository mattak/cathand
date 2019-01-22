package cathand

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func RecordContinuously(sdcardProject *Project, sdcardFiles chan []string) {
	countCh := make(chan int)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)

	go RunSequentialSignalWait(signalCh, countCh, func(count int) *exec.Cmd {
		fmt.Println("adb", "shell", "screenrecord",
			sdcardProject.VideoFile(count),
			"--size", "1280x720",
			"--bit-rate", "4000000",
			"--time-limit", "10")

		return exec.Command("adb", "shell", "screenrecord",
			sdcardProject.VideoFile(count),
			"--size", "1280x720",
			"--bit-rate", "4000000",
			"--time-limit", "10")
	})

	<-countCh

	// XXX: hack wait for hardware processing last video file
	time.Sleep(500 * time.Millisecond)

	sdcardFiles <- nil
}

func RecordContinuouslyWithStopTrigger(sdcardProject *Project, sdcardFiles chan []string, stopTrigger *sync.WaitGroup) {
	countCh := make(chan int, 1)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)

	go RunSequentialSignalWait(signalCh, countCh, func(count int) *exec.Cmd {
		fmt.Println("adb", "shell", "screenrecord",
			sdcardProject.VideoFile(count),
			"--size", "1280x720",
			"--bit-rate", "4000000",
			"--time-limit", "10")

		return exec.Command("adb", "shell", "screenrecord",
			sdcardProject.VideoFile(count),
			"--size", "1280x720",
			"--bit-rate", "4000000",
			"--time-limit", "10")
	})

	stopTrigger.Wait()
	signalCh <- syscall.SIGINT

	<-countCh

	time.Sleep(500 * time.Millisecond)

	sdcardFiles <- nil
}
