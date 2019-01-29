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

type RecordOption struct {
	Size      string
	BitRate   uint
	TimeLimit uint
}

func NewRecordOption() RecordOption {
	return RecordOption{
		Size:      "720x1280",
		BitRate:   4000000,
		TimeLimit: 10,
	}
}

func RecordContinuously(
	sdcardProject *Project, sdcardFiles chan []string, option RecordOption,
) {
	countCh := make(chan int)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)

	go RunSequentialSignalWait(signalCh, countCh, func(count int) *exec.Cmd {
		fmt.Println("adb", "shell", "screenrecord",
			sdcardProject.VideoFile(count),
			"--size", option.Size,
			"--bit-rate", fmt.Sprint(option.BitRate),
			"--time-limit", fmt.Sprint(option.TimeLimit))

		return exec.Command("adb", "shell", "screenrecord",
			sdcardProject.VideoFile(count),
			"--size", option.Size,
			"--bit-rate", fmt.Sprint(option.BitRate),
			"--time-limit", fmt.Sprint(option.TimeLimit))
	})

	<-countCh

	// XXX: hack wait for hardware processing of last video file
	time.Sleep(1000 * time.Millisecond)

	sdcardFiles <- nil
}

func RecordContinuouslyWithStopTrigger(
	sdcardProject *Project, sdcardFiles chan []string, stopTrigger *sync.WaitGroup,
	option RecordOption,
) {
	countCh := make(chan int, 1)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)

	go RunSequentialSignalWait(signalCh, countCh, func(count int) *exec.Cmd {
		fmt.Println("adb", "shell", "screenrecord",
			sdcardProject.VideoFile(count),
			"--size", option.Size,
			"--bit-rate", fmt.Sprint(option.BitRate),
			"--time-limit", fmt.Sprint(option.TimeLimit))

		return exec.Command("adb", "shell", "screenrecord",
			sdcardProject.VideoFile(count),
			"--size", option.Size,
			"--bit-rate", fmt.Sprint(option.BitRate),
			"--time-limit", fmt.Sprint(option.TimeLimit))
	})

	stopTrigger.Wait()

	// XXX: Wait for record end
	time.Sleep(1000 * time.Millisecond)

	signalCh <- syscall.SIGINT

	<-countCh

	// XXX: Wait for mp4 processing
	time.Sleep(1000 * time.Millisecond)

	sdcardFiles <- nil
}
